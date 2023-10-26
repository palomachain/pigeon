package app

import (
	"fmt"
	"os"
	"strings"
	gotime "time"

	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	consensustypes "github.com/palomachain/paloma/x/consensus/types"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	gravitytypes "github.com/palomachain/paloma/x/gravity/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	valsettypes "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain/evm"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/health"
	"github.com/palomachain/pigeon/relayer"
	"github.com/palomachain/pigeon/util/ion"
	"github.com/palomachain/pigeon/util/rotator"
	"github.com/palomachain/pigeon/util/time"
	log "github.com/sirupsen/logrus"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
)

const (
	AppName     = "pigeon"
	AppNameCaps = "PIGEON"
)

var (
	_relayer    *relayer.Relayer
	_config     *config.Config
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

func Version() string {
	if !strings.HasPrefix(version, "v") {
		version = fmt.Sprintf("v%s", version)
	}

	return version
}

func Commit() string { return commit }

func Relayer() *relayer.Relayer {
	if _relayer == nil {
		_relayer = relayer.New(
			Config(),
			PalomaClient(),
			EvmFactory(),
			Time(),
			relayer.Config{
				KeepAliveLoopTimeout:    5 * gotime.Second,
				KeepAliveBlockThreshold: 600, // Approximately 15 minutes at 1.62 blocks per second
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

func Config() *config.Config {
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

		if len(cnf.Paloma.SigningKeys) < 1 {
			log.Info("No signing key collection provided, falling back to legacy signing key")
			cnf.Paloma.SigningKeys = []string{cnf.Paloma.SigningKey}
		}

		_config = cnf
	}

	return _config
}

func PalomaClient() *paloma.Client {
	if _palomaClient == nil {
		log.Info("cfg loading...")
		palomaConfig := Config().Paloma

		log.Info("client cfg loading...")
		clientCfg := palomaClientConfig(palomaConfig)

		// HACK: \n is added at the end of a password because github.com/cosmos/cosmos-sdk@v0.45.1/client/input/input.go at line 93 would return an EOF error which then would fail
		// Should be fixed with https://github.com/cosmos/cosmos-sdk/pull/11796
		passInput := strings.NewReader(config.KeyringPassword(palomaConfig.KeyringPassEnvName) + "\n")
		log.Info("ion client construction...")
		ionClient := whoops.Must(ion.NewClient(
			clientCfg,
			passInput,
			os.Stdout,
		))

		// It's important to pass the lens config as pointer throughout the codebase for this to work!
		fn := func(s string) {
			log.Info("new lens key", "lens-key", s)
			clientCfg.Key = s
		}
		r := rotator.New(fn, palomaConfig.SigningKeys...)

		grpcWrapper := paloma.GRPCClientWrapper{W: ionClient}
		senderWrapper := paloma.PalomaMessageSender{W: ionClient, R: r}
		log.Info("new client...")
		_palomaClient = paloma.NewClient(palomaConfig, grpcWrapper, ionClient, senderWrapper, ionClient)
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

func palomaClientConfig(palomaConfig config.Paloma) *ion.ChainClientConfig {
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
					&consensustypes.MsgAddEvidence{},
					&consensustypes.MsgSetPublicAccessData{},
					&consensustypes.MsgSetErrorData{},
					&palomatypes.MsgAddStatusUpdate{},
					&gravitytypes.MsgSendToEth{},
					&gravitytypes.MsgConfirmBatch{},
					&gravitytypes.MsgSendToPalomaClaim{},
					&gravitytypes.MsgBatchSendToEthClaim{},
					&gravitytypes.MsgCancelSendToEth{},
					&gravitytypes.MsgSubmitBadSignatureEvidence{},
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
					&evmtypes.TransferERC20Ownership{},
				},
			},
		},
		MsgsImplementations: []byop.RegisterImplementation{
			{
				Iface: (*consensustypes.ConsensusMsg)(nil),
				Msgs: []proto.Message{
					&evmtypes.SubmitLogicCall{},
					&evmtypes.Message{},
					&evmtypes.ValidatorBalancesAttestation{},
					&evmtypes.ValidatorBalancesAttestationRes{},
				},
			},
		},
	})

	return &ion.ChainClientConfig{
		// TODO: FIX FIX FIX
		// Key:            palomaConfig.SigningKeys[0],
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
