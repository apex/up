---
title: Guides
---

## Subscribing to Up Pro

Up Pro provides additional features which are not available in the open-source version, such as encrypted environment variables, alerting support and more.

First sign into the platform with the following command – you'll receive an email for confirmation.

```
$ up team login

     email: tj@apex.sh
  ⠋ verify: Check your email for a confirmation link

```

Click the link in your email and you're signed in! If you're using Up Pro with one or more organizations, you should create a team to manage team members and subscriptions independently. If you plan on using Up Pro for personal use you may skip this step.

```
$ up team add "My Company"
```

Next you'll need to subscribe! You'll be asked for an optional coupon, credit card information – which never touches our servers, only Stripe via HTTPS – and finally a subscription confirmation.

```
$ up team subscribe
 ```

Now you and your team members may upgrade to the latest version of Up Pro, instead of the open-source distribution:

```
$ up upgrade
```

To view the status of your account at any time run the following:

```
$ up team

  team: apex
  subscription: Up Pro
  amount: $10.00/mo USD
  created: December 22, 2017
```

To switch to another team run the following and select the active team.

```
$ up team switch


   ❯ apex
     tj@apex.sh
```

## Inviting Team Members

To invite members use the following command:

```
$ up team members add tobi@apex.sh
$ up team members add loki@apex.sh
$ up team members add jane@apex.sh
```

At any time you can view invites and members:

```
$ up team members

team: apex

Members

 • tj@apex.sh
 • tobi@apex.sh
 • loki@apex.sh

Invites

 • jane@apex.sh
```

Your team members will receive an email with installation instructions, where they run the following to sign in – with your team id of course.

```
$ up team login --email tobi@apex.sh --team apex
```

## Development to Production Workflow

This section guides you through taking a small application from development, to production, complete with purchasing and mapping a custom domain.

### Deploying

First create `app.js` in an empty directory with the following Node.js app. Note that it must listen on __PORT__ which is passed by Up.

```js
const http = require('http')
const { PORT = 3000 } = process.env

http.createServer((req, res) => {
  res.end('Hello World\n')
}).listen(PORT)
```

Next you should give your application a name and start configuring. The `profile` name should correspond with the name in `~/.aws/credentials` so that Up knows which AWS account to deploy to, and which credentials to use.

```json
{
  "name": "up-example",
  "profile": "up-tobi"
}
```

Run `up` to deploy the application.

```
$ up

   build: 5 files, 3.9 MB (358ms)
  deploy: complete (14.376s)
   stack: complete (1m12.086s)
```

Test with `curl` to ensure everything is working:

```
$ curl `up url`
Hello World
```

### Purchasing a Domain

Domains can be mapped from existing services, or purchased directly from AWS via Route53. First check if the domain you'd like is available:

```
$ up domains check up.com

  Domain up.com is unavailable

  Suggestions:

  theupwards.com          $12.00 USD
  upwardonline.com        $12.00 USD
  myupwards.com           $12.00 USD
  theastir.com            $12.00 USD
  astironline.com         $12.00 USD
  myastir.com             $12.00 USD
  myupward.net            $11.00 USD
  cleanup.tv              $32.00 USD
  myup.tv                 $32.00 USD
  itup.tv                 $32.00 USD
  newup.tv                $32.00 USD
  thedown.net             $11.00 USD
  theupward.net           $11.00 USD
  upwardsonline.net       $11.00 USD
```

Oh no up.com is taken! Try another:

```
$ up domains check up-example.com

  Domain up-example.com is available for $12.00 USD
```

Purchase it with the following command and fill out the details required by the registrar:

```
$ up domains buy

  Domain: up-example.com
  First name: TJ
  Last name: Holowaychuk
  Email: tj@apex.sh
  Phone: +1.2501007000
  Country code: CA
  City: Victoria
  State or province: BC
  Zip code: X9X 9X9
  Address: Some address here
```

It can take a few minutes for AWS to finalize the purchase after which you should receive an email. Then you'll see it in the `up domains` output along with the automatic renewal time.

```
$ up domains

  gh-polls.com             renews Aug 28 17:17:58
  up-example.com           renews Sep 19 19:40:50
```

