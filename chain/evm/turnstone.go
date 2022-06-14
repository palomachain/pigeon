package evm

import (
	"context"
	"math/big"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/chain"
	evmtypes "github.com/palomachain/sparrow/types/paloma/x/evm/types"
	"github.com/palomachain/sparrow/types/paloma/x/valset/types"
	valsettypes "github.com/palomachain/sparrow/types/paloma/x/valset/types"
	"github.com/palomachain/sparrow/util/slice"
	"github.com/vizualni/whoops"
)

const (
	maxPower = 1 << 32
)

type Turnstone struct {
	Client
}

type valset struct {
	Validators []common.Address
	Powers     []uint32
	ValsetID   *big.Int
}
type signature struct {
	R []byte
	S []byte
	V uint8
}
type consensus struct {
	Valset     valset
	Signatures []signature
}

func (t Turnstone) updateValset(
	ctx context.Context,
	newSnapshot *valsettypes.Snapshot,
	signatures []chain.ValidatorSignature,
) error {
	return whoops.Try(func() {
		valsetID, err := t.findLastValsetMessageID(ctx)
		whoops.Assert(err)

		snapshot, err := t.Client.paloma.GetSnapshotByID(ctx, valsetID)
		whoops.Assert(err)

		con := t.buildConsensus(ctx, snapshot, signatures)

		newValset := t.buildConsensus(ctx, newSnapshot, nil).Valset

		_, err = t.callSmartContract(ctx, "update_valset", []any{
			con,
			newValset,
		})
		whoops.Assert(err)

		return
	})
}

func (t Turnstone) submitLogicCall(
	ctx context.Context,
	messageID uint64,
	msg *evmtypes.ArbitrarySmartContractCall,
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

		snapshot, err := t.Client.paloma.GetSnapshotByID(ctx, valsetID)
		whoops.Assert(err)

		con := t.buildConsensus(ctx, snapshot, signatures)

		_, err = t.callSmartContract(ctx, "submit_logic_call", []any{
			con,
			common.HexToAddress(msg.HexAddress),
			msg.Payload,
		})
		whoops.Assert(err)

		return
	})
}

func (t Turnstone) findLastValsetMessageID(ctx context.Context) (uint64, error) {
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

func (t Turnstone) isArbitraryCallAlreadyExecuted(ctx context.Context, messageID uint64) (bool, error) {
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

func (t Turnstone) buildConsensus(
	ctx context.Context,
	snapshot *types.Snapshot,
	signatures []chain.ValidatorSignature,
) consensus {

	validators := snapshot.GetValidators()
	sort.Slice(validators, func(i, j int) bool {
		// we want a reverse sort! Higher powers go first!
		return validators[i].ShareCount.GTE(validators[j].ShareCount)
	})

	var signatureMap map[string]chain.ValidatorSignature
	if signatures != nil {
		signatureMap = slice.MakeMapKeys(
			signatures,
			func(sig chain.ValidatorSignature) string {
				return sig.ValAddress.String()
			},
		)
	}

	con := consensus{}
	con.Valset.ValsetID = 123

	totalShare := slice.Reduce(validators, func(prevSum sdk.Int, val types.Validator) sdk.Int {
		return prevSum.Add(val.ShareCount)
	}).Int64()

	for _, val := range validators {
		for _, ext := range val.GetExternalChainInfos() {

			if !(ext.ChainID == t.ChainID && ext.TurnstoneID == t.turnstoneID) {
				continue
			}

			power := maxPower * (val.ShareCount.Int64() / totalShare)
			// yay
			con.Valset.Powers = append(con.Valset.Powers, uint32(power))
			con.Valset.Validators = append(con.Valset.Validators, common.HexToAddress(ext.Address))

			if signatures != nil {
				sig, ok := signatureMap[val.Address.String()]
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

			// exiting the external chain info loop given that we have found
			// what we need
			break
		}
	}

	return con
}
