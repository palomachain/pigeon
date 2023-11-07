package paloma_test

import (
	"context"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/ion"
	"github.com/stretchr/testify/require"
)

type mockMsgSender struct {
	t           *testing.T
	result      *sdk.TxResponse
	error       error
	expectedMsg sdk.Msg
	calledMsg   sdk.Msg
}

func (m *mockMsgSender) SendMsg(ctx context.Context, msg sdk.Msg, memo string, opts ...ion.SendMsgOption) (*sdk.TxResponse, error) {
	if !reflect.DeepEqual(m.expectedMsg, msg) {
		m.t.Fatalf("unexpected argument 'msg', want: %+v, got: %+v", m.expectedMsg, msg)
	}

	m.calledMsg = msg

	return m.result, m.error
}

type mockKeyRotator struct{}

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

func (m *mockKeyRotator) RotateKeys(context.Context) string { return "foobar" }

func Test_PalomaMessageSender_SendMsg(t *testing.T) {
	t.Run("must ignore messages without metadata fields", func(t *testing.T) {
		ctx := context.Background()
		sender := &mockMsgSender{t: t}
		creator := "creator"
		testee := paloma.PalomaMessageSender{
			R:          &mockKeyRotator{},
			W:          sender,
			GetCreator: func() string { return creator },
		}

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
		creator := "creator"
		testee := paloma.PalomaMessageSender{
			R:          &mockKeyRotator{},
			W:          sender,
			GetCreator: func() string { return creator },
		}

		msg := &palomatypes.MsgAddStatusUpdate{
			Creator: "foo",
			Status:  "bar",
			Level:   palomatypes.MsgAddStatusUpdate_LEVEL_INFO,
		}

		sender.expectedMsg = msg

		_, err := testee.SendMsg(ctx, msg, "")
		require.NoError(t, err)
		require.Equal(t, creator, msg.GetMetadata().GetCreator(), "must inject creator")
		require.Equal(t, []string{"foobar"}, msg.GetMetadata().GetSigners(), "must inject creator")
	})
}
