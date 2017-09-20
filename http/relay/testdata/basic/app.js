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

routes['/env'] = (req, res) => {
  const query = qs.parse(url.parse(req.url).query);
  res.end(process.env[query.key], 'utf8');
};

routes['/pid'] = (req, res) => {
  res.end(String(process.pid), 'utf8');
};

routes['/throw/random'] = (req, res) => {
  if (Math.random() > 0.75) {
    yaynode();
  }

  res.end('Hello');
};

routes['/throw/env'] = (req, res) => {
  if (process.env.UP_RESTARTS != '2') {
    yaynode();
  }

  res.end('Hello');
};

routes['/delay'] = (req, res) => {
  setTimeout(function(){
    res.end('Hello')
  }, Math.random() * 50000);
};

routes['/throw'] = (req, res) => {
  yaynode()
};


// Close the server, so we can get network errors without ending the process
routes['/close'] = (req, res) => {
  server.close(() => 1);

  res.writeHead(200, {
    Connection: 'close'
  });

  res.end('closed');
};

routes['/appError'] = (req, res) => {
  res.writeHead(200, {
    Connection: 'close'
  });

  // Do some work here
  console.log('Doing some non-idempotent work');
  setTimeout(() => {
    res.socket.end();
  }, 10);
};

routes['/swallowSignals'] = (req, res) => {
  for (const s of ['SIGINT', 'SIGTERM', 'SIGHUP', 'SIGUSR1', 'SIGUSR2']) {
    process.on(s, () => 1);
  }

  res.end('swallow');
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

// Run this to prevent the server from exiting from an empty event loop
// This mimics _real_ servers a bit more exactly.
const keepAlive = setInterval(() => 1, 5000);