By default domains purchased with Up have privacy protection enabled, hiding your contact information from [WHOIS](https://en.wikipedia.org/wiki/WHOIS).

### Deploying to Stages

Before deploying to the staging and production stages, first tweak the application a little to include the `UP_STAGE` environment variable:

```js
const http = require('http')
const { PORT = 3000, UP_STAGE } = process.env

http.createServer((req, res) => {
  res.end('Hello World from ' + UP_STAGE)
}).listen(PORT)
```

Now deploy to development and production. Note that `up` is an alias of `up deploy development`.

```
$ up
$ up deploy production
```

Open both in the browser:

```
$ up url -o
$ up url production -o
```

You should see "Hello World from production" and "Hello World from development".

### Mapping Custom Domains to Stages

Now that you have an application deployed, you probably want a fancy custom domain for it right? You can map these using the `stages` and `domain` properties.

Here we let Up know that we want `up-example.com` for production and `dev.up-example` for development.

```json
{
  "name": "up-example",
  "profile": "up-tobi",
  "stages": {
    "development": {
      "domain": "dev.up-example.com"
    },
    "production": {
      "domain": "up-example.com"
    }
  }
}
```

Note that you could map staging to a domain like `staging-myapp.com`, it does not have to be a sub-domain of your production domain.

Now when you run `up stack plan` to preview changes to your resources, it will prompt you to verify the Let's Encrypt certificate emails that AWS sends.

```
$ up stack plan

       domains: Check your email for certificate approval
     ⠧ confirm: up-example.com
```

AWS requires email verification to prove you own the domain. AWS sends an email to the 3 contact addresses listed in WHOIS when you registered the domain, and to the following 5 common system addresses for your domain:

- administrator@your_domain_name
- hostmaster@your_domain_name
- postmaster@your_domain_name
- webmaster@your_domain_name
- admin@your_domain_name

See [Validate Domain Ownership](https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-validate-email.html) for more information.

After clicking "I Approve" in one of the emails, the output will resume and you'll see some new resources Up will be creating.

```
Add AWS::ApiGateway::DomainName
  id: ApiDomainDevelopment

Add AWS::ApiGateway::BasePathMapping
  id: ApiDomainDevelopmentPathMapping

Add AWS::ApiGateway::DomainName
  id: ApiDomainProduction

Add AWS::ApiGateway::BasePathMapping
  id: ApiDomainProductionPathMapping

Add AWS::Route53::RecordSet
  id: DnsZoneDevUpExampleComRecordDevUpExampleCom

Add AWS::Route53::RecordSet
  id: DnsZoneUpExampleComRecordUpExampleCom
```

If you're curious, now that Up knows you want to map the domain(s), it will create:

- Registers ACM free SSL certificate(s) for your domain(s)
- CloudFront distribution for the API Gateway
- API Gateway stage mapping
- Route53 DNS zone and record(s) mapping to the CloudFront distribution

Now apply these changes:

```
$ up stack apply
```

After the changes have been applied, it can take roughly 10-40 minutes for CloudFront to distribute the configuration and SSL certificate globally, so until then our up-example.com domain won't work.

Once available https://up-example.com will always point to production via `up deploy production`, and https://dev.up-example.com/ will point to the latest deployment via `up`.

### Mapping Domains from External Registrars

If you purchased a domain via `up domains buy` you can skip this step, however if you used an external registrar such as Godaddy you will need to delegate to AWS for DNS management.

To do this you'll need to sign in to your registrar's site, and configure the nameservers. To figure out what values to use for the nameservers, run `up stack`, which outputs the NS records for the apex (top-level) domains of your application.

```
$ up stack

Staging

  domain: stage.up-example.com
  endpoint: d2od0udp1p8bru.cloudfront.net

Production

  domain: up-example.com
  endpoint: d72wsqljqg5cy.cloudfront.net
  nameservers:
   • ns-1495.awsdns-58.org
   • ns-103.awsdns-12.com
   • ns-1670.awsdns-16.co.uk
   • ns-659.awsdns-18.net
```

Save those four values in your registrar's interface, and you should be good to go! Note that altering DNS records can take some time to propagate.

### Mapping with Third-party DNS

If you manage DNS with a third-party such as Cloudflare, and wish to use Up only for deployment you will need to manually edit or add DNS records.

For example if your top-level domain `sloths.com` is managed by Cloudflare and you'd like point `api.sloths.com` to your app, you should first add it to your `up.json`:

```json
{
  "name": "sloths"
  "stages": {
    "production": {
      "domain": "api.sloths.com"
    }
  }
}
```

Next you will need to `up stack plan` and `up stack apply`, this will set up a CloudFront end-point for the application. To view the endpoint information, run `up stack`:

```
$ up stack

Production

  domain: api.sloths.com
  endpoint: d72wsqljqg5cy.cloudfront.net
```

In your DNS provider – Cloudflare in this example – you should create a `CNAME` record pointing to the production `endpoint`. Make sure that the `domain` you use matches the domain in Cloudflare.

### Stack Changes

The "stack" is all of the resources associated with your app. You plan changes via `up stack plan` and perform them with `up stack apply`.

Suppose you wanted to map the "staging" stage, you would first add it to `up.json`:

```json
{
  "name": "up-example",
  "profile": "up-tobi",
  "stages": {
    "development": {
      "domain": "dev.up-example.com"
    },
    "staging": {
      "domain": "stage.up-example.com"
    },
    "production": {
      "domain": "up-example.com"
    }
  }
}
```

Then run:

```
$ up stack plan
```

Review the output, it should be all "Add"s in this case, then apply:

```
$ up stack apply
```

### Deleting the App

After you're done messing around, you may want to remove all the resources and the app itself. To do so simply run:

```
$ up stack delete
```

## Deploying Applications from Continuous Integration

Up makes it easy to deploy your applications from CI, thanks to its Go binaries you can install Up in seconds in any CI provider such as Travis, Circle, Semaphore among others.

### Environment Variables

The first step is to set up environment variables so that you have access to your AWS account. You can get these values from `cat ~/.aws/credentials`:

- `AWS_ACCESS_KEY_ID` – AWS access key
- `AWS_SECRET_ACCESS_KEY` – AWS secret key

If using running Up Pro you'll need your Up credentials:

- `UP_CONFIG` – Up configuration

To obtain this run `up team ci` or `up team ci --copy` to copy it directly to your clipboard, then paste this as the env var's value.

### Commands

Now that configuration is prepared, you'll need to install Up using the following two commands:

```
$ sudo chown -R $(whoami) /usr/local/bin
$ curl -sf https://up.apex.sh/install | sh
```

If you're on Up Pro you'll need to upgrade to install the Pro binary:

```
$ up upgrade
```

After that all you need to do is deploy to your desired stage!

```
$ up deploy production
```

If your provider sets `CI=true` colors will be disabled automatically, however, you may manually specify this using `--format`:

```
$ up --format=plain deploy production
```

## Mastering Logging

This section describes how you can log from your application in a way that Up will recognize. In the future Up will support forwarding your logs to services such as Loggly, Papertrail or ELK.

### Plain Text

The first option is plain-text logs to stdout or stderr. Currently writes to stderr are considered ERROR-level logs, and stdout becomes INFO.

Writing plain-text logs is simple, for example with Node.js:

```js
console.log('User signed in')
console.error('Failed to sign in: %s', err)
```

Would be collected as:

```
 INFO: User signed in
ERROR: Failed to sign in: something broke
```

Multi-line indented logs are also supported, and are treated as a single message. For example:

```js
console.log('User signed in')
console.log('  name: %s', user.name)
console.log('  email: %s', user.email)
```

Would be collected as the single entry:

```
INFO: User signed in
  name: tj
  email: tj@apex.sh
```

This feature is especially useful for stack traces.

### JSON

The second option is structured logging with JSON events, which is preferred as it allows you to query against specific fields and treat logs like events.

JSON logs require a `level` and `message` field:

```js
console.log(`{ "level": "info", "message": "User login" }`)
```

Would be collected as:

```
INFO: User login
```

The `message` field should typically contain no dynamic content, such as user names or emails, these can be provided as fields:

```js
console.log(`{ "level": "info", "message": "User login", "fields": { "name": "Tobi", "email": "tobi@apex.sh" } }`)
```

Would be collected as:

```
INFO: User login name=Tobi email=tobi@apex.sh
```

Allowing you to perform queries such as:

```
$ up logs 'message = "User login" name = "Tobi"'
```

Or:

```
$ up logs 'name = "Tobi" or email = "tobi@*"'
```

Here's a simple JavaScript logger for reference. All you need to do is output some JSON to stdout and Up will handle the rest!

```js
function log(level, message, fields = {}) {
  const entry = { level, message, fields }
  console.log(JSON.stringify(entry))
}
```

For example, with the Go [apex/log](https://github.com/apex/log) package you'd use the `json` handler, which outputs this format.

## Log Query Language

Up supports a comprehensive query language, allowing you to perform complex filters against structured data, supporting operators, equality, substring tests and so on. This section details the options available when querying.

### AND Operator

The `and` operator is implied, and entirely optional to specify, since this is the common case.

Suppose you have the following example query to show only production errors from a the specified IP address.

```
production error ip = "207.194.32.30"
```

The parser will inject `and`, effectively compiling to:

```
production and error and ip = "207.194.38.50"
```

### Or Operator

There is of course also an `or` operator, for example showing warnings or errors.

```
production (warn or error)
```

These may of course be nested as you require:

```
(production or staging) (warn or error) method = "GET"
```

### Equality Operators

The `=` and `!=` equality operators allow you to filter on the contents of a field.

Here `=` is used to show only GET requests:

```
method = "GET"
```

Or for example `!=` may be used to show anything except GET:

```
method != "GET"
```

### Relational Operators

The `>`, `>=`, `<`, and `<=` relational operators are useful for comparing numeric values, for example response status codes:

```
status >= 200 status < 300
```

### Stages

Currently all development, staging, and production logs are all stored in the same location, however you may filter to find exactly what you need.

The keywords `production`, `staging`, and `development` expand to:

```
stage = "production"
```

For example, filtering on slow production responses:

```
production duration >= 1s
```

Is the same as:

```
stage = "production" duration >= 1s
```

### Severity Levels

Up provides request level logging with severity levels applied automatically. For example, a 5xx response is an ERROR level, while 4xx is a WARN, and 3xx or 2xx are the INFO level.

This means that instead of using the following for showing production errors:

```
production status >= 500
```

You may use:

```
production error
```

### In Operator

The `in` operator checks for the presence of a field within the set provided. For example, showing only POST, PUT and PATCH requests:

```
method in ("POST", "PUT", "PATCH")
```

### Units

The log grammar supports units for bytes and durations. For example, showing responses larger than 56kb:

```
size > 56kb
```

Or showing responses longer than 1500ms:

```
duration > 1.5s
```

Byte units are:

- `b` bytes (`123b` or `123` are equivalent)
- `kb` bytes (`5kb`, `128kb`)
- `mb` bytes (`5mb`, `15.5mb`)

Duration units are:

- `ms` milliseconds (`100ms` or `100` are equivalent)
- `s` seconds (`1.5s`, `5s`)

### Substring Matches

When filtering on strings, such as the log message, you may use the `*` character for substring matches.

For example if you want to show logs with a remote ip prefix of `207.`:

```
ip = "207.*"
```

Or a message containing the word "login":

```
message = "*login*"
```

There is also a special keyword for this case:

```
message contains "login"
```

## Hot Reloading in Development

The `up start` command uses your `proxy.command` by default, which may be inferred based on your application type, such as `node app.js` for Node.js or `./server` for Golang.

You may alter this command for `up start` with the development environment. For example with Golang you may want `go run main.go`, or hot reloading with [gin](https://github.com/codegangsta/gin) as shown here:

```json
{
  "name": "app",
  "stages": {
    "development": {
      "proxy": {
        "command": "gin --port $PORT"
      }
    }
  }
}
```

Note that the server must always listen on `PORT` which is provided by `up start`.
