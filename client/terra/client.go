package terra

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	chain "github.com/palomachain/sparrow/client"
	types "github.com/palomachain/sparrow/types/terra"
)

type TurnstoneConsensusMessage struct {
	MessageID string
	Nonce     []byte
	Payload   json.Marshaler
}

type SmartContractExecution struct {
	Contract string
	Sender   string
	Payload  json.RawMessage
	Coins    sdk.Coins
}

type Client struct {
	LensClient *chain.LensClient

	// TODO
	smartContractAddr string
}

// TODO: this is currently oversimplified. Once we start using this for real we will adapt.
// Maybe better thing would be to actually use the "Invoke" method along with the grpc client.
func (c Client) ExecuteTurnstoneConsensusMessage(
	ctx context.Context,
	messageId string,
	nonce []byte,
	payload []byte,
) (*sdk.TxResponse, error) {
	// TODO: do validations
	var msgExec types.MsgExecuteContract

	return c.LensClient.SendMsg(ctx, &msgExec)
}
