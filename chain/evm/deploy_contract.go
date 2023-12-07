package evm

import (
	"context"
	"math/big"
	"strings"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/internal/libchain"
	"github.com/palomachain/pigeon/internal/liblog"
	arbaccounts "github.com/roodeag/arbitrum/accounts"
	arbabi "github.com/roodeag/arbitrum/accounts/abi"
	arbbind "github.com/roodeag/arbitrum/accounts/abi/bind"
	arbkeystore "github.com/roodeag/arbitrum/accounts/keystore"
	arbcommon "github.com/roodeag/arbitrum/common"
	arbtypes "github.com/roodeag/arbitrum/core/types"
	log "github.com/sirupsen/logrus"
)

func (c *Client) DeployContract(
	ctx context.Context,
	chainID *big.Int,
	rawABI string,
	bytecode,
	constructorInput []byte,
) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	// HACK HACK HACK
	// Logic branching to switch to Arbitrum go-ethereum types
	if libchain.IsArbitrum(chainID) {
		return c.wrapArbitrumDeployment(ctx, chainID, rawABI, bytecode, constructorInput)
	}

	return deployContract(
		ctx,
		c.conn,
		c.keystore,
		c.addr,
		chainID,
		rawABI,
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
	rawABI string,
	bytecode []byte,
	constructorInput []byte,
	gasAdjustment float64,
	txType uint8,
) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	logger := liblog.WithContext(ctx).WithField("chainID", chainID)
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

		// https://github.com/VolumeFi/paloma/issues/1048
		txOpts.GasLimit = uint64(float64(txOpts.GasLimit) * 1.1)

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

		var contractABI abi.ABI
		contractABI, err = abi.JSON(strings.NewReader(rawABI))
		if err != nil {
			logger.WithError(err).Error("failed to parse raw JSON")
			return
		}

		// hack begins here:
		// constructor input arguments are already properly encoded, but
		// we need to unpack them here because bind.DeployContract function
		// expects arguments to come in "go" form
		constructorArgs, err := contractABI.Constructor.Inputs.Unpack(constructorInput)
		whoops.Assert(err)
		// hack ends here

		logger.Info("deploying contract")

		contractAddr, tx, _, err = bind.DeployContract(
			txOpts,
			contractABI,
			bytecode,
			ethClient,
			constructorArgs...,
		)
		constructorArgs, _ = contractABI.Constructor.Inputs.Unpack(constructorInput)

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

func deployContractArbitrum(
	ctx context.Context,
	ethClient arbbind.ContractBackend,
	ks *arbkeystore.KeyStore,
	signingAddr arbcommon.Address,
	chainID *big.Int,
	contractAbi arbabi.ABI,
	bytecode []byte,
	constructorInput []byte,
	gasAdjustment float64,
	txType uint8,
) (contractAddr arbcommon.Address, tx *arbtypes.Transaction, err error) {
	logger := log.WithField("chainID", chainID)
	err = whoops.Try(func() {
		nonce, err := ethClient.PendingNonceAt(ctx, signingAddr)
		whoops.Assert(err)

		gasPrice, err := ethClient.SuggestGasPrice(ctx)
		whoops.Assert(err)

		txOpts, err := arbbind.NewKeyStoreTransactorWithChainID(
			ks,
			arbaccounts.Account{Address: signingAddr},
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

		// https://github.com/VolumeFi/paloma/issues/1048
		txOpts.GasLimit = uint64(float64(txOpts.GasLimit) * 1.1)

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

		contractAddr, tx, _, err = arbbind.DeployContract(
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

func (c *Client) wrapArbitrumDeployment(
	ctx context.Context,
	chainID *big.Int,
	rawABI string,
	bytecode,
	constructorInput []byte,
) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	logger := liblog.WithContext(ctx).WithField("caller", "deploy-contract-log")

	var arbContractABI arbabi.ABI
	arbContractABI, err = arbabi.JSON(strings.NewReader(rawABI))
	if err != nil {
		logger.WithError(err).Error("failed to parse ABI JSON")
		return
	}

	keystore := arbkeystore.NewKeyStore(c.config.KeyringDirectory.Path(), arbkeystore.StandardScryptN, arbkeystore.StandardScryptP)
	if !keystore.HasAddress(arbcommon.Address(c.addr)) {
		err = errors.Unrecoverable(ErrAddressNotFoundInKeyStore.Format(c.config.SigningKey, c.config.KeyringDirectory.Path()))
		if err != nil {
			logger.WithError(err).Error("failed to unlock keystore")
			return
		}
	}

	var addr *arbcommon.Address = &arbcommon.Address{}
	addr.SetBytes(c.addr.Bytes())
	acc := arbaccounts.Account{Address: *addr}

	whoops.Assert(keystore.Unlock(acc, config.KeyringPassword(c.config.KeyringPassEnvName)))

	var atx *arbtypes.Transaction
	_, atx, err = deployContractArbitrum(
		ctx,
		c.arbcon,
		keystore,
		*addr,
		chainID,
		arbContractABI,
		bytecode,
		constructorInput,
		c.config.GasAdjustment,
		c.config.TxType,
	)

	if err != nil {
		logger.WithError(err).Error("failed to deploy contract to arbitrum")
		return
	}

	v, r, s := atx.RawSignatureValues()
	tx = ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     atx.Nonce(),
		GasTipCap: atx.GasTipCap(),
		GasFeeCap: atx.GasFeeCap(),
		Gas:       atx.Gas(),
		To:        (*common.Address)(atx.To()),
		Value:     atx.Value(),
		Data:      atx.Data(),
		V:         v,
		R:         r,
		S:         s,
	})
	return
}
