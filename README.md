<div align="center">
  <h1> Pigeon </h1>
    <img alt="Paloma" src="https://github.com/palomachain/pigeon/blob/master/assets/paloma.png" />
</div>

> A Golang cross-chain message relayer system
> for Paloma validators to deliver messages to any blockchain.

<div align="center">
  <a href="https://github.com/palomachain/pigeon/blob/master/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/palomachain/pigeon.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/palomachain/pigeon">
    <img src="https://pkg.go.dev/badge/github.com/palomachain/pigeon.svg" alt="Go Reference">
  </a>
  <a href="https://goreportcard.com/report/github.com/palomachain/pigeon">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/palomachain/pigeon" />
  </a>
  <a href="https://github.com/palomachain/pigeon/blob/master/.github/workflows/ci-test.yml">
    <img alt="Code Coverage" src="https://github.com/palomachain/pigeon/actions/workflows/ci-test.yml/badge.svg?branch=master" />
  </a>
  <a href="https://github.com/palomachain/pigeon/blob/master/.github/workflows/release.yml">
    <img alt="Code Coverage" src="https://github.com/palomachain/pigeon/actions/workflows/release.yml/badge.svg?branch=master" />
  </a>
</div>

For Crosschain software engineers that want simultaneous control of multiple smart contracts, on any blockchain, Paloma is decentralized and consensus-driven message delivery, fast state awareness, low cost state computation, and powerful attestation system that enables scalable, crosschain, smart contract execution with any data source.

## Table of Contents

