package evm

import (
	"context"
	"errors"
	goerrors "errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/VolumeFi/whoops"
	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/palomachain/paloma/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

const (
	maxPower            uint64 = 1 << 32
	powerThreshold      uint64 = 2_863_311_530
	SignedMessagePrefix        = "\x19Ethereum Signed Message:\n32"
)

//go:generate mockery --name=evmClienter --inpackage --testonly
type evmClienter interface {
	FilterLogs(ctx context.Context, fq etherum.FilterQuery, currBlockHeight *big.Int, fn func(logs []ethtypes.Log) bool) (bool, error)
	ExecuteSmartContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, addr common.Address, method string, arguments []any) (*etherumtypes.Transaction, error)
	DeployContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, bytecode, constructorInput []byte) (contractAddr common.Address, tx *ethtypes.Transaction, err error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (*ethtypes.Transaction, bool, error)

	BalanceAt(ctx context.Context, address common.Address, blockHeight uint64) (*big.Int, error)
	FindBlockNearestToTime(ctx context.Context, startingHeight uint64, when time.Time) (uint64, error)
	FindCurrentBlockNumber(ctx context.Context) (*big.Int, error)
	LastValsetID(ctx context.Context, addr common.Address) (*big.Int, error)
}

type compass struct {
	CompassID        string
	ChainReferenceID string

	compassAbi        *abi.ABI
	smartContractAddr common.Address
	paloma            PalomaClienter
	evm               evmClienter

	startingBlockHeight int64

	chainID *big.Int
}

func newCompassClient(
	smartContractAddrStr,
	compassID,
	chainReferenceID string,
	chainID *big.Int,
	compassAbi *abi.ABI,
	paloma PalomaClienter,
	evm evmClienter,
) compass {
	// if !ethcommon.IsHexAddress(smartContractAddrStr) {
	// 	whoops.Assert(errors.Unrecoverable(ErrInvalidAddress.Format(smartContractAddrStr)))
	// }
	return compass{
		CompassID:         compassID,
		ChainReferenceID:  chainReferenceID,
		smartContractAddr: common.HexToAddress(smartContractAddrStr),
		chainID:           chainID,
		compassAbi:        compassAbi,
		paloma:            paloma,
		evm:               evm,
	}
}

type Signature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}
type CompassValset struct {
	Validators []common.Address
	Powers     []*big.Int
	ValsetId   *big.Int
}
type CompassConsensus struct {
	Valset     CompassValset
	Signatures []Signature

	originalSignatures [][]byte
}

type CompassLogicCallArgs struct {
	LogicContractAddress common.Address
	Payload              []byte
}

func (c CompassConsensus) OriginalSignatures() [][]byte {
	return c.originalSignatures
}

func (t compass) updateValset(
	ctx context.Context,
	queueTypeName string,
	newValset *types.Valset,
	origMessage chain.MessageWithSignatures,
) (*ethtypes.Transaction, error) {
	return whoops.TryVal(func() *ethtypes.Transaction {
		currentValsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)
		logger := log.WithFields(log.Fields{
			"chain-reference-id": t.ChainReferenceID,
			"current-valset-id":  currentValsetID,
		})
		logger.Debug("update_valset")

		currentValset, err := t.paloma.QueryGetEVMValsetByID(ctx, currentValsetID, t.ChainReferenceID)
		whoops.Assert(err)

		if currentValset == nil {
			logger := log.WithFields(log.Fields{
				"current-valset-id": currentValsetID,
			})
			logger.Debug("current valset is empty")
			whoops.Assert(fmt.Errorf("oh no"))
		}

		consensusReached := isConsensusReached(currentValset, origMessage)
		if !consensusReached {
			logger := log.WithFields(log.Fields{
				"current-valset-id": currentValsetID,
			})
			logger.Debug("no consensus")
			whoops.Assert(ErrNoConsensus)
		}

		tx, err := t.callCompass(ctx, "update_valset", []any{
			BuildCompassConsensus(ctx, currentValset, origMessage.Signatures),
			TransformValsetToCompassValset(newValset),
		})
		if err != nil {
			logger := log.WithFields(log.Fields{
				"current-valset-id": currentValsetID,
			})
			logger.Debug("call_compass error")
			isSmartContractError := whoops.Must(t.SetErrorData(ctx, queueTypeName, origMessage.ID, err))
			if isSmartContractError {
				logger := log.WithFields(log.Fields{
					"current-valset-id": currentValsetID,
				})
				logger.Debug("smart contract error")
				return nil
			}
			whoops.Assert(err)
		}
		logger0 := log.WithFields(log.Fields{
			"current-valset-id": currentValsetID,
		})
		logger0.Debug("success")

		return tx
	})
}

