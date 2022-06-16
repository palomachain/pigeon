package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/app"
	"github.com/palomachain/sparrow/chain/evm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vizualni/whoops"
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
)

func init() {
	configRequired(evmDeploySmartContractCmd)
	rootCmd.AddCommand(evmCmd)

	evmDebugCmd.AddCommand(
		debugContractsCmd,
		evmDeploySmartContractCmd,
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
