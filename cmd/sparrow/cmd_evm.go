package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/client/evm"
	"github.com/palomachain/sparrow/config"
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

			c := evm.NewClient(config.EVM{
				EVMSpecificClientConfig: config.EVMSpecificClientConfig{
					SmartContractAddress: "0x5A3E98aA540B2C3545311Fc33d445A7F62EB16Bf",
				},
				ChainClientConfig: config.ChainClientConfig{
					ChainID:            "3",
					BaseRPCURL:         "https://ropsten.infura.io/v3/d697ced03e7c49209a1fe2a1c8858821",
					KeyringPassEnvName: "blaa",
					SigningKey:         "0x621307FceE20F70Dd856F4EFF91bB1E21154105E",
					KeyringDirectory:   "/tmp",
					GasAdjustment:      1.1,
				},
			})

			return c.ExecuteArbitraryMessage(context.Background())
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
