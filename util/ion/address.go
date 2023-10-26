package ion

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (cc *Client) EncodeBech32AccAddr(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(cc.Config.AccountPrefix, addr)
}

func (cc *Client) MustEncodeAccAddr(addr sdk.AccAddress) string {
	enc, err := cc.EncodeBech32AccAddr(addr)
	if err != nil {
		panic(err)
	}
	return enc
}

func (cc *Client) EncodeBech32AccPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "pub"), addr)
}

func (cc *Client) EncodeBech32ValAddr(addr sdk.ValAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valoper"), addr)
}

func (cc *Client) MustEncodeValAddr(addr sdk.ValAddress) string {
	enc, err := cc.EncodeBech32ValAddr(addr)
	if err != nil {
		panic(err)
	}
	return enc
}

func (cc *Client) EncodeBech32ValPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valoperpub"), addr)
}

func (cc *Client) EncodeBech32ConsAddr(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valcons"), addr)
}

func (cc *Client) EncodeBech32ConsPub(addr sdk.AccAddress) (string, error) {
	return sdk.Bech32ifyAddressBytes(fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valconspub"), addr)
}

func (cc *Client) DecodeBech32AccAddr(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, cc.Config.AccountPrefix)
}

func (cc *Client) DecodeBech32AccPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "pub"))
}

func (cc *Client) DecodeBech32ValAddr(addr string) (sdk.ValAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valoper"))
}

func (cc *Client) DecodeBech32ValPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valoperpub"))
}

func (cc *Client) DecodeBech32ConsAddr(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valcons"))
}

func (cc *Client) DecodeBech32ConsPub(addr string) (sdk.AccAddress, error) {
	return sdk.GetFromBech32(addr, fmt.Sprintf("%s%s", cc.Config.AccountPrefix, "valconspub"))
}
