import express from 'express'
import bodyParser from 'body-parser'
import { Logic } from '../logic.js'

export default function (globals) {
  let app = express()
  app.use(bodyParser.json())

  const name = globals.name
  const logic = new Logic(globals)

  function handlerFactory (method) {
    return function (req, res) {

      var params = JSON.stringify(req.body, null, 2)
      console.log(`\nCLI -> request from ${name}:\n---------------------------\nmethod: ${method}\nparams: ${params}\n`)

      logic[method](req.body)
      .then(result => {
        res.send(result)
      })
      .catch(error => {
        console.log(error)
        res.status(500).send({ error: error.message })
      })
    }
  }

  var routes = {
    '/test_cli': 'testCli',
    '/view_proposed_channels': 'viewProposedChannels',
    '/propose_channel': 'proposeChannel',
    // '/_accept_channel': '_acceptChannel',
    '/accept_proposed_channel': 'acceptProposedChannel',
    '/get_blockchain_channel': 'getBlockchainChannel',


    '/propose_update': 'proposeUpdate',
    '/view_my_proposed_updates': 'viewMyProposedUpdates',
    '/view_their_proposed_updates': 'viewTheirProposedUpdates',
    '/view_accepted_updates': 'viewAcceptedUpdates',
    '/accept_update': 'acceptUpdate',
    '/accept_last_update': 'acceptLastUpdate',
    '/post_last_update': 'postLastUpdate',
    '/post_update': 'postUpdate',
    
    '/start_challenge_period': 'startChallengePeriod',

    '/try_close': 'tryClose',

    '/add_proposed_channel': 'addProposedChannel',
    '/add_proposed_update': 'addProposedUpdate',
    '/add_accepted_update': 'addAcceptedUpdate'
  }

  for (var k in routes) {
    app.post(k, handlerFactory(routes[k]))
  }

  return app
}

