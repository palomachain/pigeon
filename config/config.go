package config

import (
	"fmt"
	"net/url"
)

const (
	ChainName = "paloma"
	Name      = "sparrow"
)

type Root struct {
	Paloma Paloma
	Terra  Terra
}

type Paloma struct {
	KeyringDetails keyringEnvKey

	ValidatorAccountName string
	SigningKeyName       string

	BaseRPCURL *url.URL
}

type Terra struct {
	KeyringDetails keyringEnvKey

	RewardsAccount string
	Accounts       []string
}

type keyringEnvKey string

func (k keyringEnvKey) Type() string {
	// TODO:
	// envVal := os.Getenv(k)
	return "test"
}

func (k keyringEnvKey) Password() string {
	// TODO:
	// envVal := os.Getenv(k)
	return "test"
}

func defaultConfigLocation() string {
	return fmt.Sprintf("~/.%s/%s", ChainName, Name)
}
