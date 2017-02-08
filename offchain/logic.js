/*global Uint8Array*/

import leftPad from 'left-pad'
import p from 'es6-promisify'
import { Hex, Address, Bytes32 } from './types.js'
import t from 'tcomb'
import {inspect as inspct} from 'util'

function inspect (obj) {
  return inspct(obj, { showHidden: true, depth: null, colors: true })
}

function checkSuccess (res) {
  if (!res) {
    throw new Error('no response')
  }

  if (res.error) {
    throw '(peer error) ' + res.error
  }

  // console.log('checkSuccess:', res)

  if (!res.body.success) {
    throw new Error('peer returned unexpected value')
  }
}

export class Logic {
  constructor ({ name, storage, contract, web3, post, address, url }) {
    this.name = name
    this.storage = storage
    this.contract = contract
    this.web3 = web3
    this.post = post
    this.address = address
    this.url = url

    console.log(`${name}, address: ${address}, url: ${url}, storage: ${storage._location}`)
  }

  // Test CLI
  async testCli (params) {
    return params
  }

  // View proposed channels
  async viewProposedChannels () {
    return this.storage.getItem('proposedChannels')
  }

  // View one proposed channel
  async viewProposedChannel (channelId) {
    const proposedChannels = this.storage.getItem('proposedChannels')
    return proposedChannels[channelId]
  }

  // Get the channel from the blockchain, update the local (does it actually happen? it look like it has not been coded), and return it
  async getBlockchainChannel ({ channelId }) {
    Bytes32(channelId)
    const channel = await this.contract.getChannel.call(
      channelId
    )

    return {
      channelId: channel[0], 
      address0: channel[1],
      address1: channel[2],
      phase: channel[3],
      challengePeriod: channel[4],
      closingBlock: channel[5],
      state: channel[6],
      sequenceNumber: channel[7],
    }
  }

  // Propose a new channel and send to counterparty
  async proposeChannel (params) {

    const counterpartyUrl = t.String(params.counterpartyUrl)
    const channelId = Bytes32(params.channelId)
    const address0 = Address(params.myAddress)
    const address1 = Address(params.counterpartyAddress)
    const state = Hex(params.state)
    const challengePeriod = t.Number(params.challengePeriod)

    const fingerprint = this.solSha3(
      'newChannel',
      channelId,
      address0,
      address1,
      state,
      challengePeriod
    )

    const signature0 = await p(this.web3.eth.sign)(address0, fingerprint)

    console.log(`\nfingerprint: ${fingerprint}\n, signature: ${signature0}\n`)

    checkSuccess(await this.post(counterpartyUrl + '/add_proposed_channel', {
      channelId,
      address0,
      address1,
      state,
      challengePeriod,
      signature0,
      url: this.url
    }))

    this.storeChannel({
      channelId,
      address0,
      address1,
      state,
      challengePeriod,
      signature0,
      me: 0,
      counterpartyUrl,
      theirProposedUpdates: [],
      myProposedUpdates: [],
      acceptedUpdates: []
    })
  }



  // Called by the counterparty over the http api, gets added to the proposed channel list
  async addProposedChannel ({ channelId, address0, address1, state, challengePeriod, signature0, url }) {
    t.String(url)

    var channel = {
      channelId,
      address0,
      address1,
      state,
      challengePeriod,
      signature0,
      counterpartyUrl: url
    }

    await this.verifyChannel(channel)
    
    let proposedChannels = this.storage.getItem('proposedChannels') || {}
    proposedChannels[channel.channelId] = channel
    this.storage.setItem('proposedChannels', proposedChannels)
    console.log(`- ${this.name}#storeProposedChannel: ${this.storage._location}`)

    return { success: true }
  }

  // Get a channel from the proposed channel list and accept it
  async acceptProposedChannel (channelId) {
    channelId = channelId.channelId || channelId
    const channel = this.storage.getItem('proposedChannels')[channelId]

    if (!channel) {
      throw new Error('no channel with that id')
    }

    await this._acceptChannel(channel)
  }

  // Sign the opening tx and post it to the blockchain to open the channel
  async _acceptChannel (channel) {
    const fingerprint = await this.verifyChannel(channel)
    const signature1 = await p(this.web3.eth.sign)(channel.address1, fingerprint)

    await this.contract.newChannel(
      channel.channelId,
      channel.address0,
      channel.address1,
      channel.state,
      channel.challengePeriod,
      channel.signature0,
      signature1
    )

    this.storeChannel({
      ...channel,
      me: 1,
      theirProposedUpdates: [],
      myProposedUpdates: [],
      acceptedUpdates: []
    })
  }

