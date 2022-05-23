package evm

import (
	"context"
	"embed"
	"path/filepath"
	"strings"
	"sync"

	"io/fs"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

/*
Do not delete hello.json contract. It's used for tests!
*/
var (
	//go:embed contracts/*.json
	contractsFS embed.FS

	readOnce   sync.Once
	_contracts = make(map[string]abi.ABI)
)

func StoredContracts() map[string]abi.ABI {
	readOnce.Do(func() {
		fs.WalkDir(contractsFS, ".", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			file, err := contractsFS.Open(path)

			contractName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			if err != nil {
				panic(err)
			}
			evmabi, err := abi.JSON(file)
			if err != nil {
				panic(err)
			}

			_contracts[contractName] = evmabi
			return nil
		})
	})
	return _contracts
}

type Client struct {
	ChainID  int
	HostPort string

	client any
}

func (c Client) ExecuteRemoteSmartContract(ctx context.Context, contract string, arguments []byte) error {
	// contracts := StoredContracts()
	// tx, err := contracts["hello"].Pack("bla")
	return nil
}

func (c Client) UpdateValset(ctx context.Context)            {}
func (c Client) ExecuteArbitraryMessage(ctx context.Context) {}
