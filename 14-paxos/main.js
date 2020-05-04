const net = require('net')
const _ = require('lodash')

const controlledServers = []
const controlledServersMap = {}
const constrolledServersToBeAdded = []
const client  = new net.Socket()
client.connect({
  host: '52.49.91.111',
  port: 2092
})
client.on('connect', function () {
  console.log('connected')
  for (const i in _.times(4)) {
    const subClient  = new net.Socket()
    subClient.connect({
      host: '52.49.91.111',
      port: 2092
    })
    subClient.on('data', function (d) {
      const message = d.toString()
      const serverMatch = message.match(/SERVER ID: (\d+)/)
      if (serverMatch) {
        controlledServers.push(serverMatch[1])
        constrolledServersToBeAdded.push(serverMatch[1])
        controlledServersMap[serverMatch[1]] = true
        console.log(controlledServers, serverMatch)
      }/* else {
        console.log('=======================', message)
      }*/
      })
  }
})



const getId = (function () {
  let i = 0
  return function () {
    return i++
  }
})()
const getProposalId = (function () {
  let i = 0
  return function () {
    return i++
  }
})()
let secretOwner
let messageSent = false
let subMessageSent = false
let round
let servers
let messageId = 1
client.on('data', function (d) {
  const message = d.toString()
  console.log('--', message, '----')
  const lineMessage = message.split('\n')
  const promises = []
  for (const message of lineMessage) {
    if (message.startsWith('ROUND')) {
      const serverMatch = message.match(/ROUND (\d+): 9 -> LEARN \{servers\: \[(.*)\], secret_owner: (\d+)/)
      if (serverMatch) {
        round = serverMatch[1]
        servers = serverMatch[2].split(',')
        const id = getId()
        for (const server of servers) {
          if (server == '9') {
            continue
          }
          if (!messageSent) {
            // client.write(`PREPARE {${id},9} -> ${server}\n`)
            client.write(`PREPARE {${messageId},9} -> ${server}\n`)
            subMessageSent = true
          }
        }
        secretOwner = serverMatch[3]
        // console.log('secret owner is ', secretOwner)
      }
      const promiseMatch = message.match(/(\d+) -> PROMISE \{(\d+),9\} no_proposal/)
      if (promiseMatch) {
        // console.log('promises', promiseMatch[1],promiseMatch[2])
        promises.push({server: promiseMatch[1], value: promiseMatch[2]})
        console.log('received promise')
      }
      const promiseDenialMatch = message.match(/PROMISE_DENIAL {\d+,\d+} {(\d+),(\d+)}/)
      if (promiseDenialMatch) {
        console.log('We had a denial message######')
        messageId = parseInt(promiseDenialMatch[1]) + 1
      }

    }
  }
  if (subMessageSent) {
    messageSent = true
  }
  if (promises.length >= 1) {
    const serverToRemove = _.get(_.minBy(_.filter(promises, p => p.server !== secretOwner && !(p.server in controlledServersMap)), p => parseInt(p.server)), 'server', promises[0].server)
    const maxServer = _.get(_.maxBy(_.filter(promises, p => p.server !== secretOwner && !(p.server in controlledServersMap)), p => parseInt(p.server)), 'server', promises[0].server)
    const dedupedPromises = _.uniqBy(promises, 'server')
    const newServers = servers.length === controlledServers.length + 1 ? servers : _.filter(servers, s => s !== serverToRemove)
    let newServerToAdd
    const possibleServerToAdd = _.find(constrolledServersToBeAdded, s => !_.find(servers, s1 => s1 == s))
    if (possibleServerToAdd) {
      newServers.push(possibleServerToAdd)
    }
    console.log(serverToRemove, servers, controlledServers, promises)
    const newSecretOwner = newServers.length <= 2 * controlledServers.length ? '9' : secretOwner
    const promiseValue = promises[0].value
    const serversToSend = _.map(promises, 'server').concat(controlledServers)
    for (const serverToSend of serversToSend) {
        const command = `ACCEPT {id: {${promiseValue},9}, value: {servers: [${newServers.join(',')}], secret_owner: ${newSecretOwner}}} -> ${serverToSend}\n`
        // ACCEPT {id: {<positive_int>,<positive_int>}, value: {servers: <server_list>, secret_owner: <positive_int>}} -> <dest>
        console.log(command)
        client.write(command)
    }
    messageSent = false
    subMessageSent = false
  }
})

// server \diff servers_new < 2
