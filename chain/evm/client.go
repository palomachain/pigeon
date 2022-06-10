package evm

import (
	"bytes"
	"context"
	"embed"
	"math/big"
	"path/filepath"
	"strings"
	"sync"

	"io/fs"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/errors"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

const (
	smartContractFilename = "simple"
)

type StoredContract struct {
	ABI    abi.ABI
	Source []byte
}

/*
Do not delete hello.json contract. It's used for tests!
*/
var (
	//go:embed contracts/*.json
	contractsFS embed.FS

	readOnce   sync.Once
	_contracts = make(map[string]StoredContract)
)

func StoredContracts() map[string]StoredContract {
	readOnce.Do(func() {
		fs.WalkDir(contractsFS, ".", func(path string, d fs.DirEntry, err error) error {
			logger := log.WithFields(log.Fields{
				"path": path,
			})
			if d.IsDir() {
				return nil
			}
			file, err := contractsFS.Open(path)
			if err != nil {
				logger.WithFields(log.Fields{
					"err": err,
				}).Fatal("couldn't open contract file")
			}

			contractName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

			// we need to store body locally, so reading it here first and
			// using bytes.NewBuffer few lines down.
			body := whoops.Must(ioutil.ReadAll(file))

			evmabi, err := abi.JSON(bytes.NewBuffer(body))
			if err != nil {
				logger.WithFields(log.Fields{
					"err": err,
				}).Fatal("couldn't read contract file")
			}

			_contracts[contractName] = StoredContract{
				ABI:    evmabi,
				Source: body,
			}
			return nil
		})
	})
	return _contracts
}

type PalomaClienter interface {
	DeleteJob(ctx context.Context, queueTypeName string, id uint64) error
}

type Client struct {
	config config.EVM

	smartContractAbi abi.ABI

	addr     ethcommon.Address
	keystore *keystore.KeyStore

	conn *ethclient.Client

	paloma PalomaClienter
}

func NewClient(cfg config.EVM, palomaClient PalomaClienter) Client {
	client := &Client{
		config: cfg,
		paloma: palomaClient,
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
		c.smartContractAbi = scabi.ABI

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
	abiBytes []byte,
	packedBytes []byte,
) error {
	logger := log.WithFields(log.Fields{
		"chain-id":        args.chainID,
		"arguments":       args.arguments,
		"contract-addr":   args.contract,
		"gas-adjustments": args.gasAdjustment,
		"method":          args.method,
		"signing-addr":    args.signingAddr,
	})
	return whoops.Try(func() {
		// TODO
		// packedBytes := whoops.Must(args.abi.Pack(
		// 	args.method,
		// 	args.arguments...,
		// ))

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
			logger.WithFields(log.Fields{
				"gas-price": gasPrice,
			}).Info("adusted gas price")
		}

		aabi := whoops.Must(abi.JSON(bytes.NewBuffer(abiBytes)))

		boundContract := bind.NewBoundContract(
			args.contract,
			// TODO:
			// args.abi,
			aabi,
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

		logger = logger.WithFields(log.Fields{
			"tx-opts": txOpts,
		})
		logger.Info("executing tx")

		tx := whoops.Must(boundContract.RawTransact(txOpts, packedBytes))
		logger.WithFields(log.Fields{
			"tx-hash":      tx.Hash(),
			"tx-gas-limit": tx.Gas(),
			"tx-gas-price": tx.GasPrice(),
			"tx-cost":      tx.Cost(),
		}).Info("tx executed")

		_ = tx

		// TODO: return tx hash and rest of the stuff

		return
	})
}

func (c Client) sign(ctx context.Context, bytes []byte) ([]byte, error) {
	return c.keystore.SignHash(accounts.Account{Address: c.addr}, bytes)
}

// TODO: this is just a placeholder
func (c Client) executeArbitraryMessage(
	ctx context.Context,
	// TODO
	msg *types.ArbitrarySmartContractCall,
) error {
	chainID := &big.Int{}
	chainID.SetString(c.config.ChainID, 10)

	return executeSmartContract(
		ctx,
		executeSmartContractIn{
			ethClient:     c.conn,
			chainID:       chainID,
			gasAdjustment: c.config.GasAdjustment,
			abi:           c.smartContractAbi,
			// contract:      ethcommon.HexToAddress(c.config.EVMSpecificClientConfig.SmartContractAddress),
			contract:    ethcommon.HexToAddress(msg.GetHexAddress()),
			signingAddr: c.addr,
			keystore:    c.keystore,
			method:      "store",
			arguments: []any{
				big.NewInt(111),
			},
		},
		msg.GetAbi(),
		msg.GetPayload(),
	)
}
