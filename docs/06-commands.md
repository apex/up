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
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
      --version        Show application version.

Commands:

  help                 Show help for a command.
  build                Build zip file.
  config               Show configuration after defaults and validation.
  deploy               Deploy the project.
  docs                 Open documentation website in the browser.
  domains ls           List purchased domains.
  domains check        Check availability of a domain.
  domains buy          Purchase a domain.
  env ls               List variables.
  env add              Add a variable.
  env rm               Remove a variable.
  logs                 Show log output.
  metrics              Show project metrics.
  rollback             Rollback to a previous deployment.
  run                  Run a hook.
  stack plan           Plan configuration changes.
  stack apply          Apply configuration changes.
  stack delete         Delete configured resources.
  stack status         Show status of resources.
  start                Start development server.
  team status          Status of your account.
  team switch          Switch active team.
  team login           Sign in to your account.
  team logout          Sign out of your account.
  team members add     Add invites a team member.
  team members rm      Remove a member or invite.
  team members ls      List team members and invites.
  team subscribe       Subscribe to the Pro plan.
  team unsubscribe     Unsubscribe from the Pro plan.
  team card change     Change the default card.
  team ci              Credentials for CI.
  team add             Add a new team.
  upgrade              Install the latest or specified version of Up.
  url                  Show, open, or copy a stage endpoint.
  version              Show version.

Examples:

  Deploy the project to the staging environment.
  $ up

  Deploy the project to the production stage.
  $ up deploy production

  Show the staging endpoint url.
  $ up url

  Tail project logs.
  $ up logs -f

  Show error or fatal level logs.
  $ up logs 'error or fatal'

  Run build command manually.
  $ up run build

  Show help and examples for a command.
  $ up help team

  Show help and examples for a sub-command.
  $ up help team members
```

## Deploy

Deploy the project, by default to the "staging" stage. Note that running `up` and `up deploy` are identical, as it is the default command.

```
Usage:

  up deploy [<stage>]

Flags:

  -h, --help           Output usage information.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
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

Note that since `deploy` is the default command the following are also valid:

```
$ up staging
$ up production
```

## Rollback

Rollback to a previous deployment. Note that rollbacks do not affect environment variables, as they are loaded at runtime.

```
Usage:

  up rollback [<flags>] [<version>]

Flags:

  -h, --help                Output usage information.
  -C, --chdir="."           Change working directory.
  -v, --verbose             Enable verbose log output.
      --format="text"       Output formatter.
      --version             Show application version.
  -s, --stage="production"  Target stage name.

Args:

  [<version>]  Target version for rollback.
```

### Examples

Rollback to the previous production version.

```
$ up rollback
```

Rollback to an explicit production version.

```
$ up rollback 50
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
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
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
        stage: production
      version: 5

1:36:34pm INFO response
     duration: 1ms
           id: 8ff53267-c33a-11e7-9685-15d48d102ae9
           ip: 70.66.179.182
       method: GET
         path: /static/3.jpg
         size: 0 B
        stage: production
       status: 304
      version: 5

1:36:34pm INFO request
           id: 8ff4bd57-c33a-11e7-bf4b-4f0d97c427c5
           ip: 70.66.179.182
       method: GET
         path: /static/1.png
        stage: production
      version: 5
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
$ up logs 'message = "user login" method = "GET"'
```

Show logs with a specific message with implicit `=`:

```
$ up logs '"user login" method = "GET"'
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
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
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
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
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

See [Stage Overrides](https://up.docs.apex.sh/#configuration.stage_overrides) for an example of overriding the proxy command per-stage, especially useful in development.

Up Pro supports environment variables, and these will be loaded with `up start`, and variables mapped to the "development" stage will take precedence. For example:

```
$ up env set NAME Tobi
$ up start # app has NAME available as Tobi

$ up env set NAME Loki -s development
$ up start # app has NAME available Loki
```

The `UP_STAGE` and `NODE_ENV` environment variables will be set to "development" automatically.

```
Usage:

  up start [<flags>]

Flags:

  -h, --help             Output usage information.
  -C, --chdir="."        Change working directory.
  -v, --verbose          Enable verbose log output.
      --format="text"    Output formatter.
      --version          Show application version.
  -c, --command=COMMAND  Proxy command override
  -o, --open             Open endpoint in the browser.
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

Override proxy command. Note that the server created must listen on `PORT`, which is why `--port $PORT` is required for the [gin](https://github.com/codegangsta/gin) example.

```
$ up start -c 'go run main.go'
$ up start -c 'gin --port $PORT'
$ up start -c 'node --some-flag app.js'
$ up start -c 'parcel'
```

## Domains

Manage domain names, and purchase them from AWS Route53 as the registrar.

```
Usage:

  up domains <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
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

## Team

Manage team members, plans, and billing.

```
Usage:

  up team <command> [<args> ...]

Flags:

  -h, --help           Output usage information.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
      --version        Show application version.

Subcommands:

  team status          Status of your account.
  team switch          Switch active team.
  team login           Sign in to your account.
  team logout          Sign out of your account.
  team members add     Add invites a team member.
  team members rm      Remove a member or invite.
  team members ls      List team members and invites.
  team subscribe       Subscribe to the Pro plan.
  team unsubscribe     Unsubscribe from the Pro plan.
  team card change     Change the default card.
  team ci              Credentials for CI.
  team add             Add a new team.
```

### Examples

Show active team and subscription status.

```
$ up team
```

Switch teams interactively.

```
$ up team switch
```

Sign in or create account with interactive prompt.

```
$ up team login
```

Sign in to a team.

```
$ up team login --email tj@example.com --team apex-software
```

Add a new team and automatically switch to the team.

```
$ up team add "Apex Software"
```

Subscribe to the Pro plan.

```
$ up team subscribe
```

Invite a team member to your active team.

```
$ up team members add asya@example.com
```

## Upgrade

Install the latest or specified version of Up. The OSS and Pro versions have independent semver, as bugfixes and features for one may not be relevant to the other.

If you're an Up Pro subscriber, `up upgrade` will _always_ install Up Pro, even when `--target` is specified, there is no need to specify that you want the Pro version.

```
Usage:

  up upgrade [<flags>]

Flags:

  -h, --help           Output usage information.
  -C, --chdir="."      Change working directory.
  -v, --verbose        Enable verbose log output.
      --format="text"  Output formatter.
      --version        Show application version.
  -t, --target=TARGET  Target version for upgrade.
```

### Examples

Upgrade to the latest version available.

```
$ up upgrade
```

Upgrade to the specified version.

```
$ up upgrade -t 0.4.4
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

Here are some examples:

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

Variables defined via `up env` are also available to the build hooks, as well as `up start`. The `UP_STAGE` and `NODE_ENV` are defined to the target stage name automatically, unless specified manually. See the [Env Static](https://github.com/apex/up-examples/tree/master/pro/env-static) example to see how it works with build hooks.
