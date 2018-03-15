const http = require('http')
const port = parseInt(process.env.PORT, 10)

http.createServer((req, res) => {
  res.end('Node')
}).listen(port)
