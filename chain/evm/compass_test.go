package evm

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/palomachain/paloma/v2/x/evm/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/palomachain/pigeon/chain"
	evmmocks "github.com/palomachain/pigeon/chain/evm/mocks"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	errSample                 = whoops.String("oh no")
	testPowerThreshold uint64 = 2_863_311_530
)

type fakeJsonRpcError string

var _ rpc.DataError = fakeJsonRpcError("bla")

func (f fakeJsonRpcError) Error() string  { return string(f) }
func (f fakeJsonRpcError) ErrorData() any { return string(f) }

type StatusUpdater struct{}

func (s *StatusUpdater) WithLog(status string) paloma.StatusUpdater                        { return s }
func (s *StatusUpdater) WithMsg(msg *chain.MessageWithSignatures) paloma.StatusUpdater     { return s }
func (s *StatusUpdater) WithQueueType(queueType string) paloma.StatusUpdater               { return s }
func (s *StatusUpdater) WithChainReferenceID(chainReferenceID string) paloma.StatusUpdater { return s }
func (s *StatusUpdater) WithArg(key, value string) paloma.StatusUpdater                    { return s }
func (s *StatusUpdater) Info(ctx context.Context) error                                    { return nil }
func (s *StatusUpdater) Error(ctx context.Context) error                                   { return nil }
func (s *StatusUpdater) Debug(ctx context.Context) error                                   { return nil }

var (
	smartContractAddr        = common.HexToAddress("0xDEF")
	feeMgrContractAddr       = common.HexToAddress("0xDEADBEEF")
	ethCompatibleBytesToSign = crypto.Keccak256([]byte("sign me"))

	bobPK, _   = crypto.GenerateKey()
	alicePK, _ = crypto.GenerateKey()
	frankPK, _ = crypto.GenerateKey()

	sampleTx1 = func() *ethtypes.Transaction {
		sampleTx1RawBytes := common.FromHex(string(whoops.Must(os.ReadFile("testdata/sample-tx-raw.hex"))))
		tx := new(ethtypes.Transaction)
		whoops.Assert(tx.UnmarshalBinary(sampleTx1RawBytes))
		return tx
	}()

	sampleReceiptTx1 = &ethtypes.Receipt{
		Status: ethtypes.ReceiptStatusSuccessful,
	}

	eventIdAtomic = atomic.Int64{}
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
	const maxPower uint64 = 1 << 32

	if p > 1 || p < 0 {
		panic("invalid value for percentage")
	}

	return uint64(float64(maxPower) * p)
}

func getCompassABI(t *testing.T) *abi.ABI {
	file, err := os.Open("abi/compass/compass.abi")
	require.NoError(t, err)
	compassABI, err := abi.JSON(bufio.NewReader(file))
	require.NoError(t, err)

	return &compassABI
}

func buildSubmitLogicCallBytes(t *testing.T, messageID int64) []byte {
	arguments := abi.Arguments{
		// logic_contract_address
		{Type: whoops.Must(abi.NewType("address", "", nil))},
		// payload
		{Type: whoops.Must(abi.NewType("bytes", "", nil))},
		// message_id
		{Type: whoops.Must(abi.NewType("uint256", "", nil))},
		// event_id
		{Type: whoops.Must(abi.NewType("uint256", "", nil))},
	}

	examplePayload := []byte(``)

	bytes, err := arguments.Pack(
		common.HexToAddress("0x22786Ab8091D8E8EE6809ad17B83bE2df2Ed5E7a"),
		examplePayload,
		new(big.Int).SetInt64(messageID),
		new(big.Int).SetInt64(eventIdAtomic.Add(1)),
	)
	require.NoError(t, err)

	return bytes
}

