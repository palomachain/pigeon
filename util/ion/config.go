package ion

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
)

var ModuleBasics = []module.AppModuleBasic{
	auth.AppModuleBasic{},
	authz.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	// TODO: add osmosis governance proposal types here
	// TODO: add other proposal types here
	gov.NewAppModuleBasic(
		[]client.ProposalHandler{
			paramsclient.ProposalHandler,
			upgradeclient.LegacyProposalHandler,
			upgradeclient.LegacyCancelProposalHandler,
		},
	),
	crisis.AppModuleBasic{},
	distribution.AppModuleBasic{},
	feegrant.AppModuleBasic{},
	mint.AppModuleBasic{},
	params.AppModuleBasic{},
	slashing.AppModuleBasic{},
	staking.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	transfer.AppModuleBasic{},
	ibc.AppModuleBasic{},
}

type ChainClientConfig struct {
	BlockTimeout   string                  `json:"block-timeout" yaml:"block-timeout"`
	OutputFormat   string                  `json:"output-format" yaml:"output-format"`
	KeyDirectory   string                  `json:"key-directory" yaml:"key-directory"`
	GRPCAddr       string                  `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix  string                  `json:"account-prefix" yaml:"account-prefix"`
	KeyringBackend string                  `json:"keyring-backend" yaml:"keyring-backend"`
	SignModeStr    string                  `json:"sign-mode" yaml:"sign-mode"`
	GasPrices      string                  `json:"gas-prices" yaml:"gas-prices"`
	RPCAddr        string                  `json:"rpc-addr" yaml:"rpc-addr"`
	ChainID        string                  `json:"chain-id" yaml:"chain-id"`
	Timeout        string                  `json:"timeout" yaml:"timeout"`
	Key            string                  `json:"key" yaml:"key"`
	ExtraCodecs    []string                `json:"extra-codecs" yaml:"extra-codecs"`
	Modules        []module.AppModuleBasic `json:"-" yaml:"-"`
	Slip44         int                     `json:"slip44" yaml:"slip44"`
	GasAdjustment  float64                 `json:"gas-adjustment" yaml:"gas-adjustment"`
	MinGasAmount   uint64                  `json:"min-gas-amount" yaml:"min-gas-amount"`
	Debug          bool                    `json:"debug" yaml:"debug"`
}
