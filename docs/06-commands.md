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

  help                 Show help for a command.
  account status       Status of your account.
  account login        Sign in to your account.
  account logout       Sign out of your account.
  account cards add    Add credit card.
  account cards rm     Remove credit card.
  account cards ls     List credit cards.
  account subscribe    Subscribe to the Pro plan.
  account unsubscribe  Unsubscribe from the Pro plan.
  build                Build zip file.
  config               Show configuration after defaults and validation.
  deploy               Deploy the project.
  domains ls           List purchased domains.
  domains check        Check availability of a domain.
  domains buy          Purchase a domain.
  env ls               List variables.
  env add              Add a variable.
  env rm               Remove a variable.
  logs                 Show log output.
  metrics              Show project metrics.
  run                  Run a hook.
  stack plan           Plan configuration changes.
  stack apply          Apply configuration changes.
  stack delete         Delete configured resources.
  stack status         Show status of resources.
  start                Start development server.
  upgrade              Install the latest release of Up.
  url                  Show, open, or copy a stage endpoint.
  version              Show version.
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
  -s, --since="1d"     Show logs since duration (30s, 5m, 2h, 1h30m, 3d, 1M).
  -e, --expand         Show expanded logs.

Args:

  [<query>]  Query pattern for filtering logs.
```

### Expanded Output

Use the `-e` or `--expand` flag to expand log fields:

```
$ up -e 'path = "/static/*"'

1:36:34pm INFO request
           id: 8ff53267-c33a-11e7-9685-15d48d102ae9
           ip: 70.66.179.182
       method: GET
         path: /static/3.jpg
        stage: development
      version: $LATEST

1:36:34pm INFO response
     duration: 1ms
           id: 8ff53267-c33a-11e7-9685-15d48d102ae9
           ip: 70.66.179.182
       method: GET
         path: /static/3.jpg
         size: 0 B
        stage: development
       status: 304
      version: $LATEST

1:36:34pm INFO request
           id: 8ff4bd57-c33a-11e7-bf4b-4f0d97c427c5
           ip: 70.66.179.182
       method: GET
         path: /static/1.png
        stage: development
      version: $LATEST
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

Show logs from the past day.

```
$ up logs
```

Show logs from the past 45 minutes.

```
$ up -s 45m logs
```

Show logs from the past 12 hours.

```
$ up -s 12h logs
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

## Metrics

Show project metrics and estimated cost breakdown for requests, invocation count and the time spent for Lambda invocations.

```
Usage:

  up metrics [<flags>] [<stage>]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.
  -s, --since="1M"     Show logs since duration (30s, 5m, 2h, 1h30m, 3d, 1M).

Args:

  [<stage>]  Name of the stage.
```

For example:

```
$ up metrics production -s 15d

  Requests: 13,653 ($0.01)
  Duration min: 0ms
  Duration avg: 48ms
  Duration max: 15329ms
  Duration sum: 3m6.611s ($0.00)
  Errors 4xx: 1,203
  Errors 5xx: 2
  Invocations: 12,787 ($0.00)
  Errors: 0
  Throttles: 0
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

## Domains

Manage domain names, and purchase them from AWS Route53 as the registrar.

```
Usage:

  up domains <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.

Subcommands:

  domains list    List purchased domains.
  domains check   Check availability of a domain.
  domains buy     Purchase a domain.

```

### Examples

List purchased domains.

```
$ up domains
```

Check availability of a domain.

```
$ up domains check example.com
```

Purchase a domain (with interactive form).

```
$ up domains buy
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

Show status of the stack resources and nameservers.

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


## Env

Manage encrypted environment variables. Environment variables are scoped to all stages by default, with stage-level overrides. The `up env` command is available in the Pro plan.


```
Usage:

  up env <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
  -r, --region=REGION  Override the region.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --version        Show application version.

Subcommands:

  env ls               List variables.
  env add              Add a variable.
  env rm               Remove a variable.
```

```
$ up env add DB_NAME users
$ up env add DB_USER sloth
$ up env add DB_URL stage.mydb.hosted.com
$ up env add DB_PASS passforstage
$ up env add DB_CONN "host=localhost port=5432"
```

Overriding for a stage such as production can be specified with the `-s, --stage` flag. For example `DB_NAME` and `DB_USER` would likely remain the same in production, however the `DB_URL` and `DB_PASS` would not, so we should assign them as shown here:

```
$ up env add -s production DB_URL prod.mydb.hosted.com
$ up env add -s production DB_PASS passforprod
```

Environment variables may also be plain or "clear" text, using the `-c, --clear` flag, which may be viewed in the output.

```
$ up env add -c DB_NAME users
$ up env add -c DB_USER sloth
$ up env add DB_PASS amazingpass
```

List the env vars with:

```
$ up env

all                                                   

DB_NAME  users  -  Modified 34 seconds ago by tobi
DB_PASS  -      -  Modified 22 seconds ago by tobi
DB_URL   -      -  Modified 24 seconds ago by tobi
DB_USER  sloth  -  Modified 30 seconds ago by tobi

production                                            

DB_PASS  -      -  Modified 2 seconds ago by tobi  
DB_URL   -      -  Modified 4 seconds ago by tobi  
```

Note that you can also assign descriptions with the `-d, --desc` flag. Note that `-cd` is equivalent to `-c, -d` marking the first two variables as cleartext.

```
$ up env set -cd 'MongoDB collection name' DB_NAME users
$ up env set -cd 'MongoDB name' DB_NAME users
$ up env set -d 'MongoDB address' DB_URL bar
$ up env set -d 'MongoDB password' DB_PASS foo
```

Check the output and you'll see the descriptions:

```
$ up env

all                                                                         

DB_NAME  users  MongoDB collection name  Modified 2 minutes ago by tobi  
DB_PASS  -      MongoDB password         Modified 19 seconds ago by tobi
DB_URL   -      MongoDB address          Modified 1 second ago by tobi   
DB_USER  sloth  MongoDB username         Modified 2 minutes ago by tobi  
```

Note that while changes made to the variables are effective immediately, AWS Lambda may retain idle containers with the previous values. Currently you must perform a deploy in order to receive the new values. This can be used to your advantage, as it allows you to change for example both `DB_USER` and `DB_PASS` at the same time.
