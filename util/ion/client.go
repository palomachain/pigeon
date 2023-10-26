package ion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/avast/retry-go/v4"
	provtypes "github.com/cometbft/cometbft/light/provider"
	prov "github.com/cometbft/cometbft/light/provider/http"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"gopkg.in/yaml.v3"
)

var (
	// Variables used for retries
	RtyAttNum = uint(5)
	RtyAtt    = retry.Attempts(RtyAttNum)
	RtyDel    = retry.Delay(time.Millisecond * 400)
	RtyErr    = retry.LastErrorOnly(true)
)

type Client struct {
	Codec          Codec
	Keybase        keyring.Keyring
	RPCClient      rpcclient.Client
	LightProvider  provtypes.Provider
	Input          io.Reader
	Output         io.Writer
	Config         *ChainClientConfig
	KeyringOptions []keyring.Option
}

// UnpackAny implements types.AnyUnpacker.
func (c *Client) UnpackAny(any *types.Any, iface interface{}) error {
	return c.Codec.Marshaler.UnpackAny(any, iface)
}

func (c *Client) GetKeybase() keyring.Keyring {
	return c.Keybase
}

func NewClient(cfg *ChainClientConfig, input io.Reader, output io.Writer, kro ...keyring.Option) (*Client, error) {
	c := &Client{
		KeyringOptions: kro,
		Config:         cfg,
		Input:          input,
		Output:         output,
		Codec:          makeCodec(cfg.Modules, nil),
	}

	// TODO: Needed?
	// cfg.KeyDirectory = keysDir("", cfg.Key)

	return c.init()
}

func (c *Client) init() (*Client, error) {
	// TODO: test key directory and return error if not created
	keybase, err := keyring.New(c.Config.ChainID, c.Config.KeyringBackend, c.Config.KeyDirectory, c.Input, c.Codec.Marshaler, c.KeyringOptions...)
	if err != nil {
		return nil, err
	}
	// TODO: figure out how to deal with input or maybe just make all keyring backends test?

	timeout, _ := time.ParseDuration(c.Config.Timeout)
	rpcClient, err := NewRPCClient(c.Config.RPCAddr, timeout)
	if err != nil {
		return nil, err
	}

	lightprovider, err := prov.New(c.Config.ChainID, c.Config.RPCAddr)
	if err != nil {
		return nil, err
	}

	c.RPCClient = rpcClient
	c.LightProvider = lightprovider
	c.Keybase = keybase

	return c, nil
}

func (cc *Client) GetKeyAddress() (sdk.AccAddress, error) {
	done := cc.SetSDKContext()
	defer done()

	info, err := cc.Keybase.Key(cc.Config.Key)
	if err != nil {
		return nil, err
	}
	return info.GetAddress()
}

func NewRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

// AccountFromKeyOrAddress returns an account from either a key or an address
// if empty string is passed in this returns the default key's address
func (cc *Client) AccountFromKeyOrAddress(keyOrAddress string) (out sdk.AccAddress, err error) {
	switch {
	case keyOrAddress == "":
		out, err = cc.GetKeyAddress()
	case cc.KeyExists(keyOrAddress):
		cc.Config.Key = keyOrAddress
		out, err = cc.GetKeyAddress()
	default:
		out, err = cc.DecodeBech32AccAddr(keyOrAddress)
	}
	return
}

func keysDir(home, chainID string) string {
	return path.Join(home, "keys", chainID)
}

// TODO: actually do something different here have a couple of levels of verbosity
func (cc *Client) PrintTxResponse(res *sdk.TxResponse) error {
	return cc.PrintObject(res)
}

func (cc *Client) HandleAndPrintMsgSend(res *sdk.TxResponse, err error) error {
	if err != nil {
		if res != nil {
			return fmt.Errorf("failed to withdraw rewards: code(%d) msg(%s)", res.Code, res.Logs)
		}
		return fmt.Errorf("failed to withdraw rewards: err(%w)", err)
	}
	return cc.PrintTxResponse(res)
}

func (cc *Client) PrintObject(res interface{}) error {
	var (
		bz  []byte
		err error
	)
	switch cc.Config.OutputFormat {
	case "json":
		if m, ok := res.(proto.Message); ok {
			bz, err = cc.MarshalProto(m)
		} else {
			bz, err = json.Marshal(res)
		}
		if err != nil {
			return err
		}
	case "indent":
		if m, ok := res.(proto.Message); ok {
			bz, err = cc.MarshalProto(m)
			if err != nil {
				return err
			}
			buf := bytes.NewBuffer([]byte{})
			if err = json.Indent(buf, bz, "", "  "); err != nil {
				return err
			}
			bz = buf.Bytes()
		} else {
			bz, err = json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
		}
	case "yaml":
		bz, err = yaml.Marshal(res)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output type: %s", cc.Config.OutputFormat)
	}
	fmt.Fprint(cc.Output, string(bz), "\n")
	return nil
}

func (cc *Client) MarshalProto(res proto.Message) ([]byte, error) {
	return cc.Codec.Marshaler.MarshalJSON(res)
}
