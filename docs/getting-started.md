---
title: Getting Started
---

The bare minimum Up application is a single file for the application itself. Up runs "vanilla" HTTP servers listening on the `PORT` environment variable, passed to it by Up.

For example create a new directory with the following `app.js` file:

```js
const http = require('http')
const port = ~~process.env.PORT

http.createServer((req, res) => {
  res.end('Hello World from Node.js\n')
}).listen(port)

```

Deploy it to the development stage:

```
$ up
```

Open up the URL in your browser:

```
$ up url --open
```

Or copy it to the clipboard:

```
$ up url --copy
```

Or test with curl:

```
$ curl `up url`
```

That's it! You've deployed a basic Up application. Note that the first deploy may take a minute to set up the resources required. To delete it and its resources, use the following command:

```
$ up stack delete
```

If you're not a Node.js developer here are some examples in additional languages.

For Python create `app.py`:

```python
from BaseHTTPServer import BaseHTTPRequestHandler,HTTPServer
import os

class myHandler(BaseHTTPRequestHandler):
  def do_GET(self):
  self.send_response(200)
  self.send_header('Content-type','text/html')
  self.end_headers()
  self.wfile.write("Hello World from Python\n")
  return

server = HTTPServer(('', int(os.environ['PORT'])), myHandler)
server.serve_forever()
```

For Golang create `main.go`:

```go
package main

import (
  "os"
  "fmt"
  "log"
  "net/http"
)

func main() {
  addr := ":"+os.Getenv("PORT")
  http.HandleFunc("/", hello)
  log.Fatal(http.ListenAndServe(addr, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Hello World from Go")
}
```

Finally for Crystal create `main.cr`:

```ruby
require "http/server"

port = ENV["PORT"].to_i

server = HTTP::Server.new(port) do |ctx|
  ctx.response.content_type = "text/plain"
  ctx.response.print "Hello world from Crystal"
end

server.listen
```
