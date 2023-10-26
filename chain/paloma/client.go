package paloma

import (
	"context"

	"github.com/VolumeFi/whoops"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	consensus "github.com/palomachain/paloma/x/consensus/types"
	gravity "github.com/palomachain/paloma/x/gravity/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	valset "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

type (
	ResultStatus = coretypes.ResultStatus
	Unpacker     = codectypes.AnyUnpacker
)

//go:generate mockery --name=MessageSender
type MessageSender interface {
	SendMsg(ctx context.Context, msg sdk.Msg, memo string) (*sdk.TxResponse, error)
}

//go:generate mockery --name=IonClient
type IonClient interface {
	Status(context.Context) (*ResultStatus, error)
	DecodeBech32ValAddr(string) (sdk.ValAddress, error)
	GetKeybase() keyring.Keyring
	SetSDKContext() func()
}

type Client struct {
	PalomaConfig config.Paloma // This is only needed ONCE ! Can we remove it?
	GRPCClient   grpc.ClientConn

	// TODO: Can this shit not be made private???
	Ion           IonClient
	Unpacker      Unpacker
	MessageSender MessageSender

	creator        string
	creatorValoper string
	valAddr        sdk.ValAddress
}

func NewClient(cfg config.Paloma, grpcWrapper grpc.ClientConn, ion IonClient, sender MessageSender, unpacker codectypes.AnyUnpacker) *Client {
	return (&Client{
		PalomaConfig:  cfg,
		GRPCClient:    grpcWrapper,
		Ion:           ion,
		Unpacker:      unpacker,
		MessageSender: sender,
	}).init()
}

func (c *Client) init() *Client {
	log.Info("paloma client: init")
	c.creator = getCreator(c)
	log.Info("paloma client: get as valoper")
	c.creatorValoper = getCreatorAsValoper(c)
	log.Info("paloma client: get val address")
	c.valAddr = sdk.ValAddress(getMainAddress(c).Bytes())
	return c
}

type BroadcastMessageSignatureIn struct {
	ID              uint64
	QueueTypeName   string
	Signature       []byte
	SignedByAddress string
}

// BroadcastMessageSignatures takes a list of signatures that need to be sent over to the chain.
// It build the message and sends it over.
func (c *Client) BroadcastMessageSignatures(ctx context.Context, signatures ...BroadcastMessageSignatureIn) error {
	return broadcastMessageSignatures(ctx, c.MessageSender, c.creator, signatures...)
}

func (c *Client) BlockHeight(ctx context.Context) (int64, error) {
	res, err := c.Ion.Status(ctx)
	if err != nil {
		return 0, err
	}

	if res.SyncInfo.CatchingUp {
		return 0, ErrNodeIsNotInSync
	}

	return res.SyncInfo.LatestBlockHeight, nil
}

// TODO Combine with below method
func (c *Client) SendBatchSendToEVMClaim(ctx context.Context, claim gravity.MsgBatchSendToEthClaim) error {
	_, err := c.MessageSender.SendMsg(ctx, &claim, "")
	return err
}

func (c *Client) SendSendToPalomaClaim(ctx context.Context, claim gravity.MsgSendToPalomaClaim) error {
	_, err := c.MessageSender.SendMsg(ctx, &claim, "")
	return err
}

type ChainInfoIn struct {
	ChainType        string
	ChainReferenceID string
	AccAddress       string
	PubKey           []byte
	Traits           []string
}

// AddExternalChainInfo adds info about the external chain. It adds the chain's
// account addresses that the pigeon knows about.
func (c *Client) AddExternalChainInfo(ctx context.Context, chainInfos ...ChainInfoIn) error {
	if len(chainInfos) == 0 {
		return nil
	}

	msg := &valset.MsgAddExternalChainInfoForValidator{
		Creator: c.creator,
	}

	msg.ChainInfos = slice.Map(
		chainInfos,
		func(in ChainInfoIn) *valset.ExternalChainInfo {
			return &valset.ExternalChainInfo{
				ChainType:        in.ChainType,
				ChainReferenceID: in.ChainReferenceID,
				Address:          in.AccAddress,
				Pubkey:           in.PubKey,
				Traits:           in.Traits,
			}
		},
	)

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error {
	anyProof, err := codectypes.NewAnyWithValue(proof)
	if err != nil {
		return err
	}
	msg := &consensus.MsgAddEvidence{
		Creator:       c.creator,
		Proof:         anyProof,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err = c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	msg := &consensus.MsgSetPublicAccessData{
		Creator:       c.creator,
		Data:          data,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) SetErrorData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	msg := &consensus.MsgSetErrorData{
		Creator:       c.creator,
		Data:          data,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) AddStatusUpdate(ctx context.Context, level palomatypes.MsgAddStatusUpdate_Level, status string) error {
	msg := &palomatypes.MsgAddStatusUpdate{
		Creator: c.creator,
		Status:  status,
		Level:   level,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) KeepValidatorAlive(ctx context.Context, appVersion string) error {
	msg := &valset.MsgKeepAlive{
		Creator:       c.creator,
		PigeonVersion: appVersion,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) Status(ctx context.Context) (*ResultStatus, error) {
	liblog.WithContext(ctx).Info("STATUS")
	res, err := c.Ion.Status(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) PalomaStatus(ctx context.Context) error {
	liblog.WithContext(ctx).Info("PALOMA STATUS")
	res, err := c.Status(ctx)

	if IsPalomaDown(err) {
		return whoops.Wrap(ErrPalomaIsDown, err)
	}

	if err != nil {
		return err
	}
	_ = res
	return nil
}

func (c *Client) GetValidator(ctx context.Context) (*stakingtypes.Validator, error) {
	address := c.GetValidatorAddress().String()
	_, err := c.Ion.DecodeBech32ValAddr(address)
	if err != nil {
		return nil, err
	}
	queryClient := stakingtypes.NewQueryClient(c.GRPCClient)
	req := &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: address,
	}
	res, err := queryClient.Validator(ctx, req)
	if err != nil {
		return nil, err
	}
	return &res.Validator, nil
}

func broadcastMessageSignatures(
	ctx context.Context,
	ms MessageSender,
	creator string,
	signatures ...BroadcastMessageSignatureIn,
) error {
	if len(signatures) == 0 {
		return nil
	}
	signedMessages := make([]*consensus.ConsensusMessageSignature, len(signatures))
	for i, sig := range signatures {
		signedMessages[i] = &consensus.ConsensusMessageSignature{
			Id:              sig.ID,
			QueueTypeName:   sig.QueueTypeName,
			Signature:       sig.Signature,
			SignedByAddress: sig.SignedByAddress,
		}
	}
	msg := &consensus.MsgAddMessagesSignatures{
		Creator:        creator,
		SignedMessages: signedMessages,
	}
	_, err := ms.SendMsg(ctx, msg, "")
	return err
}

func (c *Client) Keyring() keyring.Keyring {
	return c.Ion.GetKeybase()
}

func (c *Client) GetValidatorAddress() sdk.ValAddress {
	return c.valAddr
}

func (c *Client) GetCreator() string {
	return c.creator
}

func getMainAddress(c *Client) sdk.Address {
	key, err := c.Keyring().Key(c.PalomaConfig.SigningKey)
	if err != nil {
		panic(err)
	}
	address, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	return address
}

// TODO: CLEAN THIS UP???
func getCreator(c *Client) string {
	return c.addressString(getMainAddress(c))
}

func getCreatorAsValoper(c *Client) string {
	return c.addressString(sdk.ValAddress(getMainAddress(c).Bytes()))
}

func (c Client) addressString(val sdk.Address) string {
	defer (c.Ion.SetSDKContext())()
	return val.String()
}
