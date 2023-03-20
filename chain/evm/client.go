package evm

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"io/ioutil"
	"math"
	"math/big"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/VolumeFi/whoops"
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
	"github.com/gogo/protobuf/proto"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/types/paloma/x/evm/types"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
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
		err := fs.WalkDir(contractsFS, ".", func(path string, d fs.DirEntry, err error) error {
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
		if err != nil {
			log.WithField("err", err).Error("error iterating over the stored contracts")
		}
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

	conn ethClientConn

	paloma PalomaClienter
}

var _ ethClientConn = &ethclient.Client{}

//go:generate mockery --name=ethClientConn --inpackage --testonly
type ethClientConn interface {
	bind.ContractBackend
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *etherumtypes.Transaction, isPending bool, err error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*etherumtypes.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*etherumtypes.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
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

		whoops.Assert(c.keystore.Unlock(acc, config.KeyringPassword(c.config.KeyringPassEnvName)))

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
	txType        uint8

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
		"chain-id":        args.chainID,
		"contract-addr":   args.contract,
		"method":          args.method,
		"arguments":       args.arguments,
		"gas-adjustments": args.gasAdjustment,
		"signing-addr":    args.signingAddr,
	})
	return whoops.TryVal(func() *etherumtypes.Transaction {
		packedBytes, err := args.abi.Pack(
			args.method,
			args.arguments...,
		)
		if err != nil {
			logger.
				WithField("error", err).
				Error("callSmartContract: error packing input")
		}
		whoops.Assert(err)

		nonce, err := args.ethClient.PendingNonceAt(ctx, args.signingAddr)
		if err != nil {
			logger.
				WithField("error", err).
				Error("callSmartContract: error calculating pending nonce")
		}
		whoops.Assert(err)

		gasPrice, err := args.ethClient.SuggestGasPrice(ctx)
		if err != nil {
			logger.
				WithField("error", err).
				Error("callSmartContract: error calculating pending nonce")
		}
		whoops.Assert(err)

		// adjusting the gas price
		if args.txType != 2 && args.gasAdjustment > 1.0 {
			gasAdj := big.NewFloat(args.gasAdjustment)
			gasAdj = gasAdj.Mul(gasAdj, new(big.Float).SetInt(gasPrice))
			gasPrice, _ = gasAdj.Int(big.NewInt(0))
		}

		var gasTipCap *big.Int

		if args.txType == 2 {
			gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2)) // double gas price for EIP-1559 tx
			gasTipCap, err = args.ethClient.SuggestGasTipCap(ctx)
			if err != nil {
				logger.
					WithField("error", err).
					Error("callSmartContract: error calling SuggestGasTipCap")
			}
			whoops.Assert(err)
			gasPrice = gasPrice.Add(gasPrice, gasTipCap)
			logger.WithFields(log.Fields{
				"gas-max-price": gasPrice,
				"gas-max-tip":   gasTipCap,
			}).Debug("adjusted eip-1559 gas price")
		} else {
			logger.WithFields(log.Fields{
				"gas-price": gasPrice,
			}).Debug("adjusted legacy gas price")
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
		if err != nil {
			logger.
				WithField("error", err).
				Error("callSmartContract: error calling bind.NewKeyStoreTransactorWithChainID")
		}
		whoops.Assert(err)

		txOpts.Nonce = big.NewInt(int64(nonce))
		txOpts.From = args.signingAddr

		if args.txType == 2 {
			txOpts.GasFeeCap = gasPrice
			txOpts.GasTipCap = gasTipCap
			logger.WithFields(log.Fields{
				"gas-limit":     txOpts.GasLimit,
				"gas-max-price": txOpts.GasFeeCap,
				"gas-max-tip":   txOpts.GasTipCap,
				"nonce":         txOpts.Nonce,
				"from":          txOpts.From,
			}).Debug("executing eip-1559 tx")
		} else {
			txOpts.GasPrice = gasPrice
			logger.WithFields(log.Fields{
				"gas-limit": txOpts.GasLimit,
				"gas-price": txOpts.GasPrice,
				"nonce":     txOpts.Nonce,
				"from":      txOpts.From,
			}).Debug("executing legacy tx")
		}

		tx, err := boundContract.RawTransact(txOpts, packedBytes)
		if err != nil {
			logger.
				WithField("error", err).
				Error("callSmartContract: error calling boundContract.RawTransact")
		}
		whoops.Assert(err)

		if args.txType == 2 {
			logger.WithFields(log.Fields{
				"tx-hash":          tx.Hash(),
				"tx-gas-limit":     tx.Gas(),
				"tx-gas-max-price": tx.GasFeeCap(),
				"tx-gas-max-tip":   tx.GasTipCap(),
				"tx-cost":          tx.Cost(),
			}).Debug("eip-1559 tx executed")
		} else {
			logger.WithFields(log.Fields{
				"tx-hash":      tx.Hash(),
				"tx-gas-limit": tx.Gas(),
				"tx-gas-price": tx.GasPrice(),
				"tx-cost":      tx.Cost(),
			}).Debug("legacy tx executed")
		}

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
	case strings.Contains(err.Error(), "exceed maximum block range"):
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
			txType:        c.config.TxType,
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
	var bh *big.Int
	if blockHeight > 0 {
		bh = new(big.Int).SetUint64(blockHeight)
	}
	return c.conn.BalanceAt(ctx, address, bh)
}

func (c Client) FindBlockNearestToTime(ctx context.Context, startingHeight uint64, when time.Time) (uint64, error) {
	isTimeSetBeforeBlock := func(height uint64) (bool, error) {
		h, err := c.conn.HeaderByNumber(ctx, new(big.Int).SetUint64(height))
		if err != nil {
			return false, err
		}
		return h.Time < uint64(when.UTC().Unix()), nil
	}

	before, err := isTimeSetBeforeBlock(startingHeight)
	if err != nil {
		return 0, err
	}
	if !before {
		return 0, ErrStartingBlockIsInTheFuture
	}

	currBlockHeight, err := c.conn.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	from, to := startingHeight, currBlockHeight
	var res uint64
	for from <= to {
		err := whoops.Try(func() {
			mid := uint64(math.Round(float64(from+to) / 2))
			before := whoops.Must(isTimeSetBeforeBlock(mid))
			if before {
				res = mid
				from = mid + 1
			} else {
				to = mid - 1
			}
		})
		if err != nil {
			return 0, err
		}
	}

	if res == currBlockHeight {
		// there needs to be at least one block standing in between
		return 0, ErrBlockNotYetGenerated
	}

	return res, nil
}

func (c Client) FindCurrentBlockNumber(ctx context.Context) (*big.Int, error) {
	header, err := c.conn.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}