func TestIsArbitraryCallAlreadyExecuted(t *testing.T) {
	tests := []struct {
		name          string
		messageID     int64
		setup         func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter)
		expected      bool
		expectedError error
	}{
		{
			name:      "False when unable to find current block number",
			messageID: 1,
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FindCurrentBlockNumber", mock.Anything).Times(1).Return(nil, errors.New("FindCurrentBlockNumber error"))

				return evm, paloma
			},
			expected:      false,
			expectedError: errors.New("FindCurrentBlockNumber error"),
		},
		{
			name:      "False when error filtering logs",
			messageID: 1,
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FindCurrentBlockNumber", mock.Anything).Times(1).Return(big.NewInt(5000), nil)
				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, errors.New("FilterLogs error"))

				return evm, paloma
			},
			expected:      false,
			expectedError: errors.New("FilterLogs error"),
		},
		{
			name:      "False when no logs found for filter",
			messageID: 1,
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FindCurrentBlockNumber", mock.Anything).Times(1).Return(big.NewInt(5000), nil)
				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil)

				return evm, paloma
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:      "False when not found in logs",
			messageID: 1,
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FindCurrentBlockNumber", mock.Anything).Times(1).Return(big.NewInt(5000), nil)

				isArbitraryCallExecutedLogs := []etherumtypes.Log{
					{
						BlockNumber: 1,
						Data:        buildSubmitLogicCallBytes(t, 2),
					},
					{
						BlockNumber: 2,
						Data:        buildSubmitLogicCallBytes(t, 3),
					},
				}

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(isArbitraryCallExecutedLogs)
				})

				return evm, paloma
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:      "True when found in logs",
			messageID: 1,
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FindCurrentBlockNumber", mock.Anything).Times(1).Return(big.NewInt(1), nil)

				isArbitraryCallExecutedLogs := []etherumtypes.Log{
					{
						BlockNumber: 1,
						Data:        buildSubmitLogicCallBytes(t, 2),
					},
					{
						BlockNumber: 2,
						Data:        buildSubmitLogicCallBytes(t, 1),
					},
				}

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(isArbitraryCallExecutedLogs)
				})

				return evm, paloma
			},
			expected:      true,
			expectedError: nil,
		},
	}

	asserter := assert.New(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evmClienter, palomaClienter := tt.setup(t)
			comp := newCompassClient(
				smartContractAddr.Hex(),
				feeMgrContractAddr.Hex(),
				"id-123",
				"internal-chain-id",
				big.NewInt(1),
				getCompassABI(t),
				palomaClienter,
				evmClienter,
			)

			ctx := context.Background()

			actual, actualError := comp.isArbitraryCallAlreadyExecuted(ctx, uint64(tt.messageID))
			asserter.Equal(tt.expected, actual)
			asserter.Equal(tt.expectedError, actualError)
		})
	}
}

