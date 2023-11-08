package ion

import (
	"errors"
	"os"

	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"github.com/strangelove-ventures/lens/client/codecs/ethermint"
	"github.com/strangelove-ventures/lens/client/codecs/injective"
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Evmos:
	//  - secp256k1     (Cosmos)
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithms = keyring.SigningAlgoList{hd.Secp256k1, ethermint.EthSecp256k1, injective.EthSecp256k1}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Evmos for the Ledger device:
	//  - secp256k1     (Cosmos)
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{hd.Secp256k1, ethermint.EthSecp256k1, injective.EthSecp256k1}
)

// Option defines a function keys options for the ethereum Secp256k1 curve.
// It supports secp256k1 and eth_secp256k1 keys for accounts.
func LensKeyringAlgoOptions() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
	}
}

func (cc *Client) CreateKeystore(path string) error {
	keybase, err := keyring.New(cc.Config.ChainID, cc.Config.KeyringBackend, cc.Config.KeyDirectory, cc.Input, cc.Codec.Marshaler, LensKeyringAlgoOptions())
	if err != nil {
		return err
	}
	cc.Keybase = keybase
	return nil
}

func (cc *Client) KeystoreCreated(path string) bool {
	if _, err := os.Stat(cc.Config.KeyDirectory); errors.Is(err, os.ErrNotExist) {
		return false
	} else if cc.Keybase == nil {
		return false
	}
	return true
}

func (cc *Client) AddKey(name string, coinType uint32) (output *KeyOutput, err error) {
	ko, err := cc.KeyAddOrRestore(name, coinType)
	if err != nil {
		return nil, err
	}
	return ko, nil
}

func (cc *Client) RestoreKey(name, mnemonic string, coinType uint32) (address string, err error) {
	ko, err := cc.KeyAddOrRestore(name, coinType, mnemonic)
	if err != nil {
		return "", err
	}
	return ko.Address, nil
}

func (cc *Client) ShowAddress(name string) (address string, err error) {
	info, err := cc.Keybase.Key(name)
	if err != nil {
		return "", err
	}
	acc, err := info.GetAddress()
	if err != nil {
		return "", nil
	}
	out, err := cc.EncodeBech32AccAddr(acc)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (cc *Client) ListAddresses() (map[string]string, error) {
	out := map[string]string{}
	info, err := cc.Keybase.List()
	if err != nil {
		return nil, err
	}
	for _, k := range info {
		acc, err := k.GetAddress()
		if err != nil {
			return nil, err
		}
		addr, err := cc.EncodeBech32AccAddr(acc)
		if err != nil {
			return nil, err
		}
		out[k.Name] = addr
	}
	return out, nil
}

func (cc *Client) DeleteKey(name string) error {
	if err := cc.Keybase.Delete(name); err != nil {
		return err
	}
	return nil
}

func (cc *Client) KeyExists(name string) bool {
	k, err := cc.Keybase.Key(name)
	if err != nil {
		return false
	}

	return k.Name == name
}

func (cc *Client) ExportPrivKeyArmor(keyName string) (armor string, err error) {
	return cc.Keybase.ExportPrivKeyArmor(keyName, ckeys.DefaultKeyPass)
}

func (cc *Client) KeyAddOrRestore(keyName string, coinType uint32, mnemonic ...string) (*KeyOutput, error) {
	var mnemonicStr string
	var err error
	algo := keyring.SignatureAlgo(hd.Secp256k1)

	if len(mnemonic) > 0 {
		mnemonicStr = mnemonic[0]
	} else {
		mnemonicStr, err = CreateMnemonic()
		if err != nil {
			return nil, err
		}
	}

	if coinType == 60 {
		algo = keyring.SignatureAlgo(ethermint.EthSecp256k1)
		for _, codec := range cc.Config.ExtraCodecs {
			if codec == "injective" {
				algo = keyring.SignatureAlgo(injective.EthSecp256k1)
			}
		}
	}

	info, err := cc.Keybase.NewAccount(keyName, mnemonicStr, "", hd.CreateHDPath(coinType, 0, 0).String(), algo)
	if err != nil {
		return nil, err
	}

	acc, err := info.GetAddress()
	if err != nil {
		return nil, err
	}

	out, err := cc.EncodeBech32AccAddr(acc)
	if err != nil {
		return nil, err
	}
	return &KeyOutput{Mnemonic: mnemonicStr, Address: out}, nil
}

// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	Mnemonic string `json:"mnemonic" yaml:"mnemonic"`
	Address  string `json:"address" yaml:"address"`
}

// CreateMnemonic creates a new mnemonic
func CreateMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}