  // Propose an update to a channel, sign, store, and send to counterparty
  async proposeUpdate (params) {
    const channelId = Bytes32(params.channelId)
    const state = Hex(params.state)

    const channel = this.storage.getItem('channels')[channelId]

    if (!channel) {
      throw new Error('cannot find channel')
    }

    const sequenceNumber = highestProposedSequenceNumber(channel) + 1

    const fingerprint = this.solSha3(
      'updateState',
      channelId,
      sequenceNumber,
      state
    )

    const signature = await p(this.web3.eth.sign)(
      channel['address' + channel.me],
      fingerprint
    )

    const update = {
      channelId,
      sequenceNumber,
      state,
      ['signature' + channel.me]: signature
    }

    channel.myProposedUpdates.push(update)
    this.storeChannel(channel)

    console.log(channel.counterpartyUrl + '/add_proposed_update')
    checkSuccess(await this.post(channel.counterpartyUrl + '/add_proposed_update', update))

    // A confirmation that the counterparty has received the update proposal should be placed here
  }



  // Called by the counterparty over the http api, gets verified and
  // added to the proposed update list
  async addProposedUpdate (update) {
    const channel = this.storage.getItem('channels')[update.channelId]
    await this.verifyUpdate({
      channel,
      update
    })
    if (update.sequenceNumber <= highestProposedSequenceNumber(channel)) {
      throw new Error('sequenceNumber too low')
    }

    channel.theirProposedUpdates.push(update)

    this.storeChannel(channel)

    return { success: true }
  }

  async viewMyProposedUpdates ({ channelId }) {
    return this.storage.getItem('channels')[channelId].myProposedUpdates
  }

  async viewTheirProposedUpdates ({ channelId }) {
    return this.storage.getItem('channels')[channelId].theirProposedUpdates
  }

  async viewAcceptedUpdates ({ channelId }) {
    return this.storage.getItem('channels')[channelId].acceptedUpdates
  }

  // Sign the update and send it back to the counterparty
  async acceptUpdate ({ channelId, index }) { // former method sign: async acceptUpdate (update)
    // A check on the exsistance of the channel with id `channelId` is probably needed
    const channel = this.storage.getItem('channels')[channelId]
    const update = channel.theirProposedUpdates[index]
    const fingerprint = await this.verifyUpdate({ channel, update })
    const signature = await p(this.web3.eth.sign)(
      channel['address' + channel.me],
      fingerprint
    )

    update['signature' + channel.me] = signature
    channel.acceptedUpdates.push(update)
    this.storeChannel(channel)
    await this.post(channel.counterpartyUrl + '/add_accepted_update', update)
    // an ack of the counterparty acceptance of the update should perhaps be stored here
  }


  // Accepts last update from theirProposedUpdates
  async acceptLastUpdate (channelId) {
    const channel = this.storage.getItem('channels')[channelId]
    const lastUpdate = channel.theirProposedUpdates[
      channel.theirProposedUpdates.length - 1
    ]

    await this.acceptUpdate(lastUpdate)
  }

  // Called by the counterparty over the http api, gets verified and added to the accepted update list
  async addAcceptedUpdate (update) {
    const channel = this.storage.getItem('channels')[update.channelId]

    await this.verifyUpdate({
      channel,
      update,
      checkMySignature: true
    })

    if (update.sequenceNumber <= highestAcceptedSequenceNumber(channel)) {
      throw new Error('sequenceNumber too low')
    }

    channel.acceptedUpdates.push(update)
    this.storeChannel(channel)
  }



  // Post last accepted update to the blockchain
  async postLastUpdate ({ channelId }) {
    Bytes32(channelId)

    const channel = this.storage.getItem('channels')[channelId]
    const update = channel.acceptedUpdates[channel.acceptedUpdates.length - 1]

    await this.contract.updateState(
      update.channelId,
      update.sequenceNumber,
      update.state,
      update.signature0,
      update.signature1
    )
  }

  // Post an accepted state to the blockchain (useful to simulate a misbehaving party)
  async postUpdate ({ channelId, index }) {
    Bytes32(channelId)

    const channel = this.storage.getItem('channels')[channelId]
    const update = channel.acceptedUpdates[index]

    await this.contract.updateState(
      update.channelId,
      update.sequenceNumber,
      update.state,
      update.signature0,
      update.signature1
    )
  }

