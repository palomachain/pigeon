package app

import (
	"fmt"
	"os"
	"strings"

	lens "github.com/strangelove-ventures/lens/client"
	"github.com/vizualni/whoops"
	chain "github.com/volumefi/conductor/client"
	"github.com/volumefi/conductor/client/paloma"
	"github.com/volumefi/conductor/config"
	"github.com/volumefi/conductor/relayer"
)

var (
	_relayer      *relayer.Relayer
	_config       *config.Root
	_configPath   string
	_palomaClient *paloma.Client
)

func Relayer() *relayer.Relayer {
	if _relayer == nil {
		// do something
		_relayer = relayer.New(
			*Config(),
			*PalomaClient(),
		)
	}
	return _relayer
}

func SetConfigPath(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		panic("config must point to a file, not to a directory")
	}
	_configPath = path
}

func Config() *config.Root {
	if len(_configPath) == 0 {
		panic("config file path is not set")
	}
	if _config == nil {
		file, err := os.Open(_configPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		cnf, err := config.FromReader(file)
		if err != nil {
			panic(err)
		}
		_config = &cnf
	}

	return _config
}

func PalomaClient() *paloma.Client {
	if _palomaClient == nil {
		palomaConfig := Config().Paloma

		lensConfig := palomaLensClientConfig(palomaConfig.ChainClientConfig)
		lensConfig.KeyringBackend = palomaConfig.KeyringDetails.Type()

		// HACK: \n is added at the end of a password because github.com/cosmos/cosmos-sdk@v0.45.1/client/input/input.go at line 93 would return an EOF error which then would fail
		// Should be fixed with https://github.com/cosmos/cosmos-sdk/pull/11796
		passInput := strings.NewReader(palomaConfig.KeyringDetails.Password() + "\n")

		lensClient := whoops.Must(chain.NewChainClient(
			lensConfig,
			passInput,
			os.Stdout,
		))

		_palomaClient = &paloma.Client{
			L:          lensClient,
			GRPCClient: lensClient,
		}
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

func palomaLensClientConfig(palomaConfig config.ChainClientConfig) *lens.ChainClientConfig {
	fmt.Println(palomaConfig)
	return &lens.ChainClientConfig{
		ChainID:        defaultValue(palomaConfig.ChainID, "conductor"),
		RPCAddr:        defaultValue(palomaConfig.BaseRPCURL, "http://127.0.0.1:26657"),
		AccountPrefix:  defaultValue(palomaConfig.AccountPrefix, "paloma"),
		KeyringBackend: defaultValue(palomaConfig.KeyringDetails.Type(), "os"),
		GasAdjustment:  defaultValue(palomaConfig.GasAdjustment, 1.2),
		GasPrices:      defaultValue(palomaConfig.GasPrices, "0.01uatom"),
		KeyDirectory:   palomaConfig.KeyHomeDirectory,
		Debug:          false,
		Timeout:        defaultValue(palomaConfig.CallTimeout, "20s"),
		OutputFormat:   "json",
		SignModeStr:    "direct",
	}
}
