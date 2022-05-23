package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/client/evm"
	"github.com/spf13/cobra"
)

var (
	evmCmd = &cobra.Command{
		Use: "evm",
	}
	debugContractsCmd = &cobra.Command{
		Use:   "debug-contracts",
		Short: "shows info about loaded contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			contracts := evm.StoredContracts()
			for name, contract := range contracts {
				fmt.Println("%s: %#v", name, contract)
			}

			return nil
		},
	}
	debugEvmConnectToNetworkCmd = &cobra.Command{
		Use:   "debug-connect",
		Short: "tries to connect to the evm network",
		RunE: func(cmd *cobra.Command, args []string) error {
			contracts := evm.StoredContracts()
			for name, contract := range contracts {
				fmt.Println("%s: %#v", name, contract)
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

	evmCmd.AddCommand(
		debugContractsCmd,
		debugEvmConnectToNetworkCmd,
		evmKeysCmd,
	)

	evmKeysCmd.AddCommand(
		evmKeysListCmd,
		evmKeysGenerateCmd,
	)
}