func TestMessageProcessing(t *testing.T) {
	dummyErr := whoops.String("dummy")
	rpcErr := fakeJsonRpcError("bla")

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

	for _, tt := range []struct {
		name         string
		estimateOnly bool
		msgs         []chain.MessageWithSignatures
		setup        func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter)
		expErr       error
	}{
		{
			name: "submit_logic_call/message is already executed then it returns an error",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID: 666,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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

				paloma.On("NewStatus").Return(&StatusUpdater{})
				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(isArbitraryCallExecutedLogs)
				})

				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)

				sn := &valsettypes.Snapshot{
					Id: 55,
					Validators: []valsettypes.Validator{
						{
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainReferenceID: "internal-chain-id",
									Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
								},
							},
							Address: sdk.ValAddress("validator-1").Bytes(),
						},
					},
				}
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(sn, nil)
				paloma.On("SetErrorData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				return evm, paloma
			},
			expErr: ErrCallAlreadyExecuted,
		},
		{
			name: "submit_logic_call/happy path",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				evm.On("QueryUserFunds", mock.Anything, mock.Anything, mock.Anything).Return(
					big.NewInt(10_000_000),
					nil,
				)
				evm.On("SuggestGasPrice", mock.Anything).Return(
					big.NewInt(10),
					nil,
				)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{}, "submit_logic_call", mock.Anything).Return(
					tx,
					nil,
				)

				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), uint64(55), tx.Hash().Bytes()).Return(nil)
				return evm, paloma
			},
		},
		{
			name:         "estimate/submit_logic_call/happy path",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
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

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{estimateOnly: true}, "submit_logic_call", mock.Anything).Return(
					tx,
					nil,
				)

				return evm, paloma
			},
		},
		{
			name: "submit_logic_call/with target chain valset id not matching expected valset id, it should NOT return an error, but report log to Paloma",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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
				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 56,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 56,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("NewStatus").Return(&StatusUpdater{})
				return evm, paloma
			},
		},
		{
			name:         "estimate/submit_logic_call/with target chain valset id not matching expected valset id, it should NOT return an error, but report log to Paloma",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 56,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 56,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("NewStatus").Return(&StatusUpdater{})
				return evm, paloma
			},
		},
		{
			name: "submit_logic_call/happy path with mev relaying",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
									ExecutionRequirements: types.SubmitLogicCall_ExecutionRequirements{
										EnforceMEVRelay: true,
									},
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				evm.On("QueryUserFunds", mock.Anything, mock.Anything, mock.Anything).Return(
					big.NewInt(10_000_000),
					nil,
				)
				evm.On("SuggestGasPrice", mock.Anything).Return(
					big.NewInt(10),
					nil,
				)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{useMevRelay: true}, "submit_logic_call", mock.Anything).Return(
					tx,
					nil,
				)

				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), uint64(55), tx.Hash().Bytes()).Return(nil)
				return evm, paloma
			},
		},
		{
			name: "submit_logic_call/insufficient funds",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
								},
							},
						},
					},
					Signatures: []chain.ValidatorSignature{
						addValidSignature(bobPK),
					},
				},
			},
			expErr: fmt.Errorf("insufficient funds for fees: 500000 < 54000"),
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				evm.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn([]etherumtypes.Log{})
				})

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				evm.On("QueryUserFunds", mock.Anything, mock.Anything, mock.Anything).Return(
					big.NewInt(500_000),
					nil,
				)
				evm.On("SuggestGasPrice", mock.Anything).Return(
					big.NewInt(10),
					nil,
				)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)
				paloma.On("NewStatus").Return(&StatusUpdater{})
				paloma.On("SetErrorData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				return evm, paloma
			},
		},
		{
			name:         "estimate/submit_logic_call/happy path with mev relaying",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
							Action: &types.Message_SubmitLogicCall{
								SubmitLogicCall: &types.SubmitLogicCall{
									HexContractAddress: "0xABC",
									Abi:                []byte("abi"),
									Payload:            []byte("payload"),
									Deadline:           123,
									Fees: &types.Fees{
										RelayerFee:   50_000,
										CommunityFee: 3_000,
										SecurityFee:  1_000,
									},
									SenderAddress: sdk.AccAddress("message-sender").Bytes(),
									ExecutionRequirements: types.SubmitLogicCall_ExecutionRequirements{
										EnforceMEVRelay: true,
									},
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

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{useMevRelay: true, estimateOnly: true}, "submit_logic_call", mock.Anything).Return(
					tx,
					nil,
				)

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
							Assignee: sdk.ValAddress("validator-1").String(),
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				evm.On("FindCurrentBlockNumber", mock.Anything).Return(
					big.NewInt(0),
					nil,
				)

				paloma.On("NewStatus").Return(&StatusUpdater{})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.65),
							powerFromPercentage(0.35),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				return evm, paloma
			},
		},
		{
			name:         "estimate/submit_logic_call/without a consensus it returns",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("NewStatus").Return(&StatusUpdater{})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.65),
							powerFromPercentage(0.35),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)
				paloma.On("QueryGetLatestPublishedSnapshot", mock.Anything, "internal-chain-id").Return(
					&valsettypes.Snapshot{
						Id: 55,
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
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
							Assignee: sdk.ValAddress("validator-1").String(),
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

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
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{}, "update_valset", mock.Anything).Return(tx, nil)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), uint64(55), tx.Hash().Bytes()).Return(nil)
				return evm, paloma
			},
		},
		{
			name:         "estimate/update_valset/happy path",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

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
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				evm.On("ExecuteSmartContract", mock.Anything, chainID, mock.Anything, smartContractAddr, callOptions{estimateOnly: true}, "update_valset", mock.Anything).Return(tx, nil)
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
							Assignee: sdk.ValAddress("validator-1").String(),
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
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("NewStatus").Return(&StatusUpdater{})
				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
							crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.3),
							powerFromPercentage(0.3),
							powerFromPercentage(0.4),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Validators: []valsettypes.Validator{
							{
								ExternalChainInfos: []*valsettypes.ExternalChainInfo{
									{
										ChainReferenceID: "internal-chain-id",
										Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
									},
								},
								Address: sdk.ValAddress("validator-1").Bytes(),
							},
						},
					},
					nil,
				)

				return evm, paloma
			},
		},
		{
			name:         "estimate/update_valset/without a consensus it returns",
			estimateOnly: true,
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				currentValsetID := int64(55)

				evm.On("LastValsetID", mock.Anything, mock.Anything).Return(
					big.NewInt(55),
					nil,
				)

				paloma.On("NewStatus").Return(&StatusUpdater{})
				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
							crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.3),
							powerFromPercentage(0.3),
							powerFromPercentage(0.4),
						},
						ValsetID: uint64(currentValsetID),
					}, nil)
				paloma.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(&valsettypes.Snapshot{
					Validators: []valsettypes.Validator{
						{
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainReferenceID: "internal-chain-id",
									Address:          "0xDEADBEEF0ba39494ce839613fffba74279579268",
								},
							},
							Address: sdk.ValAddress("validator-1").Bytes(),
						},
					},
				}, nil)
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
							Assignee: sdk.ValAddress("validator-1").String(),
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
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(currentValsetID),
					},
					nil,
				)

				evm.On("DeployContract", mock.Anything, chainID, string(StoredContracts()["simple"].Source), []byte("bytecode"), []byte("constructor input")).Return(nil, tx, nil)

				paloma.On("SetPublicAccessData", mock.Anything, "queue-name", uint64(555), uint64(0), tx.Hash().Bytes()).Return(nil)
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
							Assignee: sdk.ValAddress("validator-1").String(),
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
						addValidSignature(alicePK),
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *evmmocks.PalomaClienter) {
				evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)

				currentValsetID := int64(0)
				paloma.On("NewStatus").Return(&StatusUpdater{})

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(currentValsetID), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{
							crypto.PubkeyToAddress(bobPK.PublicKey).Hex(),
							crypto.PubkeyToAddress(alicePK.PublicKey).Hex(),
							crypto.PubkeyToAddress(frankPK.PublicKey).Hex(),
						},
						Powers: []uint64{
							powerFromPercentage(0.3),
							powerFromPercentage(0.3),
							powerFromPercentage(0.4),
						},
						ValsetID: uint64(currentValsetID),
					},
					nil,
				)
				return evm, paloma
			},
		},
		{
			name: "upload_smart_contract/when smart contract returns an error it sends it to paloma as error data",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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

				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(0), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(55),
					},
					nil,
				)
				evm.On("DeployContract", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, rpcErr)
				paloma.On("NewStatus").Return(&StatusUpdater{})
				paloma.On("SetErrorData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return evm, paloma
			},
			expErr: rpcErr,
		},
		{
			name: "upload_smart_contract/when smart contract returns an error and sending it to paloma fails, it returns it back",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
						ID:          555,
						BytesToSign: ethCompatibleBytesToSign,
						Msg: &types.Message{
							Assignee: sdk.ValAddress("validator-1").String(),
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
				fakeErr := fakeJsonRpcError("bla")

				paloma.On("NewStatus").Return(&StatusUpdater{})
				paloma.On("QueryGetEVMValsetByID", mock.Anything, uint64(0), "internal-chain-id").Return(
					&types.Valset{
						Validators: []string{crypto.PubkeyToAddress(bobPK.PublicKey).Hex()},
						Powers:     []uint64{testPowerThreshold + 1},
						ValsetID:   uint64(55),
					},
					nil,
				)
				evm.On("DeployContract", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, fakeErr)
				paloma.On("SetErrorData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dummyErr)
				return evm, paloma
			},
			expErr: dummyErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			compassAbi := StoredContracts()["compass-evm"]
			ethClienter, palomaClienter := tt.setup(t)
			comp := newCompassClient(
				smartContractAddr.Hex(),
				feeMgrContractAddr.Hex(),
				"id-123",
				"internal-chain-id",
				chainID,
				&compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)

			_, err := comp.processMessages(ctx, "queue-name", tt.msgs, callOptions{estimateOnly: tt.estimateOnly})
			if tt.expErr != nil {
				require.ErrorContains(t, err, tt.expErr.Error())
			} else {
				require.ErrorIs(t, err, tt.expErr)
			}
		})
	}
}

