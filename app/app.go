package app

import (
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/palomachain/sparrow/attest"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/chain/evm"
	"github.com/palomachain/sparrow/chain/paloma"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/relayer"
	consensustypes "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	evmtypes "github.com/palomachain/sparrow/types/paloma/x/evm/types"
	valsettypes "github.com/palomachain/sparrow/types/paloma/x/valset/types"
	log "github.com/sirupsen/logrus"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/vizualni/whoops"
)

var (
	_relayer    *relayer.Relayer
	_config     *config.Root
	_configPath string

	_palomaClient *paloma.Client

	_evmClients    map[string]evm.Client
	_evmProcessors map[string]chain.Processor

	_attestRegistry *attest.Registry
)

func Relayer() *relayer.Relayer {
	if _relayer == nil {
		// do something
		_relayer = relayer.New(
			*Config(),
			*PalomaClient(),
			AttestRegistry(),
			GetEvmProcessors(),
		)
	}
	return _relayer
}

func SetConfigPath(path string) {
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

func GetEvmProcessors() map[string]chain.Processor {
	if _evmProcessors == nil {
		_evmProcessors = make(map[string]chain.Processor)
	}

	for chainID, client := range GetEvmClients() {
		_evmProcessors[chainID] = evm.NewProcessor(client, chainID)
	}

	return _evmProcessors
}

func GetEvmClients() map[string]evm.Client {
	if _evmClients == nil {
		_evmClients = make(map[string]evm.Client)
	}

	config := Config()
	for chainName, evmConfig := range config.EVM {
		if _, ok := _evmClients[chainName]; ok {
			log.WithFields(log.Fields{
				"chainName": chainName,
			}).Fatal("chain with chainName already registered")
		}

		_evmClients[chainName] = evm.NewClient(evmConfig, PalomaClient(), chainName)
	}

	return _evmClients
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

func AttestRegistry() *attest.Registry {
	if _attestRegistry == nil {
		_attestRegistry = attest.NewRegistry()
	}
	return _attestRegistry
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
					&consensustypes.MsgDeleteJob{},
				},
			},
		},
		MsgsImplementations: []byop.RegisterImplementation{
			{
				Iface: (*consensustypes.Message)(nil),
				Msgs: []proto.Message{
					&evmtypes.ArbitrarySmartContractCall{},
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
