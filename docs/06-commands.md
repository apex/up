---
title: Commands
---


Up provides the `up` command-line program, used to deploy the app, and manage associated resources such as domains and SSL certificates, as well as operational tasks like viewing logs.

To view details for a command at any time use `up help` or `up help <command>`.

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

  help            Show help for a command.
  build           Build zip file.
  config          Show configuration after defaults and validation.
  deploy          Deploy the project.
  domains list    List purchased domains.
  domains check   Check availability of a domain.
  domains buy     Purchase a domain.
  logs            Show log output.
  metrics         Show project metrics.
  run             Run a hook.
  stack plan      Plan configuration changes.
  stack apply     Apply configuration changes.
  stack delete    Delete configured resources.
  stack status    Show status of resources.
  start           Start development server.
  upgrade         Install the latest release of Up.
  url             Show, open, or copy a stage endpoint.
  version         Show version.
```

## Deploy

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
```

### Examples

Deploy the project to the development stage.

```
$ up
```

Deploy the project to the development stage, this is the same as running `up` without arguments.

```
$ up deploy
```

Deploy the project to the staging stage.

```
$ up deploy staging
```

Deploy the project to the production stage.

```
$ up deploy production
```

## Config

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
  "dns": {
    "zones": null
  }
}
...
```

## Logs

Show or tail log output with optional query for filtering. When viewing or tailing logs, you are viewing them from _all_ stages, see the examples below to filter on a stage name.

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
```

### JSON Output

When stdout is not a terminal Up will output the logs as JSON, which can be useful for further processing with tools such as [jq](https://stedolan.github.io/jq/).

In this contrived example the last 5 hours of production errors are piped to `jq` to produce a CSV of HTTP methods to IP address.

```
$ up logs -s 5h 'production error' | jq -r '.|[.fields.method,.fields.ip]|@csv'
```

Yielding:

```
"GET","207.194.34.24"
"GET","207.194.34.24"
"GET","207.194.34.24"
```

### Examples

Show logs from the past 5 minutes.

```
$ up logs
```

Show logs from the past 45 minutes.

```
$ up -s 45m logs
```

Show logs from the past 24 hours.

```
$ up -s 24h logs
```

Show live log output.

```
$ up logs -f
```

Show live logs from production only.

```
$ up logs -f production
```

Show live error logs from production only.

```
$ up logs -f 'production error'
```

Show error logs, which include 5xx responses.

```
$ up logs error
```

Show error and warning logs, which include 4xx and 5xx responses.

```
$ up logs 'warn or error'
```

Show logs with a specific message.

```
$ up logs 'message = "user login"'
```

Show responses with latency above 15ms.

```
$ up logs 'duration > 15'
```

Show 4xx and 5xx responses in production

```
$ up logs 'production (warn or error)'
```

Show production 5xx responses with a POST, PUT, or DELETE method.

```
$ up logs 'production error method in ("POST", "PUT", "DELETE")
```

Show 200 responses with latency above 1500ms.

```
$ up logs 'status = 200 duration > 1.5s'
```

Show responses with bodies larger than 100kb.

```
$ up logs 'size > 100kb'
```

Show 4xx and 5xx responses.

```
$ up logs 'status >= 400'
```

Show emails containing @apex.sh.

```
$ up logs 'user.email contains "@apex.sh"'
```

Show emails ending with @apex.sh.

```
$ up logs 'user.email = "*@apex.sh"'
```

Show emails starting with tj@.

```
$ up logs 'user.email = "tj@*"'
```

Show logs with a more complex query.

```
$ up logs 'method in ("POST", "PUT") ip = "207.*" status = 200 duration >= 50'
```

## URL

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
```

### Examples

Show the development endpoint.

```
$ up url
```

Open the development endpoint in the browser.

```
$ up url --open
```

Copy the development endpoint to the clipboard.

```
$ up url --copy
```

Show the production endpoint.

```
$ up url production
```

Open the production endpoint in the browser.

```
$ up url -o production
```

Copy the production endpoint to the clipboard.

```
$ up url -c production
```

## Start

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
```

### Examples

Start development server on port 3000.

```
$ up start
```

Start development server on port 5000.

```
$ up start --address :5000
```

## Stack

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
  stack plan      Plan configuration changes.
  stack apply     Apply configuration changes.
  stack delete    Delete configured resources.
  stack status    Show status of resources.
```

### Examples

Show status of the stack resources.

```
$ up stack
```

Show resource changes.

```
$ up stack plan
```

Apply resource changes.

```
$ up stack apply
```

Delete the stack resources.

```
$ up stack delete
```


## Build

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
```

### Examples

Build archive and save to ./out.zip

```
$ up build
```

Build archive and output to file via stdout.

```
$ up build > /tmp/out.zip
```

Build archive list files by size.

```
$ up build --size
```

Build archive and list size without creating out.zip.

```
$ up build --size > /dev/null
```
