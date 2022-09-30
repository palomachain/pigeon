package evm

import (
	"context"
	"errors"
	goerrors "errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/types/paloma/x/evm/types"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

const (
	maxPower            uint64 = 1 << 32
	powerThreshold      uint64 = 2_863_311_530
	SignedMessagePrefix        = "\x19Ethereum Signed Message:\n32"

	valsetUpdatedABISignature = "ValsetUpdated(bytes32,uint256)"
)

//go:generate mockery --name=evmClienter --inpackage --testonly
type evmClienter interface {
	FilterLogs(ctx context.Context, fq etherum.FilterQuery, currBlockHeight *big.Int, fn func(logs []ethtypes.Log) bool) (bool, error)
	ExecuteSmartContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, addr common.Address, method string, arguments []any) (*etherumtypes.Transaction, error)
	DeployContract(ctx context.Context, chainID *big.Int, contractAbi abi.ABI, bytecode, constructorInput []byte) (contractAddr common.Address, tx *ethtypes.Transaction, err error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (*ethtypes.Transaction, bool, error)

	BalanceAt(ctx context.Context, address common.Address, blockHeight uint64) (*big.Int, error)
	FindBlockNearestToTime(ctx context.Context, startingHeight uint64, when time.Time) (uint64, error)
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

type signature struct {
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
	Signatures []*big.Int

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

		currentValset, err := t.paloma.QueryGetEVMValsetByID(ctx, currentValsetID, t.ChainReferenceID)
		whoops.Assert(err)

		if currentValset == nil {
			whoops.Assert(fmt.Errorf("oh no"))
		}

		consensusReached := isConsensusReached(currentValset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		tx, err := t.callCompass(ctx, "update_valset", []any{
			BuildCompassConsensus(ctx, currentValset, origMessage.Signatures),
			TransformValsetToCompassValset(newValset),
		})
		if err != nil {
			isSmartContractError := whoops.Must(t.tryProvidingEvidenceIfSmartContractErr(ctx, queueTypeName, origMessage.ID, err))
			if isSmartContractError {
				return nil
			}
			whoops.Assert(err)
		}

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
			isSmartContractError := whoops.Must(t.tryProvidingEvidenceIfSmartContractErr(ctx, queueTypeName, origMessage.ID, err))
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
		contractABI, err := abi.JSON(strings.NewReader(msg.GetAbi()))
		whoops.Assert(err)

		// 0 means to get the latest valset
		latestValset, err := t.paloma.QueryGetEVMValsetByID(ctx, 0, t.ChainReferenceID)
		whoops.Assert(err)

		consensusReached := isConsensusReached(latestValset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		_, tx, err := t.evm.DeployContract(
			ctx,
			t.chainID,
			contractABI,
			msg.GetBytecode(),
			msg.GetConstructorInput(),
		)
		if err != nil {
			isSmartContractError := whoops.Must(t.tryProvidingEvidenceIfSmartContractErr(ctx, queueTypeName, origMessage.ID, err))
			if isSmartContractError {
				return nil
			}
			whoops.Assert(err)
		}

		return tx
	})
}

func (t compass) tryProvidingEvidenceIfSmartContractErr(ctx context.Context, queueTypeName string, msgID uint64, errToProcess error) (bool, error) {
	var jsonRpcErr rpc.DataError
	if !errors.As(errToProcess, &jsonRpcErr) {
		return false, nil
	}

	log.WithFields(
		log.Fields{
			"queue-type-name": queueTypeName,
			"message-id":      msgID,
			"error-message":   jsonRpcErr.Error(),
		},
	).Warn("smart contract returned an error")

	err := t.paloma.AddMessageEvidence(ctx, queueTypeName, msgID, &types.SmartContractExecutionErrorProof{
		ErrorMessage: jsonRpcErr.Error(),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (t compass) findLastValsetMessageID(ctx context.Context) (uint64, error) {
	log.Debug("fetching last valset message id")
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.smartContractAddr,
		},
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte(valsetUpdatedABISignature)),
			},
		},
		FromBlock: big.NewInt(t.startingBlockHeight),
	}

	latestMessageID := big.NewInt(0)

	var retErr error
	_, err := t.evm.FilterLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		for _, log := range logs {
			mm := make(map[string]any)
			err := t.compassAbi.Events["ValsetUpdated"].Inputs.UnpackIntoMap(mm, log.Data)
			if err != nil {
				retErr = err
				return false
			}
			id, ok := mm["valset_id"].(*big.Int)
			if !ok {
				retErr = ErrEvm.WrapS("valset_id should be big.Int, but it's not")
				return true
			}

			latestMessageID = id
			return true
		}
		return false
	})

	var g whoops.Group

	if latestMessageID.Uint64() == 0 {
		g.Add(ErrEvm.WrapS("could not find the valset_id in EVM logs"))
	}

	g.Add(retErr)
	g.Add(err)

	if g.Err() {
		return 0, g
	}
	log.WithField("valset_id", latestMessageID.Int64()).Debug("got valset_id")

	return uint64(latestMessageID.Int64()), nil
}

func (t compass) isArbitraryCallAlreadyExecuted(ctx context.Context, messageID uint64) (bool, error) {
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.smartContractAddr,
		},
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte("LagicCallEvent(address,bytes,uint256)")),
				common.Hash{},
				common.Hash{},
				crypto.Keccak256Hash(new(big.Int).SetInt64(int64(messageID)).Bytes()),
			},
		},
	}

	var found bool
	_, err := t.evm.FilterLogs(ctx, filter, big.NewInt(t.startingBlockHeight), func(logs []etherumtypes.Log) bool {
		found = len(logs) > 0
		return !found
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
			con.Signatures = append(con.Signatures, big.NewInt(0), big.NewInt(0), big.NewInt(0))
		} else {
			con.Signatures = append(con.Signatures,
				new(big.Int).SetInt64(int64(sig.Signature[64])+27),
				new(big.Int).SetBytes(sig.Signature[:32]),
				new(big.Int).SetBytes(sig.Signature[32:64]),
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

		if len(rawMsg.PublicAccessData) > 0 {
			logger.Warn("skipping the message as it already has public access data")
			continue
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
			tx, processingErr = t.updateValset(
				ctx,
				queueTypeName,
				action.UpdateValset.Valset,
				rawMsg,
			)
		case *types.Message_UploadSmartContract:
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
	var s uint64
	for i := range val.Validators {
		val, pow := val.Validators[i], val.Powers[i]
		sig, ok := signaturesMap[val]
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
		pk, err := crypto.UnmarshalPubkey(recoveredPK)
		if err != nil {
			continue
		}
		recoveredAddr := crypto.PubkeyToAddress(*pk)
		if val == recoveredAddr.Hex() {
			s += pow
		}
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