  // Start the challenge period, putting channel closing into motion
  async startChallengePeriod ({ channelId }) {
    Bytes32(channelId)

    const channel = this.storage.getItem('channels')[channelId]

    const fingerprint = this.solSha3(
      'startChallengePeriod',
      channelId
    )

    const signature = await p(this.web3.eth.sign)(
      channel['address' + channel.me],
      fingerprint
    )

    await this.contract.startChallengePeriod(
      channelId,
      signature,
      channel['address' + channel.me]
    )
  }

  async tryClose ({ channelId }) {
    Bytes32(channelId)

    // TODO is this useless?
    // const channel = this.storage.getItem('channels')[channelId]

    await this.contract.tryClose(channelId)
  }

  // Gets the channels list, adds the channel, saves the channels list
  storeChannel (channel) {
    const channels = this.storage.getItem('channels') || {}
    console.log(`- ${this.name}#storeChannel: ${this.storage._location}`)
    channels[channel.channelId] = channel
    this.storage.setItem('channels', channels)
  }



  // This checks that the signature is valid
  async verifyChannel(channel) {
    const channelId = Bytes32(channel.channelId)
    const address0 = Address(channel.address0)
    const address1 = Address(channel.address1)
    const state = Hex(channel.state)
    const challengePeriod = t.Number(channel.challengePeriod)
    const signature0 = Hex(channel.signature0)

    const fingerprint = this.solSha3(
      'newChannel',
      channelId,
      address0,
      address1,
      state,
      challengePeriod
    )

    console.log('\n-------------------')
    console.log(`fingerprint: ${fingerprint}`)
    console.log(`signature: ${signature0}`)
    console.log(`address0: ${address0}`)
    console.log('-------------------\n')

    const valid = await this.contract.ecverify.call(
      fingerprint,
      signature0,
      address0
    )

    console.log(`***** IS VALID: ${valid}`)

    if (!valid) {
      throw new Error('signature0 invalid')
    }

    return fingerprint
  }



  // This checks that their signature is valid, and optionally
  // checks my signature as well
  async verifyUpdate ({channel, update, checkMySignature}) {
    const channelId = Bytes32(update.channelId)
    const state = Hex(update.state)
    const sequenceNumber = t.Number(update.sequenceNumber)

    t.maybe(t.Boolean)(checkMySignature)

    const fingerprint = this.solSha3(
      'updateState',
      channelId,
      sequenceNumber,
      state
    )

    let valid = await this.contract.ecverify.call(
      fingerprint,
      update['signature' + swap[channel.me]],
      channel['address' + swap[channel.me]]
    )

    if (!valid) {
      throw new Error('signature' + swap[channel.me] + ' invalid')
    }

    if (checkMySignature) {
      let valid = await this.contract.ecverify.call(
        fingerprint,
        update['signature' + channel.me],
        channel['address' + channel.me]
      )

      if (!valid) {
        throw new Error('signature' + channel.me + ' invalid')
      }
    }

    return fingerprint
  }



  // Polyfill to get the sha3 to work the same as in solidity
  solSha3 (...args) {
    args = args.map(arg => {
      if (typeof arg === 'string') {
        if (arg.substring(0, 2) === '0x') {
          return arg.slice(2)
        } else {
          return this.web3.toHex(arg).slice(2)
        }
      }

      if (typeof arg === 'number') {
        return leftPad((arg).toString(16), 64, 0)
      }

      else {
        return ''
      }
    })

    args = args.join('')

    return this.web3.sha3(args, { encoding: 'hex' })
  }
}

const swap = [1, 0]

function highestAcceptedSequenceNumber (channel) {
  const highestAcceptedUpdate = channel.acceptedUpdates[
    channel.acceptedUpdates.length - 1
  ]
  return highestAcceptedUpdate && highestAcceptedUpdate.sequenceNumber
}

function highestProposedSequenceNumber (channel) {
  const myHighestSequenceNumber =
  channel.myProposedUpdates.length > 0 ?
    channel.myProposedUpdates[
      channel.myProposedUpdates.length - 1
    ].sequenceNumber
  : 0

  const theirHighestSequenceNumber =
  channel.theirProposedUpdates.length > 0 ?
    channel.theirProposedUpdates[
      channel.theirProposedUpdates.length - 1
    ].sequenceNumber
  : 0

  return Math.max(
    myHighestSequenceNumber,
    theirHighestSequenceNumber
  )
}
