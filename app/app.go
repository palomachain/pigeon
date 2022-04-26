package app

import (
	"os"

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

func Config() *config.Root {
	if _config == nil {
		_config = &config.Root{}
	}

	return _config
}

func PalomaClient() *paloma.Client {
	if _palomaClient == nil {
		_palomaClient = &paloma.Client{
			L: whoops.Must(chain.NewChainClient(
				palomaLensClientConfig("", false),
				os.Stdin,
				os.Stdout,
			)),
		}
	}
	return _palomaClient
}

// TODO: take real values and take things from _Config()_
func palomaLensClientConfig(keyHome string, debug bool) *lens.ChainClientConfig {
	return &lens.ChainClientConfig{
		Key:            "default",
		ChainID:        "cosmoshub-4",
		RPCAddr:        "https://cosmoshub-4.technofractal.com:443",
		GRPCAddr:       "https://gprc.cosmoshub-4.technofractal.com:443",
		AccountPrefix:  "cosmos",
		KeyringBackend: "test",
		GasAdjustment:  1.2,
		GasPrices:      "0.01uatom",
		KeyDirectory:   keyHome,
		Debug:          debug,
		Timeout:        "20s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
	}
}
