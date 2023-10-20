package config

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/VolumeFi/whoops"
	"gopkg.in/yaml.v2"
)

const (
	ChainName                          = "paloma"
	Name                               = "pigeon"
	cDefaultHealthServerAddressBinding = "127.0.0.1"
)

type CosmosSpecificClientConfig struct {
	KeyringType   string `yaml:"keyring-type"`
	AccountPrefix string `yaml:"account-prefix"`
	GasPrices     string `yaml:"gas-prices"`
}

type EVMSpecificClientConfig struct {
	TxType                      uint8 `yaml:"tx-type"`
	BloxrouteIntegrationEnabled bool  `yaml:"bloxroute-mev-enabled"`
}

type ChainClientConfig struct {
	BaseRPCURL         string   `yaml:"base-rpc-url"`
	KeyringPassEnvName string   `yaml:"keyring-pass-env-name"`
	SigningKey         string   `yaml:"signing-key"`
	KeyringDirectory   Filepath `yaml:"keyring-dir"`
	CallTimeout        string   `yaml:"call-timeout"`
	SigningKeys        []string `yaml:"signing-keys"`
	GasAdjustment      float64  `yaml:"gas-adjustment"`
}

type Filepath string

func (f Filepath) Path() string {
	p := string(f)
	homeDir := whoops.Must(user.Current()).HomeDir
	p = strings.ReplaceAll(p, "~", homeDir)
	return path.Clean(p)
}

type Config struct {
	HealthCheckPort    int    `yaml:"health-check-port"`
	HealthCheckAddress string `yaml:"health-check-address"`

	BloxrouteAuthorizationHeader string `yaml:"bloxroute-auth-header"`

	Paloma Paloma `yaml:"paloma"`

	EVM map[string]EVM `yaml:"evm"`
}

func (c *Config) defaults() *Config {
	if len(c.HealthCheckAddress) < 1 {
		c.HealthCheckAddress = cDefaultHealthServerAddressBinding
	}

	return c
}

func (c *Config) validate() (*Config, error) {
	if c.HealthCheckPort == 0 {
		return nil, fmt.Errorf("invalid health server port binding: %d", c.HealthCheckPort)
	}

	return c, nil
}

type EVM struct {
	EVMSpecificClientConfig `yaml:",inline"`
	ChainClientConfig       `yaml:",inline"`
}

type Paloma struct {
	CosmosSpecificClientConfig `yaml:",inline"`
	ChainClientConfig          `yaml:",inline"`
	ChainID                    string `yaml:"chain-id"`
}

func KeyringPassword(envKey string) string {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		panic(ErrUnableToLocateKeyringEnvironmentVar)
	}
	return envVal
}

func FromReader(r io.Reader) (*Config, error) {
	rawBody, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	str := string(rawBody)
	str = os.ExpandEnv(str)

	var cfg Config
	err = yaml.Unmarshal([]byte(str), &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.defaults().validate()
}
