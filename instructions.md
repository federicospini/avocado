# Avocado

>  run cli, not use npm run!

## SET ENV VARIABLES

```
ALICE_ADDR=0x0077cdf65cAD0b5ca9C2af69CD6743f396f616F0
BOB_ADDR=0x00798eaa5c12FD38ff58aD51c985D805C02eB429
ALICE_URL=http://localhost:4021
BOB_URL=http://localhost:4022
CHANNEL_ID=0x0000000000000000000000000000000000000000000000000000000000000010
```

## TEST CLI

```
node offchain/cli.js test_cli --myAddress=$ALICE_ADDR --myUrl=$ALICE_URL --counterpartyAddress=$BOB_ADDR --counterpartyUrl=$BOB_URL
```

## PROPOSE CHANNEL

Alice proposes a new channel to Bob

```
node offchain/cli.js propose_channel --myAddress=$ALICE_ADDR --myUrl=$ALICE_URL --counterpartyAddress=$BOB_ADDR --counterpartyUrl=$BOB_URL --channelId=$CHANNEL_ID --state=0x10 --challengePeriod=80
```

## VIEW PROPOSED CHANNEL

Bob checks channels proposed to him

```
node offchain/cli.js view_proposed_channels --myUrl=$BOB_URL
```

## ACCEPT PROPOSED CHANNEL

Bob accepts a channel proposed to him

```
node offchain/cli.js accept_proposed_channel --myUrl=$BOB_URL --channelId=$CHANNEL_ID
```

## PROPOSE UPDATE

Bob proposes a state update

```
node offchain/cli.js propose_update --myUrl=$BOB_URL --channelId=$CHANNEL_ID --state=0x1100
```

## VIEW PROPOSED UPDATES

Alice views state updates proposed by Bob

```
node offchain/cli.js view_their_proposed_updates --myUrl=$ALICE_URL --channelId=$CHANNEL_ID
```

Bob views state updates proposed by himself

```
node offchain/cli.js view_my_proposed_updates --myUrl=$BOB_URL --channelId=$CHANNEL_ID
```

## ACCEPT UPDATE

Alice accepts a state update proposed by Bob

```
node offchain/cli.js accept_update --myUrl=$ALICE_URL --channelId=$CHANNEL_ID --index=0
```

## VIEW ACCEPTED UPDATES

Alice views accepted states

```
node offchain/cli.js view_accepted_updates --myUrl=$ALICE_URL --channelId=$CHANNEL_ID
```

Bob views accepted states

```
node offchain/cli.js view_accepted_updates --myUrl=$BOB_URL --channelId=$CHANNEL_ID
```

## POST LAST ACCEPTES UPDATE ON BLOCKCHAIN

Bob post last accepted update on the blockchain

```
node offchain/cli.js post_last_update --myUrl=$BOB_URL --channelId=$CHANNEL_ID
```

## START CHALLENGE PERIOD

Alice starts the challenge period

```
node offchain/cli.js start_challenge_period --myUrl=$ALICE_URL --channelId=$CHANNEL_ID
```

## TRY TO CLOSE 

Bob tries to close the channel

```
node offchain/cli.js try_close --myUrl=$BOB_URL --channelId=$CHANNEL_ID
```

- - -

## TRY TO POST A PREVIOUSLY ACCEPTED UPDATE ON BLOCKCHAIN

Alice post a previously accepted update on the blockchain

```
node offchain/cli.js post_update --myUrl=$ALICE_URL --channelId=$CHANNEL_ID --index=0
```

## GET A CHANNEL FROM THE BLOCKCHAIN

Alice gets a channel from the blockchain

```
node offchain/cli.js get_blockchain_channel --myUrl=$ALICE_URL --channelId=$CHANNEL_ID
```

- - -

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

