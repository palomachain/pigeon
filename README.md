![Logo!](assets/paloma.png)

# Pigeon

> A Golang cross-chain message relayer system
> for Paloma validators to deliver messages to any blockchain.

For Crosschain software engineers that want simultaneous control of mulitiple smart contracts, on any blockchain, Paloma is decentralized and consensus-driven message delivery, fast state awareness, low cost state computation, and powerful attestation system that enables scaleable, crosschain, smart contract execution with any data source.

## Table of Contents

- [Talk To Us](#talk-to-us)
- [Releases](#releases)
- [Active Networks](#active-networks)
- [Issues](#issues)
- [Install](#install)


## Talk to us

We have active, helpful communities on Twitter and Telegram.

* [Twitter](https://twitter.com/paloma_chain)
* [Telegram](https://t.me/palomachain)

## Releases

See [Release procedure](CONTRIBUTING.md#release-procedure) for more information about the release model.

## Active Networks

### Public Testnest 8




## ISSUES

This repo does not accept issues. Please use https://github.com/palomachain/paloma/issues to submit issues and add pigeon label!


## Install

```shell
wget -O - https://github.com/palomachain/pigeon/releases/download/v0.5.1/pigeon_Linux_x86_64.tar.gz | \
tar -C /usr/local/bin -xvzf - pigeon
chmod +x /usr/local/bin/pigeon
mkdir ~/.pigeon
```

If you're upgrading to the most recent version, you will need to stop `pigeond` before removing the old binary and copying the new binary into place.

If you want to install from source use 
```go build ./cmd/pigeon```


## Set up your EVM Keys. Don't forget your passwords!

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/eth-main
```

or import existing you existing Ethereum evm private keys

```
pigeon evm keys import ~/.pigeon/keys/evm/eth-main
```

### Config setup

Make sure your Paloma Cosmos-SDK keys are stored and available on your environment.

`palomad keys add "$VALIDATOR" --recover`

Set the VALIDATOR env variable

`export VALIDATOR="$(palomad keys list --list-names | head -n1)"`

Create configuration file here `~/.pigeon/config.yaml`

```yaml
loop-timeout: 5s
health-check-port: 5757

paloma:
  chain-id: paloma-testnet-8
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
    base-rpc-url: ${ETH_RPC_URL}
    keyring-pass-env-name: ETH_PASSWORD
    signing-key: ${ETH_SIGNING_KEY}
    keyring-dir: ~/.pigeon/keys/evm/eth-main
```


### Start pigeon

First pigeon will need some keys:

```shell
cat <<EOT >~/.pigeon/env.sh
PALOMA_KEYRING_PASS=<your Paloma key password>
ETH_RPC_URL=<Your Ethereum mainnet RPC URL>
ETH_PASSWORD=<Your ETH Key Password>
ETH_SIGNING_KEY=<Your ETH SIGNING KEY>
VALIDATOR=<VALIDATOR NAME>
EOT
```

Then we can run pigeon with:

```shell
source ~/.pigeon/env.sh
pigeon start
```

#### Using systemd

Make sure you have configured `.pigeon/env.sh` as above. Then create a systemctl configuration:

```shell
cat <<EOT >/etc/systemd/system/pigeond.service
[Unit]
Description=Pigeon daemon
After=network-online.target
ConditionPathExists=/usr/local/bin/pigeon

[Service]
Type=simple
Restart=always
RestartSec=5
User=${USER}
WorkingDirectory=~
EnvironmentFile=${HOME}/.pigeon/env.sh
ExecStart=/usr/local/bin/pigeon start
ExecReload=

[Install]
WantedBy=multi-user.target
EOT
```

Then start our pigeon!

```shell
service pigeond start

# Check that it's running successfully:
service pigeond status
# Or watch the logs:
journalctl -u pigeond.service -f -n 100
```

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
	- gas-adustment:
	  - gas multiplier. The pigeon will estimate the gas to run a TX and then it will multiply it with gas-adjustment (if it's a positive number)
 - for evm -> eth-main:
	- keyring-pass-env-name: as as above for paloma.
	- signing-key
	  - address of the key from the keyring used to sign and send TXs to EVM network (one that you got when running `pigeon evm keys generate-new` from the install section)
	- keyring-dir:
	  - a directory where keys to communicate with the EVM network is stored

