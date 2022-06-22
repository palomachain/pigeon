package evm

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vizualni/whoops"
)

func powerFromPercentage(p float64) uint64 {
	if p > 1 || p < 0 {
		panic("invalid value for percentage")
	}

	return uint64(float64(maxPower) * p)
}

func valsetUpdatedEvent(blockNumber uint64, checkpoint string, valsetID int64) etherumtypes.Log {
	compassABI := StoredContracts()["compass-evm"].ABI
	// checkpoint bytes32
	// valset_id uint256
	var cp32 [32]byte
	copy(cp32[:], []byte(checkpoint))
	data, err := compassABI.Events["ValsetUpdated"].Inputs.Pack(cp32, big.NewInt(valsetID))
	if err != nil {
		panic(err)
	}
	return etherumtypes.Log{
		BlockNumber: blockNumber,
		Data:        data,
	}
}

func TestMessageProcessing(t *testing.T) {
	smartContractAddr := common.HexToAddress("0xDEF")
	ethCompatibleBytesToSign := crypto.Keccak256([]byte("sign me"))

	bobPK, _ := crypto.GenerateKey()
	alicePK, _ := crypto.GenerateKey()
	frankPK, _ := crypto.GenerateKey()

	addValidSignature := func(pk *ecdsa.PrivateKey) chain.ValidatorSignature {
		return chain.ValidatorSignature{
			SignedByAddress: crypto.PubkeyToAddress(pk.PublicKey).Hex(),
			Signature: whoops.Must(crypto.Sign(
				crypto.Keccak256(
					append([]byte(SignedMessagePrefix), ethCompatibleBytesToSign...),
				), pk)),
			PublicKey: crypto.FromECDSAPub(&pk.PublicKey),
		}
	}

	for _, tt := range []struct {
		name   string
		msgs   []chain.MessageWithSignatures
		setup  func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter)
		expErr error
	}{
		{
			name: "submit_logic_call/message is already executed then it does nothing",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID: 666,
						Msg: &types.Message{
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{},
							},
						},
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				evm, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				isArbitraryCallExecutedLogs := []etherumtypes.Log{
					{
						BlockNumber: 1,
					},
				}

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(isArbitraryCallExecutedLogs)
				})

				paloma.On("DeleteJob", mock.Anything, "queue-name", uint64(666)).Return(nil)

				return evm, paloma
			},
		},
		{
			name: "submit_logic_call/happy path",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				evm, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{})
				})

				currentValsetID := int64(55)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{
						valsetUpdatedEvent(1, "abc", currentValsetID),
					})
				})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{powerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, mock.Anything, smartContractAddr, "submit_logic_call", mock.Anything).Return(nil, nil)

				paloma.On("DeleteJob", mock.Anything, "queue-name", uint64(555)).Return(nil)
				return evm, paloma
			},
		},
		{
			name: "submit_logic_call/without a consensus it returns",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				evm, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{})
				})

				currentValsetID := int64(55)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{
						valsetUpdatedEvent(1, "abc", currentValsetID),
					})
				})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{powerFromPercentage(0.1)},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("DeleteJob", mock.Anything, "queue-name", uint64(555)).Return(nil)

				return evm, paloma
			},
		},
		{
			name: "update_valset/happy path",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_UpdateValset{
								UpdateValset: &types.UpdateValset{
									Valset: &types.Valset{
										Validators: []string{
											crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
											crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
											crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
										},
										Powers: []uint64{
											powerFromPercentage(0.4),
											powerFromPercentage(0.2),
											powerFromPercentage(0.1),
										},
										ValsetID: 123,
									},
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
						addValidSignature(alicePK),
						// frank's signature is getting ignored but putting it
						// here just in case if there is a bug in the code
						addValidSignature(frankPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				evm, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				currentValsetID := int64(55)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{
						valsetUpdatedEvent(1, "abc", currentValsetID),
					})
				})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.5),
							powerFromPercentage(0.3),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, mock.Anything, smartContractAddr, "update_valset", mock.Anything).Return(nil, nil)

				paloma.On("DeleteJob", mock.Anything, "queue-name", uint64(555)).Return(nil)
				return evm, paloma
			},
		},
		{
			name: "update_valset/without a consensus it returns",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_UpdateValset{
								UpdateValset: &types.UpdateValset{
									Valset: &types.Valset{
										Validators: []string{
											crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
											crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
											crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
										},
										Powers: []uint64{
											powerFromPercentage(0.4),
											powerFromPercentage(0.2),
											powerFromPercentage(0.1),
										},
										ValsetID: 123,
									},
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
						addValidSignature(alicePK),
						// frank's signature is getting ignored but putting it
						// here just in case if there is a bug in the code
						addValidSignature(frankPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				evm, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				currentValsetID := int64(55)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{
						valsetUpdatedEvent(1, "abc", currentValsetID),
					})
				})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.3),
							powerFromPercentage(0.3),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)

				paloma.On("DeleteJob", mock.Anything, "queue-name", uint64(555)).Return(nil)
				return evm, paloma
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			compassAbi := StoredContracts()["compass-evm"]
			ethClienter, palomaClienter := tt.setup(t)
			comp := newCompassClient(
				smartContractAddr.Hex(),
				"id-123",
				"internal-chain-id",
				compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)

			err := comp.processMessages(ctx, "queue-name", tt.msgs)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