func TestProcessingvalidatorBalancesRequest(t *testing.T) {
	ctx := context.Background()
	evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
	comp := newCompassClient(
		smartContractAddr.Hex(),
		feeMgrContractAddr.Hex(),
		"id-123",
		"internal-chain-id",
		big.NewInt(5),
		nil,
		paloma,
		evm,
	)
	comp.startingBlockHeight = 1233

	paloma.On("AddMessageEvidence", mock.Anything, "queue-name", uint64(1), &types.ValidatorBalancesAttestationRes{
		BlockHeight: 1212,
		Balances:    []string{"555", "666", "777"},
	}).Return(nil)

	evm.On("FindBlockNearestToTime", mock.Anything, uint64(comp.startingBlockHeight), time.Unix(123, 0)).Return(uint64(1212), nil)

	evm.On("BalanceAt", mock.Anything, common.HexToAddress("1"), uint64(1212)).Return(big.NewInt(555), nil)
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("2"), uint64(1212)).Return(big.NewInt(666), nil)
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("3"), uint64(1212)).Return(big.NewInt(777), nil)
	err := comp.provideEvidenceForValidatorBalance(ctx, "queue-name", []chain.MessageWithSignatures{
		{
			QueuedMessage: chain.QueuedMessage{
				ID: 1,
				Msg: &types.ValidatorBalancesAttestation{
					FromBlockTime: time.Unix(123, 0),
					HexAddresses:  []string{"1", "2", "3"},
				},
			},
		},
	})
	require.NoError(t, err)
}

