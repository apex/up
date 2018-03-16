const http = require('http')
const port = process.env.PORT

http.createServer((req, res) => {
  res.setHeader('X-Foo', 'bar')
  res.setHeader('Content-Type', 'text/plain; charset=utf-8')
  res.end('Hello World')
}).listen(port)
