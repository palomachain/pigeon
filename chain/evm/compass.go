package evm

import (
	"context"
	goerrors "errors"
	"fmt"
	"math/big"
	"strings"

	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
}

type compass struct {
	CompassID        string
	ChainReferenceID string

	compassAbi        *abi.ABI
	smartContractAddr common.Address
	paloma            palomaClienter
	evm               evmClienter

	chainID *big.Int
}

func newCompassClient(
	smartContractAddrStr,
	compassID,
	chainReferenceID string,
	chainID *big.Int,
	compassAbi *abi.ABI,
	paloma palomaClienter,
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
type valset struct {
	Validators []common.Address
	Powers     []*big.Int
	ValsetId   *big.Int
}
type consensus struct {
	Valset     valset
	Signatures []*big.Int
}

func (t compass) updateValset(
	ctx context.Context,
	newValset *types.Valset,
	origMessage chain.MessageWithSignatures,
) error {
	return whoops.Try(func() {
		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		currentValset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
		whoops.Assert(err)

		if currentValset == nil {
			whoops.Assert(fmt.Errorf("oh no"))
		}

		consensusReached := isConsensusReached(currentValset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		_, err = t.callCompass(ctx, "update_valset", []any{
			buildConsensus(ctx, currentValset, origMessage.Signatures),
			typesValsetToValset(newValset),
		})
		whoops.Assert(err)

		return
	})
}

func (t compass) submitLogicCall(
	ctx context.Context,
	msg *types.SubmitLogicCall,
	origMessage chain.MessageWithSignatures,
) error {
	return whoops.Try(func() {

		if len(origMessage.PublicProof) > 0 {
			data, pending, err := t.evm.TransactionByHash(ctx, origMessage.PublicProof)
			whoops.Assert(err)

			bz, err := data.MarshalJSON()
			whoops.Assert(err)

			return
		}

		executed, err := t.isArbitraryCallAlreadyExecuted(ctx, origMessage.ID)
		whoops.Assert(err)
		if executed {
			return
		}

		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		valset, err := t.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.ChainReferenceID)
		whoops.Assert(err)

		consensusReached := isConsensusReached(valset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		con := buildConsensus(ctx, valset, origMessage.Signatures)

		_, err = t.callCompass(ctx, "submit_logic_call", []any{
			con,
			common.HexToAddress(msg.GetHexContractAddress()),
			msg.GetPayload(),
			msg.GetDeadline(),
		})
		whoops.Assert(err)

		return
	})
}

func (t compass) uploadSmartContract(
	ctx context.Context,
	msg *types.UploadSmartContract,
	origMessage chain.MessageWithSignatures,
) error {
	return whoops.Try(func() {
		contractABI, err := abi.JSON(strings.NewReader(msg.GetAbi()))
		whoops.Assert(err)

		// 0 means to get the latest valset
		latestValset, err := t.paloma.QueryGetEVMValsetByID(ctx, 0, t.ChainReferenceID)
		whoops.Assert(err)

		consensusReached := isConsensusReached(latestValset, origMessage)
		if !consensusReached {
			whoops.Assert(ErrNoConsensus)
		}

		addr, tx, err := t.evm.DeployContract(
			ctx,
			t.chainID,
			contractABI,
			msg.GetBytecode(),
			msg.GetConstructorInput(),
		)
		// TODO: do attestation
		_ = addr
		_ = tx
		whoops.Assert(err)
		return
	})
}

func (t compass) findLastValsetMessageID(ctx context.Context) (uint64, error) {
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.smartContractAddr,
		},
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte(valsetUpdatedABISignature)),
			},
		},
	}

	var highestBlock uint64
	latestMessageID := big.NewInt(0)

	var retErr error
	_, err := t.evm.FilterLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		for _, log := range logs {
			if log.BlockNumber > highestBlock {
				highestBlock = log.BlockNumber
				mm := make(map[string]any)
				err := t.compassAbi.Events["ValsetUpdated"].Inputs.UnpackIntoMap(mm, log.Data)
				if err != nil {
					retErr = err
					return false
				}
				id, ok := mm["valset_id"].(*big.Int)
				if !ok {
					panic("valset_id should be big.Int, but it's not")
				}

				if id.Cmp(latestMessageID) == 1 {
					latestMessageID = id
				}
			}
		}
		return true
	})

	var g whoops.Group
	g.Add(retErr)
	g.Add(err)

	if g.Err() {
		return 0, g
	}

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
	_, err := t.evm.FilterLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		found = len(logs) > 0
		return !found
	})

	if err != nil {
		return false, err
	}

	return found, nil
}

func buildConsensus(
	ctx context.Context,
	v *types.Valset,
	signatures []chain.ValidatorSignature,
) consensus {

	signatureMap := slice.MakeMapKeys(
		signatures,
		func(sig chain.ValidatorSignature) string {
			return sig.SignedByAddress
		},
	)
	con := consensus{
		Valset: typesValsetToValset(v),
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
	}

	return con
}

func (t compass) processMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	var gErr whoops.Group
	for _, rawMsg := range msgs {

		if len(rawMsg.AccessableVia) > 0 {
			// maybe I don't have to process it anymore.

		}
		var processingErr error
		logger := log.WithFields(log.Fields{
			"chain-reference-id": t.ChainReferenceID,
			"queue-name":         queueTypeName,
			"msg-id":             rawMsg.ID,
		})
		logger.Info("processing")
		msg := rawMsg.Msg.(*types.Message)

		switch action := msg.GetAction().(type) {
		case *types.Message_SubmitLogicCall:
			processingErr = t.submitLogicCall(
				ctx,
				action.SubmitLogicCall,
				rawMsg,
			)
		case *types.Message_UpdateValset:
			processingErr = t.updateValset(
				ctx,
				action.UpdateValset.Valset,
				rawMsg,
			)
		case *types.Message_UploadSmartContract:
			processingErr = t.uploadSmartContract(
				ctx,
				action.UploadSmartContract,
				rawMsg,
			)
		default:
			return ErrUnsupportedMessageType.Format(action)
		}

		switch {
		case processingErr == nil:
			// TODO: this is temporary
			err := t.paloma.DeleteJob(ctx, queueTypeName, rawMsg.ID)
			gErr.Add(err)
		case goerrors.Is(processingErr, ErrNoConsensus):
			// does nothing
		default:
			gErr.Add(processingErr)
		}
	}

	if gErr.Err() {
		return gErr
	}

	return nil
}

func typesValsetToValset(val *types.Valset) valset {
	return valset{
		slice.Map(val.GetValidators(), func(s string) common.Address {
			return common.HexToAddress(s)
		}),
		slice.Map(val.GetPowers(), func(p uint64) *big.Int {
			return big.NewInt(int64(p))
		}),
		big.NewInt(int64(val.GetValsetID())),
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