func TestProcessingvalidatorBalancesRequestWithError(t *testing.T) {
	ctx := context.Background()
	evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
	comp := newCompassClient(
		smartContractAddr.Hex(),
		feeMgrContractAddr.Hex(),
		"id-123",
		"internal-chain-id",
		big.NewInt(5),
		nil,
		paloma,
		evm,
	)
	comp.startingBlockHeight = 1233

	paloma.On("AddMessageEvidence", mock.Anything, "queue-name", uint64(1), &types.ValidatorBalancesAttestationRes{
		BlockHeight: 1212,
		// The result set should include an empty field
		Balances: []string{"", "666", "777"},
	}).Return(nil)

	evm.On("FindBlockNearestToTime", mock.Anything, uint64(comp.startingBlockHeight), time.Unix(123, 0)).Return(uint64(1212), nil)

	// We fail to get the balance of the first validator
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("1"), uint64(1212)).Return(nil, errors.New("an error"))
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("2"), uint64(1212)).Return(big.NewInt(666), nil)
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("3"), uint64(1212)).Return(big.NewInt(777), nil)
	err := comp.provideEvidenceForValidatorBalance(ctx, "queue-name", []chain.MessageWithSignatures{
		{
			QueuedMessage: chain.QueuedMessage{
				ID: 1,
				Msg: &types.ValidatorBalancesAttestation{
					FromBlockTime: time.Unix(123, 0),
					HexAddresses:  []string{"1", "2", "3"},
				},
			},
		},
	})
	require.NoError(t, err)
}

