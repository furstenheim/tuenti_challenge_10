const xor = require('buffer-xor')
function decrypt () {
  const firstPlanet = Buffer.from('514;248;980;347;145;332')
  const firstMessage = Buffer.from('3633363A33353B393038383C363236333635313A353336','hex')
  const secondMessage = Buffer.from('3A3A333A333137393D39313C3C3634333431353A37363D', 'hex')
  const decrypted = xor(xor(firstMessage, firstPlanet), secondMessage).toString()
  console.log(decrypted)
}

decrypt()
