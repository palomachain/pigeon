package app

import (
	"fmt"
	"net/url"
)

const (
	ChainName = "paloma"
	Name      = "sparrow"
)

type rootConfig struct {
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
	// envVal := os.Getenv(k)
	return "test"
}

func (k keyringEnvKey) Password() string {
	// envVal := os.Getenv(k)
	return "test"
}

func defaultConfigLocation() string {
	return fmt.Sprintf("~/.%s/%s", ChainName, Name)
}
