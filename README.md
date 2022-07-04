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

# setting up the EVM keys for ropsten testnet.
# don't forget your password!
pigeon evm keys generate-new ~/.pigeon/keys/evm/ropsten
```

### Config setup

VALIDATOR="$(palomad keys list --list-names | head -n1)"
PUBKEY="$(palomad tendermint show-validator)"


```yaml
loop-timeout: 5s

paloma:
  chain-id: paloma-testnet-6
  call-timeout: 20s
  keyring-dir: ~/.paloma
  keyring-pass-env-name: PALOMA_KEYRING_PASS
  keyring-type: testnest
  signing-key: VALIDATOR
  base-rpc-url: http://localhost:26657
  gas-adjustment: 2.0
  gas-prices: 0.001ugrain
  account-prefix: paloma

evm:
  ethreum-mainnet:
    chain-id: 1
    base-rpc-url: ETH_RPC_URL
    keyring-pass-env-name: ETH_PASSWORD
    signing-key: ETH_SIGNING_KEY
    keyring-dir: ~/.pigeon/keys/evm/ethereum-mainnet
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

- Open pigeon window and look at the logs to get the TX HASH which you can look on the explorer.

- Remember to run pigeon as a systemd Service! If you have a good systemd service implementation for Ubuntu, please make a PR on this README and we will add it.


### Paloma Notes for 

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


Once you are done setting this up, you can take the config and put it here `~/.pigeon/config.yaml`.


A careful reader might notice that there can be as many EVM networks as you wish. That is true, but at the moment of writing this, only the first one will be used.
(To be more correct, all of them are going to be used, but the first one will take all the messages and send them given that there is no routing yet for sending
messages to different EVM chains.)

### Open faucet and send coins to your key

https://faucet.egorfine.com/

Send some coins to your newly generated ETH address.

### Prepare Smart Contract info that you wish to execute

Now remember that you need to take the smart contract info from the network you are trying to execute this on. You can use ropsten as I did, but you can also do it on mainnet.
Find the smart contract (or upload one) that you want to run and take its address and its ABI json encoded description.


Paste the abi here: https://abi.hashex.org/

and prepare the data for your function that you wish to run. Once you are done take the ABI encoded payload and use it to add a new job to Paloma.


### Add a new job to Paloma chain to execute a smart contract
Take the smart contract address, the ABI encoded payload and the contract's JSON ABI and run the `paloma tx evm submit-new-job` command:

```shell
palomad tx evm submit-new-job [smart contract address] [smart contract payload] [smart contract JSON abi]
```

An example of how I did it (note that the third argument (the Smart Contract JSON abi) is wrapped with `'` single quotes):

```shell
palomad tx evm submit-new-job --from my_validator --fees 200dove --broadcast-mode block -y 0x5a3e98aa540b2c3545311fc33d445a7f62eb16bf 6057361d0000000000000000000000000000000000000000000000000000000000001688 '[{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]'
```


## Can you run this on the mainnet?

Yes! Generate yourself a new key (or use the one that you've generated for testing this on ropsten) and send some tokens to your newly created address. Change the rpc-url to point to the mainnet.
Find a smart contract from the mainnet and its address to submit a job.

Right now we can't import existing keys, so you need to use generated ones.

# Notes!

- Right now there is no routing of jobs based on which network they are supposed to belong to. It simply sends it to whatever is the first chain defined in the config file under the evm key.
- pigeons are not signing anything now. All pigeons are going to try their best to get the available messages that are in the queue and they will not sign those messages. They are simply going to take them and send them away.
