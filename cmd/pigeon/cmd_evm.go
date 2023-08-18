package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/palomachain/pigeon/app"
	"github.com/palomachain/pigeon/chain/evm"
	"github.com/sirupsen/logrus"
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

			tmpF, err := os.CreateTemp("/tmp", "ooo.hex")
			if err != nil {
				return err
			}

			defer os.Remove(tmpF.Name())

			_, err = tmpF.Write([]byte(common.Bytes2Hex(b)))
			if err != nil {
				return err
			}

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
	evmKeysExportPrivateKey = &cobra.Command{
		Use:   "export [directory] [address] [export-location]",
		Short: "exports a private key",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := evm.OpenKeystore(args[0])
			acc, err := ks.Find(accounts.Account{
				Address: common.HexToAddress(args[1]),
			})
			if err != nil {
				return err
			}
			fmt.Println("Password to unlock:")
			pass := readLineFromStdin(true)
			bz, err := ks.Export(acc, pass, pass)
			if err != nil {
				return err
			}
			fmt.Println("Writing the private key in JSON format in: ", args[2])
			return os.WriteFile(args[2], bz, 0o600)
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

	evmVerifySignaturesInQueueCmd = &cobra.Command{
		Use:   "verify-signatures [queue-name]",
		Short: "Verifies signatures in the queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			msgs, err := app.PalomaClient().QueryMessagesForAttesting(ctx, args[0])
			if err != nil {
				return err
			}

			logger := logrus.WithField("queue-type-name", args[0])

			countOK, countNotOK := 0, 0
			for _, msg := range msgs {
				logger = logger.WithField("msg-id", msg.ID)
				for _, signature := range msg.Signatures {
					logger = logger.WithFields(
						logrus.Fields{
							"pub-key":   signature.PublicKey,
							"validator": signature.ValAddress.String(),
						},
					)
					foundPK, err := crypto.Ecrecover(
						crypto.Keccak256(
							append(
								[]byte(evm.SignedMessagePrefix),
								msg.BytesToSign...,
							),
						),
						signature.Signature,
					)
					if err != nil {
						logger.WithError(err).Error("processing")
						countNotOK++
						continue
					}
					pk, err := crypto.UnmarshalPubkey(foundPK)
					if err != nil {
						countNotOK++
						continue
					}

					recoveredAddr := crypto.PubkeyToAddress(*pk)

					if recoveredAddr.Hex() != signature.SignedByAddress {
						countNotOK++
						logger.Info("signature not ok")
					} else {
						countOK++
					}

				}
			}
			logger.Info("OK signatures: ", countOK, " not OK signatures: ", countNotOK)
			return nil
		},
	}
	evmVerifySignaturesAgainstValsetCmd = &cobra.Command{
		Use:   "verify-signatures-against-valset [queue-name] [chain-reference-id] [valset-id]",
		Short: "Verifies signatures in the queue",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			msgs, err := app.PalomaClient().QueryMessagesForAttesting(ctx, args[0])
			if err != nil {
				return err
			}
			valsetID, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			valset, err := app.PalomaClient().QueryGetEVMValsetByID(ctx, uint64(valsetID), args[1])
			if err != nil {
				return err
			}

			logger := logrus.WithField("queue-type-name", args[0])

			logger.Info("there are ", len(valset.Validators), " validators in the valset")

			countOK, total := 0, 0
			for _, msg := range msgs {
				logger = logger.WithField("msg-id", msg.ID)
				xconsensus := evm.BuildCompassConsensus(ctx, valset, msg.Signatures)
				logger.Info("message has ", len(msg.Signatures), " signatures")
				logger.Info("there are ", len(xconsensus.Valset.Validators), " validators in the valset validator consensus")
				for i := range xconsensus.Valset.Validators {
					valaddr := xconsensus.Valset.Validators[i]
					signature := xconsensus.OriginalSignatures()[i]
					if len(signature) == 0 {
						continue
					}

					total++
					logger = logger.WithField("validator-evm-addr", valaddr.Hex())

					foundPK, err := crypto.Ecrecover(
						crypto.Keccak256(
							append(
								[]byte(evm.SignedMessagePrefix),
								msg.BytesToSign...,
							),
						),
						signature,
					)
					if err != nil {
						logger.WithError(err).Error("processing")
						continue
					}
					pk, err := crypto.UnmarshalPubkey(foundPK)
					if err != nil {
						continue
					}

					recoveredAddr := crypto.PubkeyToAddress(*pk)

					if recoveredAddr.Hex() != valaddr.Hex() {
						logger.Info("signature not ok")
					} else {
						countOK++
					}

				}
			}
			logger.Info("total val: ", total, " OK signatures: ", countOK)
			return nil
		},
	}
)

func init() {
	configRequired(evmDeploySmartContractCmd)
	configRequired(evmVerifySignaturesInQueueCmd)
	configRequired(evmVerifySignaturesAgainstValsetCmd)

	rootCmd.AddCommand(evmCmd)

	evmDebugCmd.AddCommand(
		debugContractsCmd,
		evmDeploySmartContractCmd,
		evmVerifySignaturesInQueueCmd,
		evmVerifySignaturesAgainstValsetCmd,
	)

	evmCmd.AddCommand(
		evmDebugCmd,
		evmKeysCmd,
	)

	evmKeysCmd.AddCommand(
		evmKeysListCmd,
		evmKeysGenerateCmd,
		evmKeysImportCmd,
		evmKeysExportPrivateKey,
	)
}
