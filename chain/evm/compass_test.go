package evm

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/pigeon/chain"
	evmmocks "github.com/palomachain/pigeon/chain/evm/mocks"
	"github.com/palomachain/pigeon/types/paloma/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vizualni/whoops"
)

const (
	errSample = whoops.String("oh no")
)

var (
	smartContractAddr        = common.HexToAddress("0xDEF")
	ethCompatibleBytesToSign = crypto.Keccak256([]byte("sign me"))

	bobPK, _   = crypto.GenerateKey()
	alicePK, _ = crypto.GenerateKey()
	frankPK, _ = crypto.GenerateKey()

	sampleTx1 = func() *ethtypes.Transaction {
		sampleTx1RawBytes := common.FromHex(string(whoops.Must(ioutil.ReadFile("testdata/sample-tx-raw.hex"))))
		tx := new(ethtypes.Transaction)
		whoops.Assert(tx.UnmarshalBinary(sampleTx1RawBytes))
		return tx
	}()
)

func signMessage(bz []byte, pk *ecdsa.PrivateKey) chain.ValidatorSignature {
	return chain.ValidatorSignature{
		SignedByAddress: crypto.PubkeyToAddress(pk.PublicKey).Hex(),
		Signature: whoops.Must(crypto.Sign(
			crypto.Keccak256(
				append([]byte(SignedMessagePrefix), bz...),
			), pk)),
		PublicKey: crypto.FromECDSAPub(&pk.PublicKey),
	}
}

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

	addValidSignature := func(pk *ecdsa.PrivateKey) chain.ValidatorSignature {
		return signMessage(ethCompatibleBytesToSign, pk)
	}
	chainID := big.NewInt(5)
	tx := etherumtypes.NewTransaction(
		5,
		common.HexToAddress("0x12"),
		big.NewInt(5),
		55,
		big.NewInt(5),
		[]byte("data"),
	)

	// txbz, err := tx.MarshalBinary()
	// require.NoError(t, err)

	for _, tt := range []struct {
		name   string
		msgs   []chain.MessageWithSignatures
		setup  func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter)
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
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				isArbitraryCallExecutedLogs := []etherumtypes.Log{
					{
						BlockNumber: 1,
					},
				}

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(isArbitraryCallExecutedLogs)
				})

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
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

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

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, "submit_logic_call", mock.Anything).Return(
					tx,
					nil,
				)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), tx.Hash().Bytes()).Return(nil)
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
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

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
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

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

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, "update_valset", mock.Anything).Return(tx, nil)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), tx.Hash().Bytes()).Return(nil)
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
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

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

				return evm, paloma
			},
		},
		{
			name: "upload_smart_contract/happy path",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_UploadSmartContract{
								UploadSmartContract: &types.UploadSmartContract{
									Bytecode:         []byte("bytecode"),
									Abi:              string(StoredContracts()["simple"].Source),
									ConstructorInput: []byte("constructor input"),
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				currentValsetID := int64(0)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{powerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)

				evm.On("DeployContract", mock.Anything, chainID, StoredContracts()["simple"].ABI, []byte("bytecode"), []byte("constructor input")).Return(nil, tx, nil)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), tx.Hash().Bytes()).Return(nil)
				return evm, paloma
			},
		},
		{
			name: "upload_smart_contract/without a consensus it returns an error",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Action: &types.Message_UploadSmartContract{
								UploadSmartContract: &types.UploadSmartContract{
									Bytecode:         []byte("bytecode"),
									Abi:              string(StoredContracts()["simple"].Source),
									ConstructorInput: []byte("constructor input"),
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				currentValsetID := int64(0)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{5},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
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
				chainID,
				&compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)

			err := comp.processMessages(ctx, "queue-name", tt.msgs)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}

func TestProvidingEvidenceForAMessage(t *testing.T) {

	addValidSignature := func(pk *ecdsa.PrivateKey) chain.ValidatorSignature {
		return signMessage(ethCompatibleBytesToSign, pk)
	}
	chainID := big.NewInt(5)

	for _, tt := range []struct {
		name   string
		msgs   []chain.MessageWithSignatures
		setup  func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter)
		expErr error
	}{
		{
			name: "with a public access data it tries to get the transaction and send it",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:               555,
						BytesToSign:      ethCompatibleBytesToSign,
						PublicAccessData: []byte("tx hash"),
						Msg: &types.Message{
							Action: &types.Message_SubmitLogicCall{SubmitLogicCall: &types.SubmitLogicCall{}},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
				evm.On("TransactionByHash", mock.Anything, mock.Anything).Return(sampleTx1, false, nil)

				paloma.On("AddMessageEvidence", mock.Anything, queueTurnstoneMessage, uint64(555),
					&types.TxExecutedProof{SerializedTX: whoops.Must(sampleTx1.MarshalBinary())},
				).Return(
					nil,
				)

				return evm, paloma
			},
		},
		{
			name: "with a public access data it tries to get the transaction, but it returns an error, and the error is returned back",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:               555,
						BytesToSign:      ethCompatibleBytesToSign,
						PublicAccessData: []byte("tx hash"),
						Msg: &types.Message{
							Action: &types.Message_SubmitLogicCall{SubmitLogicCall: &types.SubmitLogicCall{}},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
				evm.On("TransactionByHash", mock.Anything, mock.Anything).Return(nil, false, errSample)
				return evm, paloma
			},
			expErr: errSample,
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
				chainID,
				&compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)
			p := Processor{
				compass: comp,
			}

			err := p.ProvideEvidence(ctx, queueTurnstoneMessage, tt.msgs)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
func TestIfTheConsensusHasBeenReached(t *testing.T) {

	addValidSignature := func(pk *ecdsa.PrivateKey) chain.ValidatorSignature {
		return signMessage(ethCompatibleBytesToSign, pk)
	}

	for _, tt := range []struct {
		name       string
		valset     *types.Valset
		msgWithSig chain.MessageWithSignatures
		expRes     bool
	}{
		{
			name:   "all valid",
			expRes: true,
			valset: &types.Valset{
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
			},
			msgWithSig: chain.MessageWithSignatures{
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
					addValidSignature(alicePK),
					addValidSignature(frankPK),
				},
			},
		},
		{
			name:   "not enough to reach consensus",
			expRes: false,
			valset: &types.Valset{
				Validators: []string{
					crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
					crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
					crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
				},
				Powers: []uint64{
					powerFromPercentage(0.1),
					powerFromPercentage(0.2),
					powerFromPercentage(0.1),
				},
			},
			msgWithSig: chain.MessageWithSignatures{
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
					addValidSignature(alicePK),
					addValidSignature(frankPK),
				},
			},
		},
		{
			name:   "incorrect pk key",
			expRes: false,
			valset: &types.Valset{
				Validators: []string{
					crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
					crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
					crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
				},
				Powers: []uint64{
					powerFromPercentage(0.5),
					powerFromPercentage(0.2),
					powerFromPercentage(0.1),
				},
			},
			msgWithSig: chain.MessageWithSignatures{
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
					// this one is invalid
					{
						SignedByAddress: crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
						Signature:       []byte("this is invalid"),
						PublicKey:       crypto.FromECDSAPub(&bobPK.PublicKey),
					},
					addValidSignature(alicePK),
					addValidSignature(frankPK),
				},
			},
		},
		{
			name:   "valid signature, but incorrect signer",
			expRes: false,
			valset: &types.Valset{
				Validators: []string{
					crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
					crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
					crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
				},
				Powers: []uint64{
					powerFromPercentage(0.5),
					powerFromPercentage(0.2),
					powerFromPercentage(0.1),
				},
			},
			msgWithSig: chain.MessageWithSignatures{
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
					// this one should be bob
					addValidSignature(alicePK),
					addValidSignature(alicePK),
					addValidSignature(frankPK),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			res := isConsensusReached(tt.valset, tt.msgWithSig)
			require.Equal(t, tt.expRes, res)
		})
	}
}
