---
title: Commands
---


Up provides the `up` command-line program. To view details for a command at any time use `up help` or `up help <command>`.

```
Usage:

  up [<flags>] <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.

Commands:

  help            Show help.
  config          Show configuration after defaults and validation.
  deploy          Deploy the project.
  logs            Show log output.
  run             Run a hook.
  stack           Show the status of the stack.
  stack delete    Delete configured resources.
  start           Start development server.
  url             Show, open, or copy a stage endpoint.

Examples:

  Deploy the project to the development stage.
  $ up

  Deploy the project to the prod stage.
  $ up deploy prod

  Tail project logs.
  $ up logs -f

  Show error or fatal level logs.
  $ up logs 'error or fatal'

  Show help and examples for a sub-command.
  $ up help logs

  Run build command manually.
  $ up run build

```

## deploy

Deploy the project, by default to the "development" stage. Note that running `up` and `up deploy` are identical, however for staging and production you must run `up deploy <stage>`.

```

  Usage:

    up deploy [<stage>]

  Flags:

    -h, --help           Output usage information.
    -r, --region=REGION  Override the region.
    -C, --chdir="."      Change working directory.
    -v, --verbose        Enable verbose log output.
        --version        Show application version.

  Args:

    [<stage>]  Target stage name.

  Examples:

    Deploy the project to the development stage.
    $ up deploy

    Deploy the project to the staging stage.
    $ up deploy staging

    Deploy the project to the production stage.
    $ up deploy production

```

## config

Validate and output configuration with defaults applied.

```
$ up config
```

```json
{
  "name": "app",
  "description": "",
  "type": "server",
  "headers": null,
  "redirects": null,
  "hooks": {
    "build": "GOOS=linux GOARCH=amd64 go build -o server *.go",
    "clean": "rm server"
  },
  "environment": null,
  "regions": [
    "us-west-2"
  ],
  "inject": null,
  "lambda": {
    "role": "arn:aws:iam::ACCOUNT:role/lambda_function",
    "memory": 128,
    "timeout": 5
  },
  "cors": null,
  "error_pages": {
    "dir": ".",
    "variables": null
  },
  "proxy": {
    "command": "./server",
    "backoff": {
      "min": 100,
      "max": 500,
      "factor": 2,
      "attempts": 3,
      "jitter": false
    }
  },
  "static": {
    "dir": "."
  },
  "logs": {
    "disable": false
  },
  "certs": null,
  "dns": {
    "zones": null
  }
}
...
```

## logs

Show or tail log output with optional query for filtering.

```
 Usage:

   up logs [<flags>] [<query>]

 Flags:

   -h, --help           Output usage information.
   -r, --region=REGION  Override the region.
   -C, --chdir="."      Change working directory.
   -v, --verbose        Enable verbose log output.
       --version        Show application version.
   -f, --follow         Follow or tail the live logs.

 Args:

   [<query>]  Query pattern for filtering logs.

 Examples:

   Show logs from the past 5 minutes.
   $ up logs

   Show live log output.
   $ up logs -f

   Show error logs.
   $ up logs error

   Show error and fatal logs.
   $ up logs 'error or fatal'

   Show non-info logs.
   $ up logs 'not info'

   Show logs with a specific message.
   $ up logs 'message = "user login"'

   Show 200 responses with latency above 150ms.
   $ up logs 'status = 200 duration > 150'

   Show 4xx and 5xx responses.
   $ up logs 'status >= 400'

   Show emails containing @apex.sh.
   $ up logs 'user.email contains "@apex.sh"'

   Show emails ending with @apex.sh.
   $ up logs 'user.email = "*@apex.sh"'

   Show emails starting with tj@.
   $ up logs 'user.email = "tj@*"'

   Show logs with a more complex query.
   $ up logs 'method in ("POST", "PUT") ip = "207.*" status = 200 duration >= 50'
```

## url

Show, open, or copy a stage endpoint.

```

Usage:

  up url [<flags>] [<stage>]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.
  -o, --open           Open endpoint in the browser.
  -c, --copy           Copy endpoint to the clipboard.

Args:

  [<stage>]  Name of the stage.

Examples:

  Show the development endpoint.
  $ up url

  Open the development endpoint in the browser.
  $ up url --open

  Copy the development endpoint to the clipboard.
  $ up url --copy

  Show the production endpoint.
  $ up url production

  Open the production endpoint in the browser.
  $ up url -o production

  Copy the production endpoint to the clipboard.
  $ up url -c production

```

## start

Start development server. The development server runs the same proxy that is used in production for serving, so you can test a static site or application locally with the same feature-set.

Currently `up start` does not work with cross-compiled languages such as Go or Crystal.

```
Usage:

  up start [<flags>]

Flags:

  -h, --help             Output usage information.
  -r, --region=REGION    Override the region.
  -C, --chdir="."        Change working directory.
  -v, --verbose          Enable verbose log output.
      --version          Show application version.
      --address=":3000"  Address for server.

Examples:

  Start development server on port 3000.
  $ up start

  Start development server on port 5000.
  $ up start --address :5000

```

## stack

Stack resource management. The stack is essentially all of the resources powering your app, which is configured by Up on the first deploy.

At any time if you'd like to delete the application simply run `$ up stack delete`. To view the status and potential errors use `$ up stack`.

```
Usage:

  up stack <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.

Subcommands:
  stack delete    Delete configured resources.
  stack show      Show status of resources.

Examples:

  Show status of the stack resources.
  $ up stack

  Delete the stack resources.
  $ up stack delete

```

## build

Build zip file, typically only helpful for inspecting its contents. If you're interested in seeing what files are causing bloat, use the `--size` flag to list files by size descending.

```

Usage:

  up build [<flags>]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.
      --size           Show zip contents size information.

Examples:

  Build archive and save to ./out.zip
  $ up build

  Build archive and output to file via stdout.
  $ up build > /tmp/out.zip

  Build archive list files by size.
  $ up build --size

  Build archive and list size without creating out.zip.
  $ up build --size > /dev/null

```
