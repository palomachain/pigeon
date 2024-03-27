package paloma_test

import (
	"context"
	"reflect"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	valsettypes "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/ion"
	"github.com/stretchr/testify/require"
)

var (
	keys     = []string{sdk.AccAddress("foo").String(), sdk.AccAddress("bar").String(), sdk.AccAddress("baz").String()}
	keyIdx   = 0
	usedKeys = make([]string, 0, 100)
)

type mockMsgSender struct {
	t           *testing.T
	error       error
	expectedMsg sdk.Msg
	calledMsg   sdk.Msg
}

func (m *mockMsgSender) SendMsg(ctx context.Context, msg sdk.Msg, memo string, opts ...ion.SendMsgOption) (*sdk.TxResponse, error) {
	if !reflect.DeepEqual(m.expectedMsg, msg) {
		m.t.Fatalf("unexpected argument 'msg', want: %+v, got: %+v", m.expectedMsg, msg)
	}

	m.calledMsg = msg

	return &sdk.TxResponse{TxHash: msg.String()}, m.error
}

type mockKeyRotator struct{}

func (m *mockKeyRotator) RotateKeys(context.Context) {
	keyIdx = (keyIdx + 1) % len(keys)
}

type mockMsg struct {
	info string
}

func (*mockMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.AccAddress("foobar"),
	}
}

func (*mockMsg) ProtoMessage()        {}
func (*mockMsg) Reset()               {}
func (*mockMsg) String() string       { return "" }
func (*mockMsg) ValidateBasic() error { return nil }

func Test_PalomaMessageSender_SendMsg(t *testing.T) {
	t.Run("must skip injecting metadata for messages with no need for it", func(t *testing.T) {
		ctx := context.Background()
		sender := &mockMsgSender{t: t}
		creator := "creator"
		signer := "signer"
		testee := paloma.NewPalomaMessageSender(&mockKeyRotator{}, sender).
			WithCreatorProvider(func() string { return creator }).
			WithSignerProvider(func() string { return signer })

		msg := &mockMsg{
			info: "foobar",
		}
		sender.expectedMsg = msg

		_, err := testee.SendMsg(ctx, msg, "")
		require.NoError(t, err)
		require.Equal(t, msg, sender.calledMsg)
	})

	t.Run("must inject metadata in messages that need it", func(t *testing.T) {
		ctx := context.Background()
		sender := &mockMsgSender{t: t}
		creator := sdk.AccAddress("creator").String()
		signer := sdk.AccAddress("signer").String()
		testee := paloma.NewPalomaMessageSender(&mockKeyRotator{}, sender).
			WithCreatorProvider(func() string { return creator }).
			WithSignerProvider(func() string { return signer })

		msg := &palomatypes.MsgAddStatusUpdate{
			Status: "bar",
			Level:  palomatypes.MsgAddStatusUpdate_LEVEL_INFO,
			Metadata: valsettypes.MsgMetadata{
				Creator: "foo",
			},
		}

		sender.expectedMsg = msg

		_, err := testee.SendMsg(ctx, msg, "")
		require.NoError(t, err)
		require.Equal(t, creator, msg.GetMetadata().GetCreator(), "must inject creator")
		require.Equal(t, []string{signer}, msg.GetMetadata().GetSigners(), "must inject signer address")
	})

	t.Run("must not reuse same key when called concurrently", func(t *testing.T) {
		ctx := context.Background()
		sender := &mockMsgSender{t: t}
		creator := "creator"
		testee := paloma.NewPalomaMessageSender(&mockKeyRotator{}, sender).
			WithCreatorProvider(func() string { return creator }).
			WithSignerProvider(func() string { return keys[keyIdx] })
		wg := &sync.WaitGroup{}

		msg := &palomatypes.MsgAddStatusUpdate{
			Status: "bar",
			Level:  palomatypes.MsgAddStatusUpdate_LEVEL_INFO,
			Metadata: valsettypes.MsgMetadata{
				Creator: "foo",
			},
		}
		sender.expectedMsg = msg

		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 33; j++ {
					r, err := testee.SendMsg(ctx, msg, "")
					usedKeys = append(usedKeys, r.TxHash)
					require.NoError(t, err)
				}
				wg.Done()
			}()
		}

		wg.Wait()

		for i, v := range usedKeys {
			require.Equal(t, v, usedKeys[i%len(keys)])
		}
	})
}
