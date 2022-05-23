package evm

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func OpenKeystore(dir string) *keystore.KeyStore {
	return keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
}
