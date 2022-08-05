package evm

import (
	"bytes"
	"context"
	"embed"
	"math/big"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"io/fs"
	"io/ioutil"

	"github.com/ethereum/go-ethereum"
	etherum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	etherumtypes "github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	proto "github.com/gogo/protobuf/proto"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/types/paloma/x/evm/types"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

const (
	smartContractFilename = "compass-evm"
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

//go:generate mockery --name=PalomaClienter
type PalomaClienter interface {
	AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error
	SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error
	QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*types.Valset, error)
}

type Client struct {
	config config.EVM

	addr     ethcommon.Address
	keystore *keystore.KeyStore

	conn *ethclient.Client

	paloma PalomaClienter
}

func (c *Client) init() error {
	return whoops.Try(func() {

		if !ethcommon.IsHexAddress(c.config.SigningKey) {
			whoops.Assert(errors.Unrecoverable(ErrInvalidAddress.Format(c.config.SigningKey)))
		}
		c.addr = ethcommon.HexToAddress(c.config.SigningKey)

		if c.keystore == nil {
			c.keystore = keystore.NewKeyStore(c.config.KeyringDirectory.Path(), keystore.StandardScryptN, keystore.StandardScryptP)
		}

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

func callSmartContract(
	ctx context.Context,
	args executeSmartContractIn,
) (*etherumtypes.Transaction, error) {
	logger := log.WithFields(log.Fields{
		"chain-id":      args.chainID,
		"contract-addr": args.contract,

		"method":    args.method,
		"arguments": args.arguments,

		"gas-adjustments": args.gasAdjustment,

		"signing-addr": args.signingAddr,
	})
	return whoops.TryVal(func() *etherumtypes.Transaction {
		packedBytes, err := args.abi.Pack(
			args.method,
			args.arguments...,
		)
		whoops.Assert(err)

		nonce, err := args.ethClient.PendingNonceAt(ctx, args.signingAddr)
		whoops.Assert(err)

		gasPrice, err := args.ethClient.SuggestGasPrice(ctx)
		whoops.Assert(err)

		logger.WithField("suggested-gas-price", gasPrice).Debug("suggested gas price")

		// adjusting the gas price
		if args.gasAdjustment > 0.0 {
			gasAdj := big.NewFloat(args.gasAdjustment)
			gasAdj = gasAdj.Mul(gasAdj, new(big.Float).SetInt(gasPrice))
			gasPrice, _ = gasAdj.Int(big.NewInt(0))
			logger.WithFields(log.Fields{
				"gas-price": gasPrice,
			}).Debug("adjusted gas price")
		}

		boundContract := bind.NewBoundContract(
			args.contract,
			args.abi,
			args.ethClient,
			args.ethClient,
			args.ethClient,
		)

		txOpts, err := bind.NewKeyStoreTransactorWithChainID(
			args.keystore,
			accounts.Account{Address: args.signingAddr},
			args.chainID,
		)
		whoops.Assert(err)

		txOpts.Nonce = big.NewInt(int64(nonce))
		txOpts.GasPrice = gasPrice
		txOpts.From = args.signingAddr

		logger = logger.WithFields(log.Fields{
			"gas-limit": txOpts.GasLimit,
			"gas-price": txOpts.GasPrice,
			"nonce":     txOpts.Nonce,
			"from":      txOpts.From,
		})

		logger.Debug("executing tx")

		tx, err := boundContract.RawTransact(txOpts, packedBytes)
		whoops.Assert(err)

		logger.WithFields(log.Fields{
			"tx-hash":      tx.Hash(),
			"tx-gas-limit": tx.Gas(),
			"tx-gas-price": tx.GasPrice(),
			"tx-cost":      tx.Cost(),
		}).Debug("tx executed")
		return tx
	})
}

func (c Client) sign(ctx context.Context, bytes []byte) ([]byte, error) {
	return c.keystore.SignHash(
		accounts.Account{Address: c.addr},
		bytes,
	)
}

// FilterLogs will gather all logs given a FilterQuery. If it encounters an
// error saying that there are too many results in the provided block window,
// then it's going to try to do this using a "binary search" approach while
// splitting the  possible set in two, recursively.
func (c Client) FilterLogs(ctx context.Context, fq etherum.FilterQuery, currBlockHeight *big.Int, fn func(logs []ethtypes.Log) bool) (bool, error) {
	found, err := filterLogs(ctx, c.conn, fq, currBlockHeight, true, fn)

	if err != nil {
		log.WithError(err).Error("error filtering logs")
	}

	return found, err
}

func (c Client) TransactionByHash(ctx context.Context, txHash common.Hash) (*ethtypes.Transaction, bool, error) {
	return c.conn.TransactionByHash(ctx, txHash)
}

func (c Client) BlockByHash(ctx context.Context, blockHash common.Hash) (*ethtypes.Block, error) {
	return c.conn.BlockByHash(ctx, blockHash)
}

//go:generate mockery --name=ethClientToFilterLogs --inpackage --testonly
type ethClientToFilterLogs interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]etherumtypes.Log, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*etherumtypes.Header, error)
}

func shouldDoBinarySearchFromError(err error) bool {

	switch {
	case strings.Contains(err.Error(), "query returned more than 10000 results"):
		return true
	case strings.Contains(err.Error(), "eth_getLogs and eth_newFilter are limited to a 10,000 blocks range"):
		return true
	case strings.Contains(err.Error(), "block range is too wide"):
		return true
	}

	return false
}

// filterLogs filters for logs in a recursive manner. If the server returns
// that the block range is too high, then it does a binary search for left and
// right sectin.
func filterLogs(
	ctx context.Context,
	ethClient ethClientToFilterLogs,
	fq etherum.FilterQuery,
	currBlockHeight *big.Int,
	// reverseOrder if set to true then it searches latest logs first
	reverseOrder bool,
	fn func(logs []ethtypes.Log) bool,
) (bool, error) {
	log.
		WithField("current-block-height", currBlockHeight).
		WithField("filter-query", fq).
		Trace("filtering logs")

	if currBlockHeight == nil {
		header, err := ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return false, err
		}
		currBlockHeight = header.Number
	}

	if fq.BlockHash == nil {
		if fq.ToBlock == nil {
			fq.ToBlock = currBlockHeight
		}
		if fq.FromBlock == nil {
			fq.FromBlock = big.NewInt(0)
		}
	}

	logs, err := ethClient.FilterLogs(ctx, fq)

	switch {
	case err == nil:
		// awesome!
		if len(logs) == 0 {
			return false, nil
		}
		slice.ReverseInplace(logs)
		return fn(logs), nil
	case shouldDoBinarySearchFromError(err):
		// this appears to be ropsten specifict, but keepeing the logic here just in case
		mid := big.NewInt(0).Sub(
			fq.ToBlock,
			fq.FromBlock,
		)
		mid.Div(mid, big.NewInt(2))
		mid.Add(fq.FromBlock, mid)

		left := func() (bool, error) {
			fqLeft := fq
			fqLeft.ToBlock = mid
			return filterLogs(
				ctx,
				ethClient,
				fqLeft,
				currBlockHeight,
				reverseOrder,
				fn,
			)
		}

		right := func() (bool, error) {
			fqRight := fq
			fqRight.FromBlock = big.NewInt(0).Add(mid, big.NewInt(1))
			return filterLogs(
				ctx,
				ethClient,
				fqRight,
				currBlockHeight,
				reverseOrder,
				fn,
			)
		}

		var first, second func() (bool, error)

		if reverseOrder {
			first, second = right, left
		} else {
			first, second = left, right
		}

		foundFirst, err := first()
		if err != nil {
			return false, err
		}

		if foundFirst {
			return true, nil
		}

		return second()

	}

	return false, err
}

