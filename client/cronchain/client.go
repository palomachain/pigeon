package cronchain

import (
	"context"
	"fmt"

	chain "github.com/volumefi/conductor/client"
	"github.com/volumefi/conductor/client/cronchain/types"
)

type Client struct {
	L *chain.LensClient
}

func (c Client) QueryMessagesForExecution(ctx context.Context) {
	qc := types.NewQueryClient(c.L)
	msgs, err := qc.QueuedMessagesForSigning(ctx, &types.QueryQueuedMessagesForSigningRequest{
		ValAddress:    "bob",
		QueueTypeName: "m",
	})
	for _, msg := range msgs.GetMsgs() {
		var m types.QueuedSignedMessageI
		err := c.L.Codec.Marshaler.UnpackAny(msg, &m)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.GetId(), m.GetMsg())
	}

	fmt.Println(err)
}