func (t compass) submitLogicCall(
	ctx context.Context,
	queueTypeName string,
	msg *types.SubmitLogicCall,
	origMessage chain.MessageWithSignatures,
) (*ethtypes.Transaction, error) {
	return whoops.TryVal(func() *ethtypes.Transaction {
		executed, err := t.isArbitraryCallAlreadyExecuted(ctx, origMessage.ID)
		whoops.Assert(err)
		if executed {
			return nil
		}

		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
		whoops.Assert(err)

		consensusReached := isConsensusReached(valset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		con := BuildCompassConsensus(ctx, valset, origMessage.Signatures)
		compassArgs := CompassLogicCallArgs{
			LogicContractAddress: common.HexToAddress(msg.GetHexContractAddress()),
			Payload:              msg.GetPayload(),
		}

		tx, err := t.callCompass(ctx, "submit_logic_call", []any{
			con,
			compassArgs,
			new(big.Int).SetInt64(int64(origMessage.ID)),
			new(big.Int).SetInt64(msg.GetDeadline()),
		})
		if err != nil {
			isSmartContractError := whoops.Must(t.SetErrorData(ctx, queueTypeName, origMessage.ID, err))
			if isSmartContractError {
				return nil
			}
			whoops.Assert(err)
		}

		return tx
	})
}

func (t compass) uploadSmartContract(
	ctx context.Context,
	queueTypeName string,
	msg *types.UploadSmartContract,
	origMessage chain.MessageWithSignatures,
) (*ethtypes.Transaction, error) {
	return whoops.TryVal(func() *etherumtypes.Transaction {
		constructorInput := msg.GetConstructorInput()
		logger := log.WithFields(log.Fields{
			"chain-id":          t.ChainReferenceID,
			"constructor-input": constructorInput,
		})
		logger.Info("upload smart contract")

		contractABI, err := abi.JSON(strings.NewReader(msg.GetAbi()))
		if err != nil {
			logger.
				WithField("error", err.Error()).
				Error("uploadSmartContract: error parsing ABI")
		}
		// todo refactor that "assert" usage. Go way is returning error, rather than panic/recover as try/catch equivalent (it is not equivalent)
		whoops.Assert(err)

		// 0 means to get the latest valset
		latestValset, err := t.paloma.QueryGetEVMValsetByID(ctx, 0, t.ChainReferenceID)
		if err != nil {
			logger.
				WithField("error", err.Error()).
				Error("uploadSmartContract: error querying valset from Paloma")
		}
		whoops.Assert(err)

		consensusReached := isConsensusReached(latestValset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		constructorArgs, err := contractABI.Constructor.Inputs.Unpack(constructorInput)
		logger.
			WithField("args", constructorArgs).
			Info("uploadSmartContract: ABI contract constructor inputs unpack")
		if err != nil {
			logger.
				WithField("error", err.Error()).
				WithField("input", constructorInput).
				Error("uploadSmartContract: error unpacking ABI contract constructor inputs")
		}

		_, tx, err := t.evm.DeployContract(
			ctx,
			t.chainID,
			contractABI,
			msg.GetBytecode(),
			constructorInput,
		)
		if err != nil {
			logger.
				WithField("error", err.Error()).
				WithField("input", constructorInput).
				Error("uploadSmartContract: error calling DeployContract")

			isSmartContractError := whoops.Must(t.SetErrorData(ctx, queueTypeName, origMessage.ID, err))
			if isSmartContractError {
				logger.Info("uploadSmartContract: error calling DeployContract was a smart contract error")
				return nil
			}
			whoops.Assert(err)
		}

		return tx
	})
}

func (t compass) SetErrorData(ctx context.Context, queueTypeName string, msgID uint64, errToProcess error) (bool, error) {
	var jsonRpcErr rpc.DataError
	if !errors.As(errToProcess, &jsonRpcErr) {
		err := t.paloma.SetErrorData(ctx, queueTypeName, msgID, []byte(errToProcess.Error()))
		return false, err
	} else {
		log.WithFields(
			log.Fields{
				"queue-type-name": queueTypeName,
				"message-id":      msgID,
				"error-message":   jsonRpcErr.Error(),
			},
		).Warn("smart contract returned an error")

		err := t.paloma.SetErrorData(ctx, queueTypeName, msgID, []byte(jsonRpcErr.Error()))
		if err != nil {
			return false, err
		}

		return true, nil
	}
}

func (t compass) findLastValsetMessageID(ctx context.Context) (uint64, error) {
	log.Debug("fetching last valset message id")
	id, err := t.evm.LastValsetID(ctx, t.smartContractAddr)
	if err != nil {
		log.
			WithField("error", err.Error()).
			WithField("addr", t.smartContractAddr.String()).
			Error("error getting LastValsetID")
		return 0, fmt.Errorf("error getting LastValsetID")
	}

	return id.Uint64(), nil
}

func (t compass) isArbitraryCallAlreadyExecuted(ctx context.Context, messageID uint64) (bool, error) {
	blockNumber, err := t.evm.FindCurrentBlockNumber(ctx)
	if err != nil {
		return false, err
	}
	fromBlock := *big.NewInt(0)
	fromBlock.Sub(blockNumber, big.NewInt(9999))
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.smartContractAddr,
		},
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte("LogicCallEvent(address,bytes,uint256)")),
				common.Hash{},
				common.Hash{},
				crypto.Keccak256Hash(new(big.Int).SetInt64(int64(messageID)).Bytes()),
			},
		},
		FromBlock: blockNumber,
	}

	var found bool
	_, err = t.evm.FilterLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		found = len(logs) > 0
		return found
	})

	if err != nil {
		return false, err
	}

	return found, nil
}

