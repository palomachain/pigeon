package cronchain

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
)

type QueuedSignedMessageI interface {
	proto.Message
	GetId() uint64
	GetMsg() *codectypes.Any
}

func (msg *MsgAddMessagesSignatures) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddMessagesSignatures) ValidateBasic() error {
	return nil
}
