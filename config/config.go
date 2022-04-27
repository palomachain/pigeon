package config

import (
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	ChainName = "paloma"
	Name      = "sparrow"
)

type ChainClientConfig struct {
	ChainID          string        `yaml:"chain-id"`
	BaseRPCURL       string        `yaml:"base-rpc-url"`
	KeyringDetails   keyringEnvKey `yaml:"keyring-env-name"`
	KeyHomeDirectory string        `yaml:"key-home"`
	CallTimeout      string        `yaml:"call-timeout"`
	GasAdjustment    float64       `yaml:"gas-adjustment"`
	AccountPrefix    string        `yaml:"account-prefix"`
	GasPrices        string        `yaml:"gas-prices"`
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

type keyringEnvKey string

type keyringDetails struct {
	typ  string
	pass string
}

func parseKeyringEnvValue(value string) (keyringDetails, error) {
	var zero keyringDetails
	values := strings.SplitN(value, ";", 2)

	if len(values) != 2 {
		return zero, ErrUnableToParseKeyringDetails
	}

	return keyringDetails{
		typ:  values[0],
		pass: values[1],
	}, nil
}

func (k keyringEnvKey) Type() string {
	envVal := os.Getenv(string(k))
	details, err := parseKeyringEnvValue(envVal)
	if err != nil {
		panic(err)
	}
	return details.typ
}

func (k keyringEnvKey) Password() string {
	envVal := os.Getenv(string(k))
	details, err := parseKeyringEnvValue(envVal)
	if err != nil {
		panic(err)
	}
	return details.pass
}

func FromReader(r io.Reader) (Root, error) {
	var cnf Root
	err := yaml.NewDecoder(r).Decode(&cnf)
	if err != nil {
		return Root{}, err
	}

	return cnf, nil
}
