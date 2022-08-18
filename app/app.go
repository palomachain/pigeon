package app

import (
	"os"
	"strings"
	gotime "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/health"
	"github.com/palomachain/pigeon/relayer"
	consensustypes "github.com/palomachain/pigeon/types/paloma/x/consensus/types"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	valsettypes "github.com/palomachain/pigeon/types/paloma/x/valset/types"
	"github.com/palomachain/pigeon/util/time"
	log "github.com/sirupsen/logrus"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/vizualni/whoops"
)

const (
	AppName     = "pigeon"
	AppNameCaps = "PIGEON"
)

var (
	_relayer    *relayer.Relayer
	_config     *config.Root
	_configPath string

	_palomaClient *paloma.Client

	_evmFactory *evm.Factory

	_timeAdapter time.Time

	_healthCheckService *health.Service
)

var (
	version = ""
	commit  = ""
)

func Version() string { return version }

func Commit() string { return commit }

func Relayer() *relayer.Relayer {
	if _relayer == nil {
		// do something
		_relayer = relayer.New(
			*Config(),
			*PalomaClient(),
			EvmFactory(),
			Time(),
			relayer.Config{
				KeepAliveLoopTimeout: 30 * gotime.Second,
				KeepAliveThreshold:   1 * gotime.Minute,
			},
		)
	}
	return _relayer
}

func SetConfigPath(path string) {
	path = config.Filepath(path).Path()
	fi, err := os.Stat(path)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("couldn't stat config file")
	}
	if fi.IsDir() {
		log.WithFields(log.Fields{
			"path": path,
		}).Fatal("path must be a file, not a directory")
	}
	_configPath = path
}

func EvmFactory() *evm.Factory {
	if _evmFactory == nil {
		_evmFactory = evm.NewFactory(PalomaClient())
	}

	return _evmFactory
}

func Config() *config.Root {
	if len(_configPath) == 0 {
		log.Fatal("config file path is not set")
	}
	if _config == nil {
		file, err := os.Open(_configPath)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("couldn't open config file")
		}
		defer file.Close()
		cnf, err := config.FromReader(file)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("couldn't read config file")
		}
		_config = &cnf
	}

	return _config
}

func PalomaClient() *paloma.Client {
	if _palomaClient == nil {
		palomaConfig := Config().Paloma

		lensConfig := palomaLensClientConfig(palomaConfig)

		// HACK: \n is added at the end of a password because github.com/cosmos/cosmos-sdk@v0.45.1/client/input/input.go at line 93 would return an EOF error which then would fail
		// Should be fixed with https://github.com/cosmos/cosmos-sdk/pull/11796
		passInput := strings.NewReader(config.KeyringPassword(palomaConfig.KeyringPassEnvName) + "\n")

		lensClient := whoops.Must(chain.NewChainClient(
			lensConfig,
			passInput,
			os.Stdout,
		))

		_palomaClient = &paloma.Client{
			L:             lensClient,
			GRPCClient:    lensClient,
			MessageSender: lensClient,
			PalomaConfig:  palomaConfig,
		}
		_palomaClient.Init()
	}
	return _palomaClient
}

func defaultValue[T comparable](proposedVal T, defaultVal T) T {
	var zero T
	if proposedVal == zero {
		return defaultVal
	}
	return proposedVal
}

func palomaLensClientConfig(palomaConfig config.Paloma) *lens.ChainClientConfig {
	modules := lens.ModuleBasics[:]

	modules = append(modules, byop.Module{
		ModuleName: "paloma",
		MsgsInterfaces: []byop.RegisterInterface{
			{
				Name:  "paloma",
				Iface: (*sdk.Msg)(nil),
				Msgs: []proto.Message{
					&consensustypes.MsgAddMessagesSignatures{},
					&valsettypes.MsgAddExternalChainInfoForValidator{},
					&valsettypes.MsgKeepAlive{},
					&consensustypes.MsgDeleteJob{},
					&consensustypes.MsgAddEvidence{},
					&consensustypes.MsgSetPublicAccessData{},
				},
			},
			{
				Name:  "any-messages",
				Iface: (*proto.Message)(nil),
				Msgs: []proto.Message{
					&evmtypes.TxExecutedProof{},
					&evmtypes.SmartContractExecutionErrorProof{},
					&evmtypes.ValidatorBalancesAttestation{},
					&evmtypes.ValidatorBalancesAttestationRes{},
				},
			},
		},
		MsgsImplementations: []byop.RegisterImplementation{
			{
				Iface: (*consensustypes.Message)(nil),
				Msgs: []proto.Message{
					&evmtypes.ArbitrarySmartContractCall{},
					&evmtypes.Message{},
					&evmtypes.ValidatorBalancesAttestation{},
					&evmtypes.ValidatorBalancesAttestationRes{},
				},
			},
		},
	})

	return &lens.ChainClientConfig{
		Key:            palomaConfig.SigningKey,
		ChainID:        defaultValue(palomaConfig.ChainID, "paloma"),
		RPCAddr:        defaultValue(palomaConfig.BaseRPCURL, "http://127.0.0.1:26657"),
		AccountPrefix:  defaultValue(palomaConfig.AccountPrefix, "paloma"),
		KeyringBackend: defaultValue(palomaConfig.KeyringType, "os"),
		GasAdjustment:  defaultValue(palomaConfig.GasAdjustment, 1.2),
		GasPrices:      defaultValue(palomaConfig.GasPrices, "0.01uatom"),
		KeyDirectory:   palomaConfig.KeyringDirectory.Path(),
		Debug:          false,
		Timeout:        defaultValue(palomaConfig.CallTimeout, "20s"),
		OutputFormat:   "json",
		SignModeStr:    "direct",
		Modules:        modules,
	}
}

func Time() time.Time {
	if _timeAdapter == nil {
		_timeAdapter = time.New()
	}
	return _timeAdapter
}

func HealthCheckService() health.Service {
	if _healthCheckService == nil {
		_healthCheckService = &health.Service{
			Checks: []health.Checker{
				Relayer(),
			},
		}
	}
	return *_healthCheckService
}