func BuildCompassConsensus(
	ctx context.Context,
	v *types.Valset,
	signatures []chain.ValidatorSignature,
) CompassConsensus {
	signatureMap := slice.MakeMapKeys(
		signatures,
		func(sig chain.ValidatorSignature) string {
			return sig.SignedByAddress
		},
	)
	con := CompassConsensus{
		Valset: TransformValsetToCompassValset(v),
	}

	for i := range v.GetValidators() {
		sig, ok := signatureMap[v.GetValidators()[i]]
		if !ok {
			con.Signatures = append(con.Signatures,
				Signature{
					V: big.NewInt(0),
					R: big.NewInt(0),
					S: big.NewInt(0),
				})
		} else {
			con.Signatures = append(con.Signatures,
				Signature{
					V: new(big.Int).SetInt64(int64(sig.Signature[64]) + 27),
					R: new(big.Int).SetBytes(sig.Signature[:32]),
					S: new(big.Int).SetBytes(sig.Signature[32:64]),
				},
			)
		}
		con.originalSignatures = append(con.originalSignatures, sig.Signature)
	}

	return con
}

func (t compass) processMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	var gErr whoops.Group
	logger := log.WithField("queue-type-name", queueTypeName)
	for _, rawMsg := range msgs {
		logger = logger.WithField("message-id", rawMsg.ID)

		if ctx.Err() != nil {
			logger.Debug("exiting processing message context")
			break
		}

		var processingErr error
		var tx *ethtypes.Transaction
		msg := rawMsg.Msg.(*types.Message)
		logger := log.WithFields(log.Fields{
			"chain-reference-id": t.ChainReferenceID,
			"queue-name":         queueTypeName,
			"msg-id":             rawMsg.ID,
			"message-type":       fmt.Sprintf("%T", msg.GetAction()),
		})
		logger.Debug("processing")

		switch action := msg.GetAction().(type) {
		case *types.Message_SubmitLogicCall:
			tx, processingErr = t.submitLogicCall(
				ctx,
				queueTypeName,
				action.SubmitLogicCall,
				rawMsg,
			)
		case *types.Message_UpdateValset:
			logger := log.WithFields(log.Fields{
				"chain-reference-id":     t.ChainReferenceID,
				"queue-name":             queueTypeName,
				"msg-id":                 rawMsg.ID,
				"msg-bytes-to-sign":      rawMsg.BytesToSign,
				"msg-msg":                rawMsg.Msg,
				"msg-nonce":              rawMsg.Nonce,
				"msg-public-access-data": rawMsg.PublicAccessData,
				"message-type":           "Message_UpdateValset",
			})
			logger.Debug("switch-case-message-update-valset")
			tx, processingErr = t.updateValset(
				ctx,
				queueTypeName,
				action.UpdateValset.Valset,
				rawMsg,
			)
		case *types.Message_UploadSmartContract:
			logger := log.WithFields(log.Fields{
				"chain-reference-id":     t.ChainReferenceID,
				"queue-name":             queueTypeName,
				"msg-id":                 rawMsg.ID,
				"msg-bytes-to-sign":      rawMsg.BytesToSign,
				"msg-msg":                rawMsg.Msg,
				"msg-nonce":              rawMsg.Nonce,
				"msg-public-access-data": rawMsg.PublicAccessData,
				"message-type":           "Message_UploadSmartContract",
			})
			logger.Debug("switch-case-message-upload-contract")
			tx, processingErr = t.uploadSmartContract(
				ctx,
				queueTypeName,
				action.UploadSmartContract,
				rawMsg,
			)
		default:
			return ErrUnsupportedMessageType.Format(action)
		}

		processingErr = whoops.Enrich(
			processingErr,
			FieldMessageID.Val(rawMsg.ID),
			FieldMessageType.Val(msg.GetAction()),
		)

		switch {
		case processingErr == nil:
			if tx != nil {
				logger.Debug("setting public access data")
				if err := t.paloma.SetPublicAccessData(ctx, queueTypeName, rawMsg.ID, tx.Hash().Bytes()); err != nil {
					gErr.Add(err)
					return gErr
				}
			}
		case goerrors.Is(processingErr, ErrNoConsensus):
			// does nothing
		default:
			logger.WithError(processingErr).Error("processing error")
			gErr.Add(processingErr)
		}
	}

	return gErr.Return()
}