func TestProcessingvalidatorBalancesRequestWithAllErrors(t *testing.T) {
	ctx := context.Background()
	evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
	comp := newCompassClient(
		smartContractAddr.Hex(),
		feeMgrContractAddr.Hex(),
		"id-123",
		"internal-chain-id",
		big.NewInt(5),
		nil,
		paloma,
		evm,
	)
	comp.startingBlockHeight = 1233

	evm.On("FindBlockNearestToTime", mock.Anything, uint64(comp.startingBlockHeight), time.Unix(123, 0)).Return(uint64(1212), nil)

	// We fail to get the balance of the first validator
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("1"), uint64(1212)).Return(nil, errors.New("an error"))
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("2"), uint64(1212)).Return(nil, errors.New("an error"))
	evm.On("BalanceAt", mock.Anything, common.HexToAddress("3"), uint64(1212)).Return(nil, errors.New("an error"))
	err := comp.provideEvidenceForValidatorBalance(ctx, "queue-name", []chain.MessageWithSignatures{
		{
			QueuedMessage: chain.QueuedMessage{
				ID: 1,
				Msg: &types.ValidatorBalancesAttestation{
					FromBlockTime: time.Unix(123, 0),
					HexAddresses:  []string{"1", "2", "3"},
				},
			},
		},
	})
	require.Error(t, err)
}

func TestProcessingReferenceBlockRequest(t *testing.T) {
	ctx := context.Background()
	conn := newMockEthClientConn(t)
	evm, paloma := newMockEvmClienter(t), evmmocks.NewPalomaClienter(t)
	comp := newCompassClient(
		smartContractAddr.Hex(),
		feeMgrContractAddr.Hex(),
		"id-123",
		"internal-chain-id",
		big.NewInt(5),
		nil,
		paloma,
		evm,
	)
	comp.startingBlockHeight = 1233

	paloma.On("AddMessageEvidence", mock.Anything, "queue-name", uint64(1), &types.ReferenceBlockAttestationRes{
		BlockHeight: 1400,
		BlockHash:   "0xa6d113404b5b2591c0a98b6ed1c6ca1760421c6fbb3571d942a47b131b88f51b",
	}).Return(nil)

	evm.On("FindBlockNearestToTime", mock.Anything, uint64(comp.startingBlockHeight), time.Unix(123, 0)).Return(uint64(1400), nil)
	evm.On("GetEthClient").Return(conn)

	conn.On("HeaderByNumber", mock.Anything, big.NewInt(1400)).
		Return(&ethtypes.Header{
			Number: big.NewInt(1400),
		}, nil)
	err := comp.provideEvidenceForReferenceBlock(ctx, "queue-name", []chain.MessageWithSignatures{
		{
			QueuedMessage: chain.QueuedMessage{
				ID: 1,
				Msg: &types.ReferenceBlockAttestation{
					FromBlockTime: time.Unix(123, 0),
				},
			},
		},
	})
	require.NoError(t, err)
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
				evm.On("TransactionReceipt", mock.Anything, mock.Anything).
					Return(sampleReceiptTx1, nil)

				paloma.On("AddMessageEvidence", mock.Anything, queue.QueueSuffixTurnstone, uint64(555),
					&types.TxExecutedProof{
						SerializedTX:      whoops.Must(sampleTx1.MarshalBinary()),
						SerializedReceipt: whoops.Must(sampleReceiptTx1.MarshalBinary()),
					},
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
				feeMgrContractAddr.Hex(),
				"id-123",
				"internal-chain-id",
				chainID,
				&compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)
			p := Processor{
				compass: &comp,
			}

			err := p.ProvideEvidence(ctx, queue.QueueSuffixTurnstone, tt.msgs)
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
					powerFromPercentage(0.3),
					powerFromPercentage(0.3),
					powerFromPercentage(0.4),
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
			res := isConsensusReached(context.Background(), tt.valset, tt.msgWithSig)
			require.Equal(t, tt.expRes, res)
		})
	}
}

func TestCompassBytesToPaloma(t *testing.T) {
	for _, tt := range []struct {
		name string
		src  [32]byte
		res  sdk.AccAddress
		err  bool
	}{
		{
			name: "20 byte address",
			src:  [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			res:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "20 byte address starting with 0",
			src:  [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			res:  []byte{0, 0, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "32 byte address",
			src:  [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			res:  []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "32 byte address starting with 0",
			src:  [32]byte{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			res:  []byte{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "invalid address",
			src:  [32]byte{},
			err:  true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			res, err := compassBytesToPalomaAddress(tt.src)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.res, res)
		})
	}
}
