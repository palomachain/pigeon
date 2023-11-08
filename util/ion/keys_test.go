package ion_test

import (
	"testing"

	"github.com/palomachain/pigeon/util/ion"
)

// TestKeyRestore restores a test mnemonic
func TestKeyRestore(t *testing.T) {
	t.Skip()
	keyName := "test_key"
	mnemonic := "blind master acoustic speak victory lend kiss grab glad help demand hood roast zone lend sponsor level cheap truck kingdom apology token hover reunion"
	expectedAddress := "cosmos15cw268ckjj2hgq8q3jf68slwjjcjlvxy57je2u"
	coinType := uint32(118) // Cosmos coin type used in address derivation

	homepath := t.TempDir()
	cl, err := ion.NewClient(
		ion.GetCosmosHubConfig(homepath, true), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = cl.DeleteKey(keyName) // Delete if test is being run again
	address, err := cl.RestoreKey(keyName, mnemonic, coinType)
	if err != nil {
		t.Fatalf("Error while restoring mnemonic: %v", err)
	}
	if address != expectedAddress {
		t.Fatalf("Restored address: %s does not match expected: %s", address, expectedAddress)
	}
}

// TestKeyRestoreEth restores a test mnemonic
func TestKeyRestoreEth(t *testing.T) {
	t.Skip()
	keyName := "test_key"
	mnemonic := "three elevator silk family street child flip also leaf inmate call frame shock little legal october vivid enable fetch siege sell burger dolphin green"
	expectedAddress := "evmos1dea7vlekr9e34vugwkvesulglt8fx4e457vk9z"
	coinType := uint32(60) // Ethereum coin type used in address derivation

	cl, err := ion.NewClient(
		&ion.ChainClientConfig{
			Key:            "default",
			ChainID:        "evmos_9001-2",
			AccountPrefix:  "evmos",
			KeyringBackend: "test",
			GasAdjustment:  1.2,
			GasPrices:      "0.01uevmos",
			Timeout:        "20s",
			OutputFormat:   "json",
			SignModeStr:    "direct",
		}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = cl.DeleteKey(keyName) // Delete if test is being run again
	address, err := cl.RestoreKey(keyName, mnemonic, coinType)
	if err != nil {
		t.Fatalf("Error while restoring mnemonic: %v", err)
	}
	if address != expectedAddress {
		t.Fatalf("Restored address: %s does not match expected: %s", address, expectedAddress)
	}
}

// TestKeyRestoreinj restores a test mnemonic
func TestKeyRestoreInj(t *testing.T) {
	t.Skip()
	keyName := "inj_key"
	mnemonic := "three elevator silk family street child flip also leaf inmate call frame shock little legal october vivid enable fetch siege sell burger dolphin green"
	expectedAddress := "inj1dea7vlekr9e34vugwkvesulglt8fx4e4uk2udj"
	coinType := uint32(60) // Ethereum coin type used in address derivation

	cl, err := ion.NewClient(
		&ion.ChainClientConfig{
			Key:            "default",
			ChainID:        "injective-1",
			AccountPrefix:  "inj",
			KeyringBackend: "test",
			GasAdjustment:  1.2,
			GasPrices:      "0.01inj",
			Timeout:        "20s",
			OutputFormat:   "json",
			SignModeStr:    "direct",
			ExtraCodecs:    []string{"injective"},
		}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = cl.DeleteKey(keyName) // Delete if test is being run again
	address, err := cl.RestoreKey(keyName, mnemonic, coinType)
	if err != nil {
		t.Fatalf("Error while restoring mnemonic: %v", err)
	}
	if address != expectedAddress {
		t.Fatalf("Restored address: %s does not match expected: %s", address, expectedAddress)
	}
	_, err = cl.ListAddresses()
	if err != nil {
		t.Fatalf("Error while restoring mnemonic: %v", err)
	}
	err = cl.DeleteKey(keyName) // Delete if test is being run again
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}
}
