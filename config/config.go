package config

import (
	"io"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/vizualni/whoops"

	"gopkg.in/yaml.v2"
)

const (
	ChainName = "paloma"
	Name      = "sparrow"
)

type ChainClientConfig struct {
	ChainID            string   `yaml:"chain-id"`
	BaseRPCURL         string   `yaml:"base-rpc-url"`
	KeyringPassEnvName string   `yaml:"keyring-pass-env-name"`
	KeyringType        string   `yaml:"keyring-type"`
	KeyringDirectory   filepath `yaml:"keyring-dir"`
	CallTimeout        string   `yaml:"call-timeout"`
	GasAdjustment      float64  `yaml:"gas-adjustment"`
	AccountPrefix      string   `yaml:"account-prefix"`
	GasPrices          string   `yaml:"gas-prices"`
}

type filepath string

func (f filepath) Path() string {
	p := string(f)
	homeDir := whoops.Must(user.Current()).HomeDir
	p = strings.ReplaceAll(p, "~", homeDir)
	return path.Clean(p)
}

type Root struct {
	Paloma Paloma `yaml:"paloma"`
	Terra  Terra  `yaml:"terra"`
}

type Paloma struct {
	ChainClientConfig `yaml:",inline"`

	SigningKeyName   string `yaml:"signing-key-name"`
	ValidatorAddress string `yaml:"validator-address"`
}

type Terra struct {
	ChainClientConfig `yaml:",inline"`

	RewardsAccount string   `yaml:"rewards-account"` // TODO: rethink this
	Accounts       []string `yaml:"acc-addresses"`
}

func KeyringPassword(envKey string) string {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		panic(ErrUnableToLocateKeyringEnvironmentVar)
	}
	return envVal
}

func FromReader(r io.Reader) (Root, error) {
	var cnf Root
	rawBody, err := io.ReadAll(r)
	if err != nil {
		return Root{}, err
	}
	str := string(rawBody)
	str = os.ExpandEnv(str)
	err = yaml.Unmarshal([]byte(str), &cnf)
	if err != nil {
		return Root{}, err
	}

	return cnf, nil
}
