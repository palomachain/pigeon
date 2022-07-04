# Pigeon

## ISSUES

Please use https://github.com/palomachain/paloma/issues to submit issues and add pigeon label!


# Instructions to test sending messages to EVM

## Install

```shell
wget -O - https://github.com/palomachain/pigeon/releases/download/{version}/{version}_Linux_x86_64.tar.gz | \
tar -C /usr/local/bin -xvzf - pigeon
chmod +x /usr/local/bin/pigeon
mkdir ~/.pigeon

# setting up the EVM keys for ethereuem mainnet.
# don't forget your password!
pigeon evm keys generate-new ~/.pigeon/keys/evm/eth-main
```

# or import existing you existing Ethereum evm private keys 
pigeon evm keys import ~/.pigeon/keys/evm/eth-main


### Config setup

`VALIDATOR="$(palomad keys list --list-names | head -n1)"`

Create configuration file here `~/.pigeon/config.yaml`

```yaml
loop-timeout: 5s

paloma:
  chain-id: paloma-testnet-6
  call-timeout: 20s
  keyring-dir: ~/.paloma
  keyring-pass-env-name: PALOMA_KEYRING_PASS
  keyring-type: test
  signing-key: ${VALIDATOR}
  base-rpc-url: http://localhost:26657
  gas-adjustment: 1.5
  gas-prices: 0.001ugrain
  account-prefix: paloma

evm:
  eth-main:
    chain-id: 1
    base-rpc-url: $ETH_RPC_URL
    keyring-pass-env-name: ETH_PASSWORD
    signing-key: ETH_SIGNING_KEY
    keyring-dir: ~/.pigeon/keys/evm/eth-main
```


### Start pigeon

Open a new terminal window and run:

```shell
PALOMA_KEYRING_PASS=<your Paloma key password>
ETH_RPC_URL=<Your Ethereum mainnet RPC URL>
ETH_PASSWORD=<Your ETH Key Password>
ETH_SIGNING_KEY=<Your ETH SIGNING KEY>
pigeon start
```

- Remember to run pigeon as a systemd Service! If you have a good systemd service implementation for Ubuntu, please make a PR on this README and we will add it.


### Definitions and Descriptions of Pigeons Variables
  - for paloma key:
	- keyring-dir
      - right now it's not really super important where this points. The important things for the future is that pigeon needs to send transactions to Paloma using it's validator (operator) key!
	  - it's best to leave it as is
	- keyring-pass-env-name
	  - this one is super important!
	  - it is the name of the ENV variable where password to unlock the keyring is stored!
	  - you are not writing password here!! You are writing the ENV variable's name where the password is stored.
	  - you should obviously use a bit more advanced method than shown here, but here is the example:
	    - if the `keyring-pass-env-name` is set to `MY_SUPER_SECRET_PASS` then you should provide ENV variable `MY_SUPER_SECRET_PASS` and store the password there
	    - e.g. `MY_SUPER_SECRET_PASS=abcd pigeon start`
	- keyring-type
	  - it should be the same as it's defined for paloma's client. Look under the ~/.paloma/config/client.toml
	- signing-key
	  - right now it's again not important which key we are using. It can be any key that has enough balance to submit TXs to Paloma. It's best to use the same key that's set up for the validator.
	- validator-address
	  - a bit redundant, but it's **your** validator's address.
	- gas-adustment:
	  - gas multiplier. The pigeon will estimate the gas to run a TX and then it will multiply it with gas-adjustment (if it's a positive number)
 - for evm -> ropsten:
    - base-rpc-url
	  - this one is my private url but you are more than welcome to provide your own
	  - you can also change this to the mainnet if you want to
	- keyring-pass-env-name: as as above for paloma.
	- signing-key
	  - address of the key from the keyring used to sign and send TXs to EVM network (one that you got when running `pigeon evm keys generate-new` from the install section)
	- keyring-dir:
	  - a directory where keys to communicate with the EVM network is stored

