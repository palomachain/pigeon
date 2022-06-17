package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/app"
	"github.com/palomachain/sparrow/chain/evm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vizualni/whoops"
)

const (
	signaturePrefix = "\x19Ethereum Signed Message:\n32"
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
				fmt.Printf("%s: %#v\n", name, contract.ABI)
			}

			return nil
		},
	}

	evmKeysCmd = &cobra.Command{
		Use: "keys",
	}

	evmKeysListCmd = &cobra.Command{
		Use:   "list [directory]",
		Short: "lists accounts in the keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore(args[0])
			for _, acc := range ks.Accounts() {
				fmt.Println(acc)
			}
			return nil
		},
	}

	evmDeploySmartContractCmd = &cobra.Command{
		Use:   "deploy-smart-contract [chainID] [abi] [bytecode] [packed-input]",
		Short: "deploys a smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID, contractABIbz, bytecode, packedInput := args[0], args[1], args[2], args[3]
			c := app.GetEvmClients()[chainID]
			contractABI := whoops.Must(abi.JSON(strings.NewReader(contractABIbz)))
			constructorArgs, err := contractABI.Constructor.Inputs.Unpack(common.FromHex(packedInput))
			if err != nil {
				return err
			}
			addr, tx, err := c.DeployContract(
				cmd.Context(),
				contractABI,
				common.FromHex(bytecode),
				constructorArgs,
			)
			if err != nil {
				return err
			}
			log.WithFields(log.Fields{
				"address": addr,
				"tx":      tx,
			}).Info("smart contract deployed")
			return nil
		},
	}
	// TODO: add import
	evmKeysGenerateCmd = &cobra.Command{
		Use:   "generate-new [directory]",
		Short: "generates a new account and adds it to keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore(args[0])
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
			fmt.Println("Key created in", args[0], "directory")
			fmt.Println("Don't lose your password! Otherwise you'd lose access to your key!")
			fmt.Println(acc)
			return nil
		},
	}

	evmKeysImportCmd = &cobra.Command{
		Use:   "import [directory]",
		Short: "generates a new account and adds it to keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore(args[0])
			fmt.Println("Paste the private key in HEX format:")
			pkHex := readLineFromStdin(false)
			pk, err := crypto.HexToECDSA(pkHex)
			if err != nil {
				return err
			}
			pass := doubleReadInput("Password to encode your key: ", true, 3)
			acc, err := ks.ImportECDSA(pk, pass)
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println("Key created in", args[0], "directory")
			fmt.Println("Don't lose your password! Otherwise you'd lose access to your key!")
			fmt.Println(acc)
			return nil
		},
	}

	zero            [32]byte
	evmDebugSignCmd = &cobra.Command{
		Use:   "sign-test [directory]",
		Short: "generates a new account and adds it to keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore(args[0])
			hash := crypto.Keccak256([]byte{})
			protectedHash := crypto.Keccak256Hash(append([]byte(signaturePrefix), hash...))

			sig, err := ks.SignHashWithPassphrase(accounts.Account{
				Address: common.HexToAddress("0xe4Ab6f4D62Ba7e0bBC4CF6c5E8153e105108FBa9"),
			}, "aaaaaaaa", protectedHash.Bytes())
			if err != nil {
				panic(err)
			}

			fmt.Println(new(big.Int).SetBytes([]byte{(sig[64])}))
			fmt.Println(new(big.Int).SetBytes(sig[:32]))
			fmt.Println(new(big.Int).SetBytes(sig[32:64]))
			return nil
		},
	}
)

func init() {
	configRequired(evmDeploySmartContractCmd)
	rootCmd.AddCommand(evmCmd)

	evmDebugCmd.AddCommand(
		debugContractsCmd,
		evmDeploySmartContractCmd,
		evmDebugSignCmd,
	)

	evmCmd.AddCommand(
		evmDebugCmd,
		evmKeysCmd,
	)

	evmKeysCmd.AddCommand(
		evmKeysListCmd,
		evmKeysGenerateCmd,
		evmKeysImportCmd,
	)
}
