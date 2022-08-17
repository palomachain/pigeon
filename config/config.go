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
	Name      = "pigeon"
)

type CosmosSpecificClientConfig struct {
	KeyringType   string `yaml:"keyring-type"`
	AccountPrefix string `yaml:"account-prefix"`
	GasPrices     string `yaml:"gas-prices"`
}

type ChainClientConfig struct {
	BaseRPCURL         string   `yaml:"base-rpc-url"`
	KeyringPassEnvName string   `yaml:"keyring-pass-env-name"`
	SigningKey         string   `yaml:"signing-key"`
	KeyringDirectory   Filepath `yaml:"keyring-dir"`
	CallTimeout        string   `yaml:"call-timeout"`
	GasAdjustment      float64  `yaml:"gas-adjustment"`
}

type Filepath string

func (f Filepath) Path() string {
	p := string(f)
	homeDir := whoops.Must(user.Current()).HomeDir
	p = strings.ReplaceAll(p, "~", homeDir)
	return path.Clean(p)
}

type Root struct {
	HealthCheckPortRaw int `yaml:"health-check-port"`

	Paloma Paloma `yaml:"paloma"`

	EVM map[string]EVM `yaml:"evm"`
}

func (r *Root) HealthCheckPort() int {
	if r.HealthCheckPortRaw == 0 {
		panic(whoops.String("invalid health check port in pigeon's config file"))
	}
	return r.HealthCheckPortRaw
}

func (r *Root) init() {
	(&r.Paloma).init()
}

type EVM struct {
	ChainClientConfig `yaml:",inline"`
}

type Paloma struct {
	CosmosSpecificClientConfig `yaml:",inline"`
	ChainClientConfig          `yaml:",inline"`
	ChainID                    string `yaml:"chain-id"`
}

func (p *Paloma) init() {
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

	(&cnf).init()

	return cnf, nil
}
