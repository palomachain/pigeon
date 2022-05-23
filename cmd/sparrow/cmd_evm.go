package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/palomachain/sparrow/client/evm"
	"github.com/spf13/cobra"
)

var (
	evmCmd = &cobra.Command{
		Use: "evm",
	}
	evmDebugCmd = &cobra.Command{
		Use:    "debug",
		Hidden: true,
	}
	debugContractsCmd = &cobra.Command{
		Use:   "contracts",
		Short: "shows info about loaded contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			contracts := evm.StoredContracts()
			for name, contract := range contracts {
				fmt.Printf("%s: %#v\n", name, contract)
			}

			return nil
		},
	}

	debugEvmConnectToNetworkCmd = &cobra.Command{
		Use:   "connect",
		Short: "tries to connect to the evm network",
		RunE: func(cmd *cobra.Command, args []string) error {
			// addr := accounts.("0x621307FceE20F70Dd856F4EFF91bB1E21154105E")
			// ks := evm.OpenKeystore("/tmp")
			hw, err := abi.JSON(strings.NewReader(`[{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
			if err != nil {
				panic(err)
			}
			pks := "a68db8652c3d31c4e40b05fe7f0bf020d0e4b119a276ee385824d1a73dd134a9"
			pk, err := crypto.HexToECDSA(pks)
			if err != nil {
				panic(err)
			}
			addr := crypto.PubkeyToAddress(pk.PublicKey)

			conn, err := ethclient.Dial("https://ropsten.infura.io/v3/d697ced03e7c49209a1fe2a1c8858821")
			if err != nil {
				log.Fatalf("Failed to connect to the Ethereum client: %v", err)
			}

			nonce, err := conn.PendingNonceAt(context.Background(), addr)
			fmt.Println(nonce)
			if err != nil {
				log.Fatalf("Failed to connect to the Ethereum client: %v", err)
			}

			fmt.Println(conn.BalanceAt(context.Background(), addr, nil))
			bc := bind.NewBoundContract(common.HexToAddress("0x5A3E98aA540B2C3545311Fc33d445A7F62EB16Bf"), hw, conn, conn, conn)

			// opts, err := bind.NewKeyStoreTransactor(ks, ks.Accounts()[0])
			opts, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(3))
			if err != nil {
				log.Fatalf("Failed to create authorized transactor: %v", err)
			}
			packedBytes, err := hw.Pack(
				"store",
				big.NewInt(1337),
			)
			gasPrice, err := conn.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatalf("Failed to create authorized transactor: %v", err)
			}
			// value := big.NewInt(100000000000000) // in wei (1 eth)
			fmt.Println(gasPrice)

			opts.Nonce = big.NewInt(int64(nonce))
			// opts.Value = value

			opts.GasLimit = 210000
			opts.GasPrice = gasPrice

			opts.From = addr

			tx, err := bc.RawTransact(opts, packedBytes)
			fmt.Println(err)
			if tx != nil {
				fmt.Println(tx)
				fmt.Println(tx.Hash())
			}

			return nil
		},
	}
	debugEvmDeploySmartContractCmd = &cobra.Command{
		Use:   "deploy",
		Short: "tries to connect to the evm network",
		RunE: func(cmd *cobra.Command, args []string) error {
			// addr := "0x621307FceE20F70Dd856F4EFF91bB1E21154105E"
			// ks := evm.OpenKeystore("/tmp")

			// ks.Find()
			return nil
		},
	}

	evmKeysCmd = &cobra.Command{
		Use: "keys",
	}
	evmKeysListCmd = &cobra.Command{
		Use:   "list",
		Short: "lists accounts in the keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore("/tmp")
			for _, acc := range ks.Accounts() {
				fmt.Println(acc)
			}
			return nil
		},
	}
	evmKeysGenerateCmd = &cobra.Command{
		Use:   "generate-new",
		Short: "generates a new account and adds it to keystore",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore("/tmp")
			ecdsaPK, err := crypto.GenerateKey()
			if err != nil {
				return err
			}
			pass := doubleReadInput("Password to encode your key: ", true, 3)
			acc, err := ks.ImportECDSA(ecdsaPK, pass)
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println(acc)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(evmCmd)

	evmDebugCmd.AddCommand(
		debugContractsCmd,
		debugEvmConnectToNetworkCmd,
	)

	evmCmd.AddCommand(
		evmDebugCmd,
		evmKeysCmd,
	)

	evmKeysCmd.AddCommand(
		evmKeysListCmd,
		evmKeysGenerateCmd,
	)
}
