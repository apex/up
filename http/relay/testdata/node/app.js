const http = require('http')
const port = parseInt(process.env.UP_PORT, 10)

http.createServer((req, res) => {
  res.end('Node')
}).listen(port)
