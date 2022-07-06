package types

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

func (msg *MsgDeleteJob) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteJob) ValidateBasic() error {
	return nil
}

func (msg *MsgAddEvidence) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddEvidence) ValidateBasic() error {
	return nil
}

func (msg *MsgSetPublicAccessData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetPublicAccessData) ValidateBasic() error {
	return nil
}
