const http = require('http')
const { PORT } = process.env

http.createServer((req, res) => {
  res.end('Hello World')
}).listen(PORT)
