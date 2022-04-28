package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	ChainName = "paloma"
	Name      = "sparrow"
)

type ChainClientConfig struct {
	ChainID            string  `yaml:"chain-id"`
	BaseRPCURL         string  `yaml:"base-rpc-url"`
	KeyringPassEnvName string  `yaml:"keyring-pass-env-name"`
	KeyringType        string  `yaml:"keyring-type"`
	KeyHomeDirectory   string  `yaml:"key-home"`
	CallTimeout        string  `yaml:"call-timeout"`
	GasAdjustment      float64 `yaml:"gas-adjustment"`
	AccountPrefix      string  `yaml:"account-prefix"`
	GasPrices          string  `yaml:"gas-prices"`
}

type Root struct {
	Paloma Paloma `yaml:"paloma"`
	Terra  Terra  `yaml:"terra"`
}

type Paloma struct {
	ChainClientConfig `yaml:",inline"`

	ValidatorAccountName string `yaml:"validator-account-name"`
	SigningKeyName       string `yaml:"signing-key-name"`
}

type Terra struct {
	ChainClientConfig `yaml:",inline"`

	RewardsAccount string   `yaml:"rewards-account"` // TODO: rethink this
	Accounts       []string `yaml:"acc-addresses"`
}

func KeyringPassword(envKey string) string {
	envVal := os.Getenv(envKey)
	return envVal
}

func FromReader(r io.Reader) (Root, error) {
	var cnf Root
	err := yaml.NewDecoder(r).Decode(&cnf)
	if err != nil {
		return Root{}, err
	}

	return cnf, nil
}
