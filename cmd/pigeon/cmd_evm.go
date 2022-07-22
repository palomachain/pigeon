package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/palomachain/pigeon/app"
	"github.com/palomachain/pigeon/chain/evm"
	"github.com/spf13/cobra"
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
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := ethclient.Dial(app.Config().EVM["ropsten"].BaseRPCURL)
			if err != nil {
				return err
			}
			tx, _, _ := c.TransactionByHash(cmd.Context(), common.HexToHash("0xa5b948628c719cd8771d36d52c22ae6d41ce94ad0951f7b382e04d5b0028d246"))
			fmt.Println(tx)

			b, err := tx.MarshalBinary()
			if err != nil {
				return err
			}

			ioutil.WriteFile("/tmp/ooo.hex", []byte(common.Bytes2Hex(b)), 0o666)

			fmt.Println(b)
			fmt.Println(common.Bytes2Hex(b))
			fmt.Println(sha256.Sum256(b))
			return nil
			// chainID, contractABIbz, bytecode, packedInput := args[0], args[1], args[2], args[3]
			// c := app.GetEvmClients()[chainID]
			// contractABI := whoops.Must(abi.JSON(strings.NewReader(contractABIbz)))
			// addr, tx, err := c.DeployContract(
			// 	cmd.Context(),
			// 	contractABI,
			// 	common.FromHex(bytecode),
			// 	common.FromHex(packedInput),
			// )
			// if err != nil {
			// 	return err
			// }
			// log.WithFields(log.Fields{
			// 	"address": addr,
			// 	"tx":      tx,
			// }).Info("smart contract deployed")
			return nil
		},
	}
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
		Short: "imports an existing key into the keyring's directory",
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
		evmKeysImportCmd,
	)
}
