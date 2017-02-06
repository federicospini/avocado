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

  app.post('/test_cli', handlerFactory('testCli'))
  app.post('/view_proposed_channels', handlerFactory('viewProposedChannels'))
  app.post('/propose_channel', handlerFactory('proposeChannel'))
  app.post('/accept_channel', handlerFactory('acceptChannel'))
  app.post('/accept_proposed_channel', handlerFactory('acceptProposedChannel'))
  app.post('/propose_update', handlerFactory('proposeUpdate'))
  app.post('/accept_update', handlerFactory('acceptUpdate'))
  app.post('/accept_last_update', handlerFactory('acceptLastUpdate'))
  app.post('/post_update', handlerFactory('postUpdate'))
  app.post('/start_challenge_period', handlerFactory('startChallengePeriod'))

  app.post('/add_proposed_channel', handlerFactory('addProposedChannel'))
  app.post('/add_proposed_update', handlerFactory('addProposedUpdate'))
  app.post('/add_accepted_update', handlerFactory('addAcceptedUpdate'))

  return app
}

