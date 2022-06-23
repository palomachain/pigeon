package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

func (c Client) DeployContract(ctx context.Context, contractAbi abi.ABI, bytecode, constructorInput []byte) (contractAddr common.Address, tx *ethtypes.Transaction, err error) {
	return deployContract(
		ctx,
		c.conn,
		c.keystore,
		c.addr,
		c.config.GetChainID(),
		contractAbi,
		bytecode,
		constructorInput,
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
		txOpts.GasPrice = gasPrice
		txOpts.From = signingAddr

		logger = logger.WithFields(log.Fields{
			"tx-opts": txOpts,
		})

		// hack begins here
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
		whoops.Assert(err)

		logger.WithFields(log.Fields{
			"tx-hash":      tx.Hash(),
			"tx-gas-limit": tx.Gas(),
			"tx-gas-price": tx.GasPrice(),
			"tx-cost":      tx.Cost(),
		}).Info("tx executed")
	})
	return
}

func rawDeploy(opts *bind.TransactOpts, abi abi.ABI, bytecode []byte, backend bind.ContractBackend, packedConstructorInput []byte) (common.Address, *ethtypes.Transaction, error) {
	// Otherwise try to deploy the contract
	c := bind.NewBoundContract(common.Address{}, abi, backend, backend, backend)

	tx, err := c.RawTransact(opts, append(bytecode, packedConstructorInput...))
	if err != nil {
		return common.Address{}, nil, err
	}
	return crypto.CreateAddress(opts.From, tx.Nonce()), tx, nil
}
