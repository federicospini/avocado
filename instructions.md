# Avocado

>  run cli, not use npm run!

## TEST CLI

```
node offchain/cli.js test_cli --myAddress=0x0077cdf65cAD0b5ca9C2af69CD6743f396f616F0 --myUrl=http://localhost:4021 --counterpartyAddress=0x00798eaa5c12FD38ff58aD51c985D805C02eB429 --counterpartyUrl=http://localhost:4021 
```

## PROPOSE CHANNEL

Alice proposes a new channel to Bob

```
node offchain/cli.js propose_channel --myAddress=0x0077cdf65cAD0b5ca9C2af69CD6743f396f616F0 --myUrl=http://localhost:4021 --counterpartyAddress=0x00798eaa5c12FD38ff58aD51c985D805C02eB429 --counterpartyUrl=http://localhost:4022 --channelId=0x0000000000000000000000000000000000000000000000000000000000000100 --state=0x10 --challengePeriod=80
```

## VIEW PROPOSED CHANNEL

Bob checks channels proposed to him

```
node offchain/cli.js view_proposed_channels --counterpartyAddress=0x00798eaa5c12FD38ff58aD51c985D805C02eB429 --counterpartyUrl=http://localhost:4022
```

## ACCEPT PROPOSED CHANNEL

Bob accepts a channel proposed to him

```
node offchain/cli.js accept_proposed_channel --myAddress=0x0077cdf65cAD0b5ca9C2af69CD6743f396f616F0 --myUrl=http://localhost:4021 --counterpartyAddress=0x00798eaa5c12FD38ff58aD51c985D805C02eB429 --counterpartyUrl=http://localhost:4022 --channelId=0x0000000000000000000000000000000000000000000000000000000000000100 
```

## PROPOSE UPDATE

```
node offchain/cli.js propose_update --channelId=0x0000000000000000000000000000000000000000000000000000000000000100 --state=0x0101
```


## How to unlock account 0

> from ethconsole, parity needs to be started with --geth

```
web3.eth.getAccounts(function (err, accounts) {
  web3.personal.unlockAccount(accounts[0], 'btclab', function (err, res) {
    console.log(accounts[0])
    console.log(err, res)
  })
})
```

## How to unlock account 2

```
web3.eth.getAccounts(function (err, accounts) {
  web3.personal.unlockAccount(accounts[1], 'btclab', function (err, res) {
    console.log(accounts[1])
    console.log(err, res)
  })
})
```

## Hot so sign a message via `ethconsole`

```
var msg = '0xfe349692aa98980aa748afbe9f480e5f05cbe7fa09fd3fb365e4c8ae004211ac'
// var msg = 'fe349692aa98980aa748afbe9f480e5f05cbe7fa09fd3fb365e4c8ae004211ac'
web3.eth.getAccounts(function (err, accounts) {
  var account0 = accounts[0]
  web3.personal.unlockAccount(account0, 'btclab', function (err, res) {
    console.log(account0, res ? 'unlocked' : 'still locked')
    web3.eth.sign(account0, msg, function (e, sign) {
      console.log(sign)
    })
  })
})
```

