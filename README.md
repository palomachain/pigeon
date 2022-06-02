# Conductor

## ISSUES

Please use https://github.com/palomachain/paloma/issues to submit issues and add sparrow label!


# Instructions to test sending messages to EVM

## Install

```shell
wget -O - https://github.com/palomachain/sparrow/releases/download/{version}/{version}_Linux_x86_64.tar.gz | \
tar -C /usr/local/bin -xvzf - sparrow
chmod +x /usr/local/bin/sparrow
mkdir ~/.sparrow

# setting up the EVM keys for ropsten testnet.
# don't forget your password!
sparrow evm keys generate-new ~/.sparrow/keys/evm/ropsten
```

### Config setup

```yaml
loop-timeout: 5s

paloma:
  chain-id: paloma
  call-timeout: 20s
  keyring-dir: ~/.paloma
  keyring-pass-env-name: PALOMA_KEYRING_PASS
  keyring-type: test
  signing-key: my_validator
  validator-address: palomavaloper107ur3w38qvjnx44a732ayphp0yfh4f0jmxsn0e
  base-rpc-url: http://localhost:26657
  gas-adjustment: 2.0
  gas-prices: 0.01dove
  account-prefix: paloma

evm:
  ropsten:
    chain-id: 3
    base-rpc-url: https://ropsten.infura.io/v3/d697ced03e7c49209a1fe2a1c8858821
    keyring-pass-env-name: ROPSTEN_PASS
    signing-key: 0x378d6991F6b5207d7cC2b5270AD2Afb3Dd328E82
    keyring-dir: ~/.sparrow/keys/evm/ropsten
```

Important things to set up are:
  - for paloma key:
	- keyring-dir
      - right now it's not really super important where this points. The important things for the future is that Sparrow needs to send transactions to Paloma using it's validator (operator) key!
	  - it's best to leave it as is
	- keyring-pass-env-name
	  - this one is super important!
	  - it is the name of the ENV variable where password to unlock the keyring is stored!
	  - you are not writing password here!! You are writing the ENV variable's name where the password is stored.
	  - you should obviously use a bit more advanced method than shown here, but here is the example:
	    - if the `keyring-pass-env-name` is set to `MY_SUPER_SECRET_PASS` then you should provide ENV variable `MY_SUPER_SECRET_PASS` and store the password there
	    - e.g. `MY_SUPER_SECRET_PASS=abcd sparrow start`
	- keyring-type
	  - it should be the same as it's defined for paloma's client. Look under the ~/.paloma/config/client.toml
	- signing-key
	  - right now it's again not important which key we are using. It can be any key that has enough balance to submit TXs to Paloma. It's best to use the same key that's set up for the validator.
	- validator-address
	  - a bit redundant, but it's **your** validator's address.
	- gas-adustment:
	  - gas multiplier. The Sparrow will estimate the gas to run a TX and then it will multiply it with gas-adjustment (if it's a positive number)
 - for evm -> ropsten:
    - base-rpc-url
	  - this one is my private url but you are more than welcome to provide your own
	  - you can also change this to the mainnet if you want to
	- keyring-pass-env-name: as as above for paloma.
	- signing-key
	  - address of the key from the keyring used to sign and send TXs to EVM network (one that you got when running `sparrow evm keys generate-new` from the install section)
	- keyring-dir:
	  - a directory where keys to communicate with the EVM network is stored


Once you are done setting this up, you can take the config and put it here `~/.sparrow/config.yaml`.


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


### Start sparrow

Open a new terminal window and run:

```shell
YOUR_ENV_VARIABLE_WITH_PALOMA_PASS=abcd
YOUR_ENV_VARIABLE_WITH_EVM_ROPSTEN_PASS=efgh
sparrow start
```

- Open sparrow window and look at the logs to get the TX HASH which you can look on the explorer.
- Feel free to add more jobs to the queue while Sparrows are running

## Can you run this on the mainnet?

Yes! Generate yourself a new key (or use the one that you've generated for testing this on ropsten) and send some tokens to your newly created address. Change the rpc-url to point to the mainnet.
Find a smart contract from the mainnet and its address to submit a job.

Right now we can't import existing keys, so you need to use generated ones.

# Notes!

- Right now there is no routing of jobs based on which network they are supposed to belong to. It simply sends it to whatever is the first chain defined in the config file under the evm key.
- Sparrows are not signing anything now. All sparrows are going to try their best to get the available messages that are in the queue and they will not sign those messages. They are simply going to take them and send them away.