func (t compass) provideEvidenceForValidatorBalance(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	var g whoops.Group
	logger := log.WithField("queue-type-name", queueTypeName)
	logger.Debug("start processing validator balance request")
	for _, msg := range msgs {
		g.Add(
			whoops.Try(func() {
				vb := msg.Msg.(*types.ValidatorBalancesAttestation)
				height := whoops.Must(t.evm.FindBlockNearestToTime(ctx, uint64(t.startingBlockHeight), vb.FromBlockTime))

				logger1 := logger.WithFields(
					log.Fields{
						"height":          height,
						"nearest-to-time": vb.FromBlockTime,
					},
				)
				logger1.Debug("got height for time")

				res := &types.ValidatorBalancesAttestationRes{
					BlockHeight: height,
					Balances:    make([]string, 0, len(vb.HexAddresses)),
				}

				for _, addrHex := range vb.HexAddresses {
					addr := common.HexToAddress(addrHex)
					balance := whoops.Must(t.evm.BalanceAt(ctx, addr, height))
					logger1.WithFields(log.Fields{
						"evm-address": addr,
						"balance":     balance,
					}).Info("got balance")
					res.Balances = append(res.Balances, balance.Text(10))
				}

				whoops.Assert(t.paloma.AddMessageEvidence(ctx, queueTypeName, msg.ID, res))
			}),
		)
	}

	return g.Return()
}

