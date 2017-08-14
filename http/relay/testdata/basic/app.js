const http = require('http')
const port = process.env.PORT

http.createServer((req, res) => {
  if (req.url.indexOf('/echo') == 0) {
    const buffers = []
    req.on('data', b => buffers.push(b))
    req.on('end', _ => {
      const body = Buffer.concat(buffers).toString()
      res.setHeader('Content-Type', 'application/json')
      res.end(JSON.stringify({
        header: req.headers,
        url: req.url,
        body
      }, null, 2))
    })
    return
  }

  if (req.url.indexOf('/throw/random') == 0) {
    if (Math.random() > 0.75) {
      yaynode()
    }

    res.end('Hello')
    return
  }

  if (req.url.indexOf('/throw/env') == 0) {
    if (process.env.UP_RESTARTS != '2') {
      yaynode()
    }

    res.end('Hello')
    return
  }

  if (req.url.indexOf('/delay') == 0) {
    setTimeout(function(){
      res.end('Hello')
    }, Math.random() * 50000)
    return
  }

  if (req.url.indexOf('/throw') == 0) {
    yaynode()
  }

  res.setHeader('Content-Type', 'text/plain')
  res.end('Hello World')
}).listen(port)