- [Talk To Us](#talk-to-us)
- [Releases](#releases)
- [Active Networks](#active-networks)
- [Issues](#issues)
- [Install](#install)

## Talk to us

We have active, helpful communities on Twitter and Telegram.

- [Twitter](https://twitter.com/paloma_chain)
- [Telegram](https://t.me/palomachain)
- [Discord](https://discord.gg/HtUvgxvh5N)
- [Forum](https://forum.palomachain.com/)

## Releases

See [Release procedure](CONTRIBUTING.md#release-procedure) for more information about the release model.

## Active Networks

- Paloma Testnest `paloma-testnet-16` (May 23, 2024)
- Paloma Mainnet `tumbler` (April 22, 2024)
- Arbitrum Mainnet (relay)
- Base Mainnet (relay)
- Blast Mainnet (relay)
- Binance Smart Chain Mainnet (relay)
- Ethereum Mainnet (relay)
- Gnosis Mainnet (relay)
- Optimism Mainnet (relay)
- Polygon Mainnet (relay)

## ISSUES

This repo does not accept issues. Please use <https://github.com/palomachain/paloma/issues> to submit issues and add pigeon label!

## Install

**If you are upgrading from a prior tesntet confirm that you added the `health-check-port: 5757` to your pigeon yaml configuration file and upgrade the paloma chain-id field to `paloma-testnet-16` for Testnet or `tumbler` for Mainnet (see example below).**

> #### Note
>
> If you're joining while testnet didn't boot up yet you may see a log line saying `not staking. waiting`. That's OK.
> If you see this after the chains starts producing blocks, then it means that your validator has been jailed.

> #### Note
>
> Some have seen errors with GLIBC version differences with the downloaded binaries.  This is caused by a difference in the libraries of the host that built the binary and the host running the binary.
>
> If you experience these errors, please pull down the code and build it, rather than downloading the prebuilt binary

### To get the latest prebuilt `pigeon` binary

```shell
wget -O - https://github.com/palomachain/pigeon/releases/download/v1.12.2/pigeon_Linux_x86_64.tar.gz  | \
  sudo tar -C /usr/local/bin -xvzf - pigeon
sudo chmod +x /usr/local/bin/pigeon

mkdir ~/.pigeon
```

### To build pigeon using latest

```shell
git clone https://github.com/palomachain/pigeon.git
cd pigeon
git checkout v1.12.2
make build
sudo mv ./build/pigeon /usr/local/bin/pigeon

mkdir ~/.pigeon
```

If you're upgrading to the most recent version, you will need to stop `pigeond` before removing the old binary and copying the new binary into place.

## Set up your EVM Keys. Don't forget your passwords

### Create a new key

Ethereum Mainnet (eth-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/eth-main
```

Binance Smart Chain Mainnet (bnb-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/bnb-main
```

Polygon Mainnet (matic-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/matic-main
```

Optimism Mainnet (op-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/op-main
```

Base Mainnet (base-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/base-main
```

Arbitrum Mainnet (arbitrum-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/arbitrum-main
```

Gnosis Mainnet (gnosis-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/gnosis-main
```

Blast Mainnet (blast-main)

```
pigeon evm keys generate-new ~/.pigeon/keys/evm/blast-main
```

### or import existing you existing Ethereum evm private keys

Ethereum Mainnet (eth-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/eth-main
```

Binance Smart Chain Mainnet (bnb-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/bnb-main
```

Polygon Mainnet (matic-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/matic-main
```

Optimism Mainnet (op-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/op-main
```

Base Mainnet (base-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/base-main
```

Arbitrum Mainnet (arbitrum-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/arbitrum-main
```

Gnosis Mainnet (gnosis-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/gnosis-main
```

Blast Mainnet (blast-main)

```
pigeon evm keys import ~/.pigeon/keys/evm/blast-main
```

### Config setup

**IMPORTANT VALIDATOR NOTE:** `gas-adjustment` is important in your pigeon settings. The gas adjustment values in the example below are set to ensure that your relay is able to increase gas required during periods of high congestion on the target chain. We propose that for `tx-type: 2`, `gas-adjustment` is set to `2` except for `bnb-main` chain. However, given that Paloma assigns messages based on pigeon performance, make sure your gas adjustment for each target chain maximizes your message delivery, while minimizing your relay costs.

#### Paloma supported chains on `tumbler` mainnet as of January 5, 2024

Validators joining the Paloma network for the first time will need to fund the minimum onchain gas-fee balances to be included in the validator set. All validators must support all chains for relaying. All validators must maintain the minimum balance required for each target chain. Below are the list of supported chains and their minimums.

| Chain Name    | Paloma Chain-ID | Native Token | Minimum Balance Required | Governance Proposal|
| ------------- | --------------- |------------- |-------------- |---------------- |
| BNB Chain / Binance Smart Chain |bnb-main | BNB  | 0.005 BNB  | [PIP 27 - Paloma Mainnet Support for BNB Chain](https://paloma.explorers.guru/proposal/9) |
| Ethereum  | eth-main  | ETH  | 0.005 ETH  | [PIP 28 - Paloma Compass-EVM and Mainnet Support for Ethereum Mainnet](https://paloma.explorers.guru/proposal/10)|
| Polygon MATIC  | matic-main  | MATIC  | 0.005 MATIC  | [PIP 34 - Paloma Messenger Mainnet Validator Support for Polygon Mainnet](https://paloma.explorers.guru/proposal/15)|
| Optimism  | op-main  | ETH  | 0.005 ETH  | [PIP 36 - Paloma Mainnet Support for Optimistic EVM blockchain, Optimism](https://paloma.explorers.guru/proposal/22)|
| Base  | base-main  | ETH  | 0.001 ETH  | [PIP 37 - Paloma Mainnet Support for Base EVM blockchain](https://paloma.explorers.guru/proposal/23)|
| Gnosis  | base-main  | xDAI  | 10 xDAI  | [PIP 46 - Paloma Mainnet Support for Gnosis EVM blockchain](https://paloma.explorers.guru/proposal/36)|
| Arbitrum  | arbitrum-main  | ETH  | 0.005 ETH  | [PIP 52 - Paloma Mainnet Support for Arbitrum EVM blockchain](https://paloma.explorers.guru/proposal/43)|
| BLAST  | blast-main  | ETH  | 0.005 ETH  | [PIP-55 - Paloma Messenger Support for BLAST!](https://paloma.explorers.guru/proposal/46)|

#### Configuration updates

Please make sure you restart Pigeon after making changes to your configuration files for the updates to come into effect.

Make sure your Paloma Cosmos-SDK keys are stored and available on your environment.

`palomad keys add "$VALIDATOR" --recover`

Set the VALIDATOR env variable

`export VALIDATOR="$(palomad keys list --list-names | head -n1)"`

Create configuration file here `~/.pigeon/config.yaml`

```yaml
loop-timeout: 5s
health-check-port: 5757

paloma:
  chain-id: tumbler
  call-timeout: 20s
  keyring-dir: ~/.paloma
  keyring-pass-env-name: PALOMA_KEYRING_PASS
  keyring-type: os
  validator-key: ${VALIDATOR}
  base-rpc-url: http://localhost:26657
  gas-adjustment: 3.0
  gas-prices: 0.01ugrain
  account-prefix: paloma

evm:
  eth-main:
    chain-id: 1
    base-rpc-url: ${ETH_RPC_URL}
    keyring-pass-env-name: ETH_PASSWORD
    signing-key: ${ETH_SIGNING_KEY}
    keyring-dir: ~/.pigeon/keys/evm/eth-main
    gas-adjustment: 2
    tx-type: 2

  bnb-main:
    chain-id: 56
    base-rpc-url: ${BNB_RPC_URL}
    keyring-pass-env-name: BNB_PASSWORD
    signing-key: ${BNB_SIGNING_KEY}
    keyring-dir: ~/.pigeon/keys/evm/bnb-main
    gas-adjustment: 1
    tx-type: 0

  matic-main:
    chain-id: 137
    base-rpc-url: ${MATIC_RPC_URL}
    keyring-pass-env-name: MATIC_PASSWORD
    signing-key: ${MATIC_SIGNING_KEY}
    keyring-dir: ~/.pigeon/keys/evm/matic-main
    gas-adjustment: 2
    tx-type: 2
  
  op-main:
    chain-id: 10
    base-rpc-url: ${OP_RPC_URL}
    keyring-pass-env-name: OP_PASSWORD
    signing-key: ${OP_SIGNING_KEY}
    keyring-dir: /root/.pigeon/keys/evm/op-main
    gas-adjustment: 2
    tx-type: 2

  base-main:
    chain-id: 8453
    base-rpc-url: ${BASE_RPC_URL}
    keyring-pass-env-name: BASE_PASSWORD
    signing-key: ${BASE_SIGNING_KEY}
    keyring-dir: /root/.pigeon/keys/evm/base-main
    gas-adjustment: 2
    tx-type: 2

  arbitrum-main:
    chain-id: 42161
    base-rpc-url: ${ARB_RPC_URL}
    keyring-pass-env-name: ARB_PASSWORD
    signing-key: ${ARB_SIGNING_KEY}
    keyring-dir: /root/.pigeon/keys/evm/arbitrum-main
    gas-adjustment: 2
    tx-type: 2

 gnosis-main:
    chain-id: 100
    base-rpc-url: ${GNOSIS_RPC_URL}
    keyring-pass-env-name: GNOSIS_PASSWORD
    signing-key: ${GNOSIS_SIGNING_KEY}
    keyring-dir: /root/.pigeon/keys/evm/gnosis-main
    gas-adjustment: 2
    tx-type: 2

  blast-main:
    chain-id: 81457
    base-rpc-url: ${BLAST_RPC_URL}
    keyring-pass-env-name: BLAST_PASSWORD
    signing-key: ${BLAST_SIGNING_KEY}
    keyring-dir: ~/.pigeon/keys/evm/blast-main
    gas-adjustment: 2
    tx-type: 2
```

#### Support for multiple signing keys

By default, Pigeon will use your validator key to sign any transactions sent to Paloma. In high throughput environments, this may lead to `account sequence mismatch` errors.
It's possible to define more than one signing key to be used in rotation to combat this issue. In order to do so, you will need to create a number of new keys and register them with Pigeon like this:

```bash
# First, create a number of new keys. You may create an arbitrary amount of keys.
palomad keys add pigeon-operator-alpha
palomad keys add pigeon-operator-bravo
palomad keys add pigeon-operator-charlie

# Second, your new addresses will need to receive an active feegrant from your validator address.
# This step is very important. It's not enough to simply fund those addresses manually.
# The active feegrant is considered a "permission" to send transactions from your validator address".
palomad tx feegrant grant $VALIDATOR_ADDRESS $(palomad keys show pigeon-operator-alpha-$(hostname) -a) --fees 500ugrain -y
palomad tx feegrant grant $VALIDATOR_ADDRESS $(palomad keys show pigeon-operator-bravo-$(hostname) -a) --fees 500ugrain -y
palomad tx feegrant grant $VALIDATOR_ADDRESS $(palomad keys show pigeon-operator-charlie-$(hostname) -a) --fees 500ugrain -y
```

After creating your signing keys, all you need to do is register them with Pigeon by adding the following to your pigeon config. Make sure to restart the service after making these changes.

```yaml
paloma:
  signing-keys:
    - pigeon-operator-alpha
    - pigeon-operator-bravo
    - pigeon-operator-charlie
```

### Start pigeon

First pigeon will need some keys:

```shell
cat <<EOT >~/.pigeon/env.sh
PALOMA_KEYRING_PASS=<your Paloma key password>
ETH_RPC_URL=<Your Ethereum mainnet RPC URL>
ETH_PASSWORD=<Your ETH Key Password>
ETH_SIGNING_KEY=<Your ETH SIGNING KEY>
BNB_RPC_URL=<Your Binance mainnet RPC URL>
BNB_PASSWORD=<Your BNB Key Password>
BNB_SIGNING_KEY=<Your BNB SIGNING KEY>
MATIC_RPC_URL=<Your Binance mainnet RPC URL>
MATIC_PASSWORD=<Your BNB Key Password>
MATIC_SIGNING_KEY=<Your BNB SIGNING KEY>
BASE_RPC_URL=<Your Base mainnet RPC URL>
BASE_PASSWORD=<Your Base Key Password>
BASE_SIGNING_KEY=<Your Base SIGNING KEY>
GNOSIS_RPC_URL=<Your Gnosis mainnet RPC URL>
GNOSIS_PASSWORD=<Your Gnosis Key Password>
GNOSIS_SIGNING_KEY=<Your Gnosis SIGNING KEY>
ARB_RPC_URL=<Your Arbitrum mainnet RPC URL>
ARB_PASSWORD=<Your Arbitrum Key Password>
ARB_SIGNING_KEY=<Your Arbitrum SIGNING KEY>
BLAST_RPC_URL=<Your Blast mainnet RPC URL>
BLAST_PASSWORD=<Your Blast Key Password>
BLAST_SIGNING_KEY=<Your Blast SIGNING KEY>
OP_RPC_URL=<Your Optimism mainnet RPC URL>
OP_PASSWORD=<Your Optimism Key Password>
OP_SIGNING_KEY=<Your Optimism SIGNING KEY>
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
  - right now it's not really super important where this points. The important things for the future is that pigeon needs to send transactions to Paloma using its validator (operator) key!
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
  - keyring-pass-env-name: same as above for paloma.
  - signing-key
    - address of the key from the keyring used to sign and send TXs to EVM network (one that you got when running `pigeon evm keys generate-new` from the install section)
  - keyring-dir:
    - a directory where keys to communicate with the EVM network is stored
