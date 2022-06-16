package evm

import (
	"context"
	"math/big"

	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	"github.com/palomachain/sparrow/util/slice"
	"github.com/vizualni/whoops"
)

type Compass struct {
	Client

	CompassID []byte
	ChainID   string
}

type signature struct {
	R []byte
	S []byte
	V uint8
}
type consensus struct {
	Valset     *types.Valset
	Signatures []signature
}

func (t Compass) updateValset(
	ctx context.Context,
	newValset *types.Valset,
	signatures []chain.ValidatorSignature,
) error {
	return whoops.Try(func() {
		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		currentValset, err := t.Client.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.internalChainID)
		whoops.Assert(err)

		con := t.buildConsensus(ctx, currentValset, signatures)

		_, err = t.callSmartContract(ctx, "update_valset", []any{
			con,
			newValset,
		})
		whoops.Assert(err)

		return
	})
}

func (t Compass) submitLogicCall(
	ctx context.Context,
	messageID uint64,
	msg *types.SubmitLogicCall,
	signatures []chain.ValidatorSignature,
) error {
	return whoops.Try(func() {
		executed, err := t.isArbitraryCallAlreadyExecuted(ctx, messageID)
		whoops.Assert(err)
		if executed {
			return
		}

		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		valset, err := t.Client.paloma.QueryGetEVMValsetByID(ctx, valsetID, t.internalChainID)
		whoops.Assert(err)

		con := t.buildConsensus(ctx, valset, signatures)

		_, err = t.callSmartContract(ctx, "submit_logic_call", []any{
			con,
			common.HexToAddress(msg.GetHexContractAddress()),
			msg.GetPayload(),
			msg.GetDeadline(),
		})
		whoops.Assert(err)

		return
	})
}

// func (t Compass) uploadSmartContract(
// 	ctx context.Context,
// 	messageID uint64,
// 	turnstoneID []byte,
// 	msg *types.UploadSmartContract,
// 	signatures []chain.ValidatorSignature,
// ) error {
// 	return whoops.Try(func() {
// 		executed, err := t.isArbitraryCallAlreadyExecuted(ctx, messageID)
// 		whoops.Assert(err)
// 		if executed {
// 			return
// 		}

// 		valsetID, err := t.findLastValsetMessageID(ctx)
// 		whoops.Assert(err)

// 		snapshot, err := t.Client.paloma.GetSnapshotByID(ctx, valsetID)
// 		whoops.Assert(err)

// 		bind.DeployContract()

// 		con := t.buildConsensus(ctx, snapshot, signatures)

// 		_, err = t.callSmartContract(ctx, "submit_logic_call", []any{
// 			con,
// 			common.HexToAddress(msg.GetHexContractAddress()),
// 			msg.GetPayload(),
// 			msg.GetDeadline(),
// 		})
// 		whoops.Assert(err)

// 		return
// 	})
// }

func (t Compass) findLastValsetMessageID(ctx context.Context) (uint64, error) {
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.turnstoneEVMContract,
		},
		Topics: [][]common.Hash{
			{
				crypto.Keccak256Hash([]byte("ValsetUpdated(bytes32,uint256)")),
			},
		},
	}

	var highestBlock uint64
	latestMessageID := big.NewInt(0)

	var retErr error
	_, err := t.processAllLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		for _, log := range logs {
			if log.BlockNumber > highestBlock {
				highestBlock = log.BlockNumber
				mm := make(map[string]any)
				err := t.smartContractAbi.Events["ValsetUpdated"].Inputs.UnpackIntoMap(mm, log.Data)
				if err != nil {
					retErr = err
					return false
				}
				id, ok := mm["valset_id"].(*big.Int)
				if !ok {
					panic("NOW WHAT")
				}

				if id.Cmp(latestMessageID) == 1 {
					latestMessageID = id
				}
			}
		}
		return true
	})

	if err != nil {
		return 0, err
	}

	if retErr != nil {
		return 0, retErr
	}

	return uint64(latestMessageID.Int64()), nil
}

func (t Compass) isArbitraryCallAlreadyExecuted(ctx context.Context, messageID uint64) (bool, error) {
	filter := etherum.FilterQuery{
		Addresses: []common.Address{
			t.turnstoneEVMContract,
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
	_, err := t.processAllLogs(ctx, filter, nil, func(logs []etherumtypes.Log) bool {
		found = true
		return false
	})

	if err != nil {
		return false, err
	}

	return found, nil
}

func (t Compass) buildConsensus(
	ctx context.Context,
	valset *types.Valset,
	signatures []chain.ValidatorSignature,
) consensus {

	signatureMap := slice.MakeMapKeys(
		signatures,
		func(sig chain.ValidatorSignature) string {
			return sig.SignedByAddress
		},
	)
	con := consensus{
		Valset: valset,
	}

	for i := range valset.HexAddress {
		sig, ok := signatureMap[valset.HexAddress[i]]
		if !ok {
			con.Signatures = append(con.Signatures, signature{})
		} else {
			con.Signatures = append(con.Signatures, signature{
				R: sig.Signature[:32],
				S: sig.Signature[32:64],
				V: uint8(int(sig.Signature[64])) + 27,
			})
		}
	}

	return con
}

func (t Compass) processMessages(ctx context.Context, msgs []chain.MessageWithSignatures) error {
	for _, rawMsg := range msgs {
		msg := rawMsg.Msg.(*types.Message)

		switch action := msg.GetAction().(type) {
		case *types.Message_SubmitLogicCall:
			return t.submitLogicCall(
				ctx,
				rawMsg.ID,
				action.SubmitLogicCall,
				rawMsg.Signatures,
			)
		case *types.Message_UpdateValset:
			return t.updateValset(
				ctx,
				action.UpdateValset.Valset,
				rawMsg.Signatures,
			)
		}
	}
	return nil
}
