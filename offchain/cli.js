var request = require('request')
var clc = require('cli-color')
var argv = require('minimist')(process.argv.slice(2), {
  string: [
    'myAddress',
    'counterpartyAddress',
    'channelId',
    'state',
    'address0',
    'address1',
    'signature0',
    'signature1'
  ]
})

var toAddress = clc.red.bold(argv.counterpartyAddress)
var toUrl = clc.green.italic(argv.counterpartyUrl)
var fromAddress = clc.red.bold(argv.myAddress)
var fromUrl = clc.green.italic(argv.myUrl)
var method = clc.yellow.italic(argv._[0])
var params = Object.assign({}, argv)
delete params['_']
params = JSON.stringify(params, null, 2)

console.log(`\nRequest:\n--------\nfrom: ${fromAddress} - ${fromUrl}\n  to: ${toAddress} - ${toUrl}\n\nmethod: ${method}\nparams: ${params}\n`)

// TODO: check exsistence of mandatory parameters

post(argv.counterpartyUrl + '/' + argv._[0], argv, response)

function response (err, res, body) {
  console.log('\nResponse:\n---------\n', err || '', body, '\n')
}

function post(url, body, callback) {
  request.post({
    url,
    body,
    json: true,
  }, callback)
}