func (c Client) ExecuteSmartContract(
	ctx context.Context,
	chainID *big.Int,
	contractAbi abi.ABI,
	addr common.Address,
	method string,
	arguments []any,
) (*etherumtypes.Transaction, error) {
	return callSmartContract(
		ctx,
		executeSmartContractIn{
			ethClient:     c.conn,
			chainID:       chainID,
			gasAdjustment: c.config.GasAdjustment,
			abi:           contractAbi,
			contract:      addr,
			signingAddr:   c.addr,
			keystore:      c.keystore,

			method:    method,
			arguments: arguments,
		},
	)
}

func (c Client) BalanceAt(ctx context.Context, address common.Address, blockHeight uint64) (*big.Int, error) {
	return c.conn.BalanceAt(ctx, address, new(big.Int).SetUint64(blockHeight))
}

func (c Client) FindBlockNearestTo(ctx context.Context, startFromHeight uint64, when time.Time) (*etherumtypes.Block, error) {
	blockHeight, err := c.conn.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	from, to := new(big.Int).SetUint64(startFromHeight), new(big.Int).SetUint64(blockHeight)
	var block *etherumtypes.Block
	var err error
	for from.{
		block, err = c.conn.BlockByNumber(ctx, from)
		if err != nil {
			return nil, err
		}

		blockTime := time.Unix(int64(block.Time()), 0)
		if blockTime.Before(when) {
			from = block.Number()
		} else {
			to = block.Number()
		}
	}

	return block, nil
}