// provideTxProof provides a very simple proof which is a transaction object
func (t compass) provideTxProof(ctx context.Context, queueTypeName string, rawMsg chain.MessageWithSignatures) error {
	log.WithFields(log.Fields{
		"queue-type-name":    queueTypeName,
		"msg-id":             rawMsg.ID,
		"public-access-data": rawMsg.PublicAccessData,
	}).Debug("providing proof")
	txHash := common.BytesToHash(rawMsg.PublicAccessData)
	tx, _, err := t.evm.TransactionByHash(ctx, txHash)
	if err != nil {
		return err
	}

	txProof, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	return t.paloma.AddMessageEvidence(ctx, queueTypeName, rawMsg.ID, &types.TxExecutedProof{
		SerializedTX: txProof,
	})
}

// provideErrorProof provides a pass-through proof for an error during relaying
func (t compass) provideErrorProof(ctx context.Context, queueTypeName string, rawMsg chain.MessageWithSignatures) error {
	log.WithFields(log.Fields{
		"queue-type-name":    queueTypeName,
		"msg-id":             rawMsg.ID,
		"public-access-data": rawMsg.PublicAccessData,
	}).Debug("providing error proof")

	return t.paloma.AddMessageEvidence(ctx, queueTypeName, rawMsg.ID, &types.SmartContractExecutionErrorProof{
		ErrorMessage: string(rawMsg.PublicAccessData),
	})
}

func TransformValsetToCompassValset(val *types.Valset) CompassValset {
	return CompassValset{
		Validators: slice.Map(val.GetValidators(), func(s string) common.Address {
			return common.HexToAddress(s)
		}),
		Powers: slice.Map(val.GetPowers(), func(p uint64) *big.Int {
			return big.NewInt(int64(p))
		}),
		ValsetId: big.NewInt(int64(val.GetValsetID())),
	}
}

func isConsensusReached(val *types.Valset, msg chain.MessageWithSignatures) bool {
	signaturesMap := make(map[string]chain.ValidatorSignature)
	for _, sig := range msg.Signatures {
		signaturesMap[sig.SignedByAddress] = sig
	}
	log.WithFields(log.Fields{
		"validators-size": len(val.Validators),
	}).Debug("confirming consensus reached")
	var s uint64
	for i := range val.Validators {
		val, pow := val.Validators[i], val.Powers[i]
		sig, ok := signaturesMap[val]
		log.WithFields(log.Fields{
			"i":         i,
			"validator": val,
			"power":     pow,
		}).Debug("checking consensus")
		if !ok {
			continue
		}
		bytesToVerify := crypto.Keccak256(append(
			[]byte(SignedMessagePrefix),
			msg.BytesToSign...,
		))
		recoveredPK, err := crypto.Ecrecover(bytesToVerify, sig.Signature)
		if err != nil {
			continue
		}
		log.WithFields(log.Fields{
			"i": i,
		}).Debug("good ecrecover")
		pk, err := crypto.UnmarshalPubkey(recoveredPK)
		if err != nil {
			continue
		}
		log.WithFields(log.Fields{
			"i": i,
		}).Debug("good unmarshal")
		recoveredAddr := crypto.PubkeyToAddress(*pk)
		if val == recoveredAddr.Hex() {
			s += pow
		}
		log.WithFields(log.Fields{
			"i": i,
		}).Debug("good consensus")
	}
	if s >= powerThreshold {
		return true
	}
	return false
}

func (c compass) callCompass(
	ctx context.Context,
	method string,
	arguments []any,
) (*etherumtypes.Transaction, error) {
	if c.compassAbi == nil {
		return nil, ErrABINotInitialized
	}
	return c.evm.ExecuteSmartContract(ctx, c.chainID, *c.compassAbi, c.smartContractAddr, method, arguments)
}
