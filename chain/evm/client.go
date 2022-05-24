package evm

import (
	"context"
	"embed"
	"math/big"
	"path/filepath"
	"strings"
	"sync"

	"io/fs"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/errors"
	"github.com/vizualni/whoops"
)

const (
	smartContractFilename = "simple"
)

/*
Do not delete hello.json contract. It's used for tests!
*/
var (
	//go:embed contracts/*.json
	contractsFS embed.FS

	readOnce   sync.Once
	_contracts = make(map[string]abi.ABI)
)

func StoredContracts() map[string]abi.ABI {
	readOnce.Do(func() {
		fs.WalkDir(contractsFS, ".", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			file, err := contractsFS.Open(path)

			contractName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			if err != nil {
				panic(err)
			}
			evmabi, err := abi.JSON(file)
			if err != nil {
				panic(err)
			}

			_contracts[contractName] = evmabi
			return nil
		})
	})
	return _contracts
}

type Client struct {
	config config.EVM

	smartContractAbi abi.ABI

	addr     ethcommon.Address
	keystore *keystore.KeyStore

	conn *ethclient.Client
}

func NewClient(cfg config.EVM) Client {
	client := &Client{
		config: cfg,
	}

	whoops.Assert(client.init())

	return *client
}

func (c *Client) init() error {
	return whoops.Try(func() {
		contracts := StoredContracts()
		scabi, ok := contracts[smartContractFilename]
		if !ok {
			whoops.Assert(errors.Unrecoverable(ErrSmartContractNotFound.Format(smartContractFilename)))
		}
		c.smartContractAbi = scabi

		if !ethcommon.IsHexAddress(c.config.SigningKey) {
			whoops.Assert(errors.Unrecoverable(ErrInvalidAddress.Format(c.config.SigningKey)))
		}
		c.addr = ethcommon.HexToAddress(c.config.SigningKey)

		c.keystore = keystore.NewKeyStore(c.config.KeyringDirectory.Path(), keystore.StandardScryptN, keystore.StandardScryptP)
		if !c.keystore.HasAddress(c.addr) {
			whoops.Assert(errors.Unrecoverable(ErrAddressNotFoundInKeyStore.Format(c.config.SigningKey, c.config.KeyringDirectory.Path())))
		}
		acc := accounts.Account{Address: c.addr}

		c.keystore.Unlock(acc, config.KeyringPassword(c.config.KeyringPassEnvName))

		c.conn = whoops.Must(ethclient.Dial(c.config.BaseRPCURL))
	})
}

//go:generate mockery --name=ethClienter --inpackage --testonly
type ethClienter interface {
	bind.ContractBackend
}

type executeSmartContractIn struct {
	ethClient ethClienter

	chainID       *big.Int
	gasAdjustment float64

	abi      abi.ABI
	contract common.Address

	signingAddr common.Address
	keystore    *keystore.KeyStore

	method    string
	arguments []any
}

func executeSmartContract(
	ctx context.Context,
	args executeSmartContractIn,
) error {
	return whoops.Try(func() {
		packedBytes := whoops.Must(args.abi.Pack(
			args.method,
			args.arguments...,
		))

		nonce := whoops.Must(
			args.ethClient.PendingNonceAt(ctx, args.signingAddr),
		)

		gasPrice := whoops.Must(
			args.ethClient.SuggestGasPrice(ctx),
		)

		// adjusting the gas price
		if args.gasAdjustment > 0.0 {
			gasAdj := big.NewFloat(args.gasAdjustment)
			gasAdj = gasAdj.Mul(gasAdj, new(big.Float).SetInt(gasPrice))
			gasPrice, _ = gasAdj.Int(big.NewInt(0))
		}

		boundContract := bind.NewBoundContract(
			args.contract,
			args.abi,
			args.ethClient,
			args.ethClient,
			args.ethClient,
		)

		txOpts := whoops.Must(
			bind.NewKeyStoreTransactorWithChainID(
				args.keystore,
				accounts.Account{Address: args.signingAddr},
				args.chainID,
			),
		)

		txOpts.Nonce = big.NewInt(int64(nonce))
		txOpts.GasPrice = gasPrice
		txOpts.From = args.signingAddr

		tx := whoops.Must(boundContract.RawTransact(txOpts, packedBytes))

		_ = tx

		// TODO: return tx hash and rest of the stuff

		return
	})
}

func (c Client) UpdateValset(ctx context.Context) {}

// TODO: this is just a placeholder
func (c Client) ExecuteArbitraryMessage(ctx context.Context) error {

	chainID := &big.Int{}
	chainID.SetString(c.config.ChainID, 10)

	return executeSmartContract(
		ctx,
		executeSmartContractIn{
			ethClient:     c.conn,
			chainID:       chainID,
			gasAdjustment: c.config.GasAdjustment,
			abi:           c.smartContractAbi,
			contract:      ethcommon.HexToAddress(c.config.EVMSpecificClientConfig.SmartContractAddress),
			signingAddr:   c.addr,
			keystore:      c.keystore,
			method:        "store",
			arguments: []any{
				big.NewInt(111),
			},
		},
	)
}
