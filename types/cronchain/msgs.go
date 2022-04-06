package cronchain

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	proto "github.com/gogo/protobuf/proto"
)

type QueuedSignedMessageI interface {
	proto.Message
	GetId() uint64
	GetMsg() *codectypes.Any
}
