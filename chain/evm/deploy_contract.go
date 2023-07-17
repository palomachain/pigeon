package evm

import (
	"context"
	"math/big"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func (c *Client) DeployContract(
	ctx context.Context,
	chainID *big.Int,
	contractAbi abi.ABI,
	bytecode,
	constructorInput []byte,
) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	return deployContract(
		ctx,
		c.conn,
		c.keystore,
		c.addr,
		chainID,
		contractAbi,
		bytecode,
		constructorInput,
		c.config.GasAdjustment,
		c.config.TxType,
	)
}

func deployContract(
	ctx context.Context,
	ethClient bind.ContractBackend,
	ks *keystore.KeyStore,
	signingAddr common.Address,
	chainID *big.Int,
	contractAbi abi.ABI,
	bytecode []byte,
	constructorInput []byte,
	gasAdjustment float64,
	txType uint8,
) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	logger := log.WithField("chainID", chainID)
	err = whoops.Try(func() {
		nonce, err := ethClient.PendingNonceAt(ctx, signingAddr)
		whoops.Assert(err)

		gasPrice, err := ethClient.SuggestGasPrice(ctx)
		whoops.Assert(err)

		txOpts, err := bind.NewKeyStoreTransactorWithChainID(
			ks,
			accounts.Account{Address: signingAddr},
			chainID,
		)
		whoops.Assert(err)

		txOpts.Nonce = big.NewInt(int64(nonce))
		txOpts.From = signingAddr
		// adjusting the gas price
		if txType != 2 && gasAdjustment > 1.0 {
			gasAdj := big.NewFloat(gasAdjustment)
			gasAdj = gasAdj.Mul(gasAdj, new(big.Float).SetInt(gasPrice))
			gasPrice, _ = gasAdj.Int(big.NewInt(0))
		}

		var gasTipCap *big.Int

		if txType == 2 {
			gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2)) // double gas price for EIP-1559 tx
			gasTipCap, err = ethClient.SuggestGasTipCap(ctx)
			whoops.Assert(err)
			gasPrice = gasPrice.Add(gasPrice, gasTipCap)
			logger.WithFields(log.Fields{
				"gas-max-price": gasPrice,
				"gas-max-tip":   gasTipCap,
			}).Info("adjusted eip-1559 gas price")
			txOpts.GasFeeCap = gasPrice
			txOpts.GasTipCap = gasTipCap
			logger = logger.WithFields(log.Fields{
				"gas-limit":     txOpts.GasLimit,
				"gas-max-price": txOpts.GasFeeCap,
				"gas-max-tip":   txOpts.GasTipCap,
				"nonce":         txOpts.Nonce,
				"from":          txOpts.From,
				"tx-type":       2,
			})
		} else {
			logger.WithFields(log.Fields{
				"gas-price": gasPrice,
			}).Info("adjusted legacy gas price")
			txOpts.GasPrice = gasPrice
			logger = logger.WithFields(log.Fields{
				"gas-limit": txOpts.GasLimit,
				"gas-price": txOpts.GasPrice,
				"nonce":     txOpts.Nonce,
				"from":      txOpts.From,
				"tx-type":   0,
			})
		}

		// hack begins here:
		// constructor input arguments are already properly encoded, but
		// we need to unpack them here because bind.DeployContract function
		// expects arguments to come in "go" form
		constructorArgs, err := contractAbi.Constructor.Inputs.Unpack(constructorInput)
		whoops.Assert(err)
		// hack ends here

		logger.Info("deploying contract")

		contractAddr, tx, _, err = bind.DeployContract(
			txOpts,
			contractAbi,
			bytecode,
			ethClient,
			constructorArgs...,
		)
		constructorArgs, _ = contractAbi.Constructor.Inputs.Unpack(constructorInput)

		whoops.Assert(err)
		if tx.Type() == 2 {
			logger.WithFields(log.Fields{
				"tx-hash":          tx.Hash(),
				"tx-gas-limit":     tx.Gas(),
				"tx-gas-max-price": tx.GasFeeCap(),
				"tx-gas-max-tip":   tx.GasTipCap(),
				"tx-cost":          tx.Cost(),
				"tx-type":          tx.Type(),
			}).Info("eip-1559 tx executed")
		} else {
			logger.WithFields(log.Fields{
				"tx-hash":      tx.Hash(),
				"tx-gas-limit": tx.Gas(),
				"tx-gas-price": tx.GasPrice(),
				"tx-cost":      tx.Cost(),
				"tx-type":      tx.Type(),
			}).Info("legacy tx executed")
		}
	})
	return
}
