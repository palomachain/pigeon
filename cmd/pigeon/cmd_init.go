package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [keyring-backend] [keyring-location]",
	Short: "initializes the pigeon",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if the signing key exists,
		// if it does not, it creates one which it reads from the passed in config.
		key := "signing-key"

		clientCtx := client.GetClientContextFromCmd(cmd)
		kr, err := keyring.New("pigeon", args[0], args[1], os.Stdin, clientCtx.Codec)
		if err != nil {
			return err
		}
		_, err = kr.Key(key)
		if err == nil {
			// nothing to do
			// TODO; there could be other things that need to be initialised,
			// so this message might go away from this line of code.
			fmt.Println("nothing to do for me")
			return nil
		}

		fmt.Println("Adding a new signing key to the keyring")
		info, _, err := kr.NewMnemonic(key, keyring.English, types.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
		if err != nil {
			return err
		}
		fmt.Printf("key '%s' added to the keyring in %s\n", key, args[1])
		fmt.Println(info)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
