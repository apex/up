---
title: Getting Started
---

The simplest Up application is a single file for the application itself, with zero dependencies, and an `up.json` file which requires only a `name`.

If the directory does not contain an `up.json` file, the first execution of `up` will prompt you to create it, or you can manually create an `up.json` with some preferences:

```json
{
  "name": "appname-api",
  "profile": "companyname",
  "regions": ["us-west-2"]
}
```

Up runs "vanilla" HTTP servers listening on the `PORT` environment variable, which is passed to your program by Up. For example create a new directory with the following `app.js` file:

```js
const http = require('http')
const { PORT = 3000 } = process.env

http.createServer((req, res) => {
  res.end('Hello World from Node.js\n')
}).listen(PORT)
```

Deploy it to the staging environment:

```
$ up
```

Open up the URL in your browser:

```
$ up url --open
```

Or test with curl:

```
$ curl `up url`
```

That's it! You've deployed a basic Up application. Note that the first deploy may take a minute to set up the resources required.

To view help for commands use:

```
$ up help
$ up help COMMAND
$ up help COMMAND SUBCOMMAND
```

To delete it and its resources, use the following command:

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
