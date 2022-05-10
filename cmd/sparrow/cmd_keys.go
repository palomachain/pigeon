package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var (
	keysCmd = &cobra.Command{
		Use: "keys",
	}
	keysConvertCmd = &cobra.Command{
		Use:   "convert-prefix [address] [from-prefix] [to-prefix]",
		Short: "initializes the sparrow",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			address, from, to := args[0], args[1], args[2]
			hrp, err := sdk.GetFromBech32(address, from)
			if err != nil {
				return err
			}
			newAddress, err := sdk.Bech32ifyAddressBytes(to, hrp)
			if err != nil {
				return err
			}
			fmt.Println(newAddress)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(keysConvertCmd)

	noConfigRequired(keysCmd)
	noConfigRequired(keysConvertCmd)
}
