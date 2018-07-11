const http = require('http');
const url = require('url');
const qs = require('querystring');
const port = process.env.PORT;

let server;

const routes = {};

routes['/echo'] = (req, res) => {
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
  });
};

routes['/timeout'] = (req, res) => {
  setTimeout(function(){
    res.end('Hello')
  }, 50000);
};

routes['/throw'] = (req, res) => {
  yaynode()
};

routes['/exit'] = (req, res) => {
  process.exit()
};

server = http.createServer((req, res) => {
  const r = Object.keys(routes).find(pattern => req.url.indexOf(pattern) === 0);
  const handler = r && routes[r];
  if (handler) {
    handler(req, res);
    return;
  }

  res.setHeader('Content-Type', 'text/plain')
  res.end('Hello World')
}).listen(port);
