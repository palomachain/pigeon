package cronchain

import (
	"context"
	"fmt"

	chain "github.com/volumefi/conductor/client"
	cronchain "github.com/volumefi/conductor/types/cronchain"
)

type Client struct {
	L *chain.LensClient
}

func (c Client) QueryMessagesForExecution(ctx context.Context) {
	qc := cronchain.NewQueryClient(c.L)
	msgs, err := qc.QueuedMessagesForSigning(ctx, &cronchain.QueryQueuedMessagesForSigningRequest{
		ValAddress:    "bob",
		QueueTypeName: "m",
	})
	for _, msg := range msgs.GetMsgs() {
		var m cronchain.QueuedSignedMessageI
		err := c.L.Codec.Marshaler.UnpackAny(msg, &m)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.GetId(), m.GetMsg())
	}

	fmt.Println(err)
}
