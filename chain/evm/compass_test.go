package evm

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMessageProcessing(t *testing.T) {
	for _, tt := range []struct {
		name   string
		msgs   []chain.MessageWithSignatures
		setup  func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter)
		expErr error
	}{
		{
			name: "submit_logic_call/happy path",
			msgs: []chain.MessageWithSignatures{
				{
					QueuedMessage: chain.QueuedMessage{
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
						{
							ValAddress:      sdk.ValAddress("abc"),
							SignedByAddress: "bob",
							Signature:       []byte("abc"),
							PublicKey:       []byte("pk"),
						},
					},
				},
			},
			setup: func(t *testing.T) (*mockEvmClienter, *mockPalomaClienter) {
				eth, paloma := newMockEvmClienter(t), newMockPalomaClienter(t)

				valsetUpdatedLogs := []etherumtypes.Log{
					{
						BlockNumber: 1,
						log.Data(),
					},
				}

				eth.On("FilterLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Times(1).Return(false, nil).Run(func(args mock.Arguments) {
					fn := args.Get(3).(func([]etherumtypes.Log) bool)
					fn(valsetUpdatedLogs)
				})
				return eth, paloma
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			compassAbi := StoredContracts()["compass-evm"]
			ethClienter, palomaClienter := tt.setup(t)
			comp := newCompassClient(
				common.HexToAddress("0xDEF").Hex(),
				"id-123",
				"internal-chain-id",
				compassAbi.ABI,
				palomaClienter,
				ethClienter,
			)

			err := comp.processMessages(ctx, "any-queue-here", tt.msgs)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
