---
title: Configuration
---

Configuration for your app lives in the `up.json` within your project's directory. This section details each of the options available.

## Name

The name of the application, which is used to name resources such as the Lambda function or API Gateway. By default this is the directory name.

```json
{
  "name": "api"
}
```

## Profile

The `profile` property is equivalent to setting `AWS_PROFILE` for referencing AWS credentials in the `~/.aws` directory. Use of this property is preferred it prevents accidents with environment variables.

```json
{
  "profile": "someapp"
}
```

## Regions

You may specify one or more target regions for deployment using the `regions` array. Glob style patterns may be used to match region ids. By default "us-west-2" is used unless the `AWS_REGION` environment variable is defined.

A single region:

```json
{
  "regions": ["us-west-2"]
}
```

Several regions:

```json
{
  "regions": ["us-west-2", "us-east-1", "ca-central-1"]
}
```

USA and Canada only:

```json
{
  "regions": ["us-*", "ca-*"]
}
```

Western USA only:

```json
{
  "regions": ["us-west-*"]
}
```

All regions like a boss:

```json
{
  "regions": ["*"]
}
```

Currently Lambda supports the following regions:

- **us-east-2** – US East (Ohio)
- **us-east-1** – US East (N. Virginia)
- **us-west-1** – US West (N. California)
- **us-west-2** – US West (Oregon)
- **ap-northeast-2** – Asia Pacific (Seoul)
- **ap-south-1** – Asia Pacific (Mumbai)
- **ap-southeast-1** – Asia Pacific (Singapore)
- **ap-southeast-2** – Asia Pacific (Sydney)
- **ap-northeast-1** – Asia Pacific (Tokyo)
- **ca-central-1** – Canada (Central)
- **eu-central-1** – EU (Frankfurt)
- **eu-west-1** – EU (Ireland)
- **eu-west-2** – EU (London)
- **sa-east-1** – South America (São Paulo)

**WARNING**: multi-region support won't be complete until https://github.com/apex/up/issues/134 is closed.

## Lambda Settings

The following Lambda-specific settings are available:

- `role` – IAM role ARN, defaulting to the one Up creates for you
- `memory` – Function memory in mb (Default `128`, Min `128`, Max `1536`)

For example:

```json
{
  "name": "api",
  "lambda": {
    "memory": 512
  }
}
```

View the [Lambda Pricing](https://aws.amazon.com/lambda/pricing/) page for more information regarding the `memory` setting.

Lambda timeout is implied from the [Reverse Proxy](#configuration.reverse_proxy) `timeout` setting.

## Hook Scripts

Up provides "hooks" which are commands invoked at certain points within the the deployment workflow for automating builds, linting and so on. The following hooks are available:

- `build` – Run on deploy, before building the zip file
- `clean` – Run after a build to clean up artifacts

Here's an example using Browserify to bundle a Node application. Use the `-v` verbose log flag to see how long each hook takes.

```json
{
  "name": "app",
  "hooks": {
    "build": "browserify --node app.js > server.js",
    "clean": "rm server.js"
  }
}
```

Up performs runtime inference to discover what kind of application you're using, and does its best to provide helpful defaults. See the "Runtimes" section.

Multiple commands be provided by using arrays, and are run in separate shells:

```json
{
  "name": "app",
  "hooks": {
    "build": [
      "mkdir -p build",
      "cp -fr static build",
      "browserify --node index.js > build/client.js"
    ],
    "clean": "rm -fr build"
  }
}
```

## Static File Serving

Up ships with a robust static file server, to enable it specify the app `type` as `"static"`.

```json
{
  "type": "static"
}
```

By default the current directory (`.`) is served, however you can change this using the `dir` setting. The following configuration restricts only serving of files in `./public/*`, any attempts to read files from outside of this root directory will fail.

```json
{
  "name": "app",
  "type": "static",
  "static": {
    "dir": "public"
  }
}
```

Note that `static.dir` only tells Up which directory to serve – it does not exclude other files from the directory – see [Ignoring Files](#configuration.ignoring_files).

## Environment Variables

The `environment` object may be used for plain-text environment variables. Note that these are not encrypted, and are stored in up.json which is typically committed to GIT, so do not store secrets here.

```json
{
  "name": "api",
  "environment": {
    "API_FEATURE_FOO": "1",
    "API_FEATURE_BAR": "0"
  }
}
```

These become available to you via `process.env.API_FEATURES_FOO`, `os.Getenv("API_FEATURES_FOO")` or similar in your language of choice.

The following environment variables are provided by Up:

- `PORT` – port number such as "3000"
- `UP_STAGE` – stage name such as "development" or "production"

## Header Injection

The `headers` object allows you to map HTTP header fields to paths. The most specific pattern takes precedence.

Here's an example of two header fields specified for `/*` and `/*.css`:

```json
{
  "name": "app",
  "type": "static",
  "headers": {
    "/*": {
      "X-Something": "I am applied to everything"
    },
    "/*.css": {
      "X-Something-Else": "I am applied to styles"
    }
  }
}
```

Requesting `GET /` will match the first pattern, injecting `X-Something`:

```
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 200
Content-Type: text/html; charset=utf-8
Last-Modified: Fri, 21 Jul 2017 20:42:51 GMT
X-Powered-By: up
X-Something: I am applied to everything
Date: Mon, 31 Jul 2017 20:49:33 GMT
```

Requesting `GET /style.css` will match the second, more specific pattern, injecting `X-Something-Else`:

```json
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 50
Content-Type: text/css; charset=utf-8
Last-Modified: Fri, 21 Jul 2017 20:42:51 GMT
X-Powered-By: up
X-Something-Else: I am applied to styles
Date: Mon, 31 Jul 2017 20:49:35 GMT
```

## Error Pages

By default Up will serve a minimalistic error page for you if the client accepts `text/html`, and it is explicitly enabled.

You may customize the default template's `color` and optionally provide a `support_email` to allow customers to contact your support team.

```json
{
  "name": "site",
  "type": "static",
  "error_pages": {
    "enable": true,
    "variables": {
      "support_email": "support@apex.sh",
      "color": "#228ae6"
    }
  }
}
```

If you'd like to provide custom templates you may create one or more of the following files. The most specific file takes precedence.

- `error.html` – Matches any 4xx or 5xx
- `5xx.html` – Matches any 5xx error
- `4xx.html` – Matches any 4xx error
- `CODE.html` – Matches a specific code such as 404.html

Variables specified via `variables`, as well as `.StatusText` and `.StatusCode` may be used in the template.

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{.StatusText}} - {{.StatusCode}}</title>
    <link rel="stylesheet" href="/css/style.css">
  </head>
  <body>
    <h1>{{.StatusText}}</h1>
    {{with .Variables.support_email}}
      <span class="message">Please try your request again or <a href="mailto:{{.}}">contact support</a>.</span>
    {{else}}
      <span class="message">Please try your request again or contact support.</span>
    {{end}}
  </body>
</html>
```

## Script Injection

Scripts, styles, and other tags may be injected to HTML pages before the closing `</head>` tag or closing `</body>` tag.

In the following example the `<link rel="/style.css">` is injected to the head, as well as the inlining the `scripts/config.js` file. A `<script src="/app.js"></script>` is then injected into the body.


```json
{
  "name": "site",
  "type": "static",
  "inject": {
    "head": [
      {
        "type": "style",
        "value": "/style.css"
      },
      {
        "type": "inline script",
        "file": "scripts/config.js"
      }
    ],
    "body": [
      {
        "type": "script",
        "value": "/app.js"
      }
    ]
  }
}
```

Currently you may specify the following types:

- `literal` – A literal string
- `comment` – An html comment
- `style` – A style `href`
- `script` – A script `src`
- `inline style` – An inline style
- `inline script` – An inline script
- `google analytics` – Google Analytics snippet with API key
- `segment` – Segment snippet with API key

All of these require a `value`, which sets the `src`, `href`, or inline content. Optionally you can populate `value` via a `file` path to a local file on disk, this is typically more convenient for inline scripts or styles. For example:

- `{ "type": "literal", "value": "<meta name=...>" }`
- `{ "type": "comment", "value": "Just a boring comment" }`
- `{ "type": "script", "value": "/feedback.js" }`
- `{ "type": "style", "value": "/feedback.css" }`
- `{ "type": "inline script", "file": "/feedback.js" }`
- `{ "type": "inline style", "file": "/feedback.css" }`
- `{ "type": "script", "value": "var config = {};" }`
- `{ "type": "google analytics", "value": "API_KEY" }`
- `{ "type": "segment", "value": "API_KEY" }`

## Redirects and Rewrites

Up supports redirects and URL rewriting via the `redirects` object, which maps path patterns to a new location. If `status` is omitted (or 200) then it is a rewrite, otherwise it is a redirect.

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/blog": {
      "location": "https://blog.apex.sh/",
      "status": 301
    },
    "/docs/:section/guides/:guide": {
      "location": "/help/:section/:guide",
      "status": 302
    },
    "/store/*": {
      "location": "/shop/:splat"
    }
  }
}
```

In the previous example `/blog` will redirect to a different site, while `/docs/ping/guides/alerting` will redirect to `/help/ping/alerting`. Finally `/store/ferrets` and nested paths such as `/store/ferrets/tobi` will redirect to `/shop/ferrets/tobi` and so on.

A common use-case for rewrites is for SPAs or Single Page Apps, where you want to serve the `index.html` file regardless of the path. The other common requirement for SPAs is that you of course can serve scripts and styles, so by default if a file is found, it will not be rewritten to `location`.

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/*": {
      "location": "/",
      "status": 200
    }
  }
}
```

If you wish to force the rewrite regardless of a file existing, set `force` to `true` as shown here:

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/*": {
      "location": "/",
      "status": 200,
      "force": true
    }
  }
}
```

Note that more specific target paths take precedence over those which are less specific, for example `/blog` will win over and `/*`.

## Cross-Origin Resource Sharing

CORS is a mechanism which allows requests originating from a different host to make requests to your API. Several options are available to restrict this access, however to enable CORS for everyone all you need is to add:

```json
{
  "cors": {
    "enable": true
  }
}
```

Suppose you have `https://api.myapp.com`, you may want to customize `cors` to allow access only from `https://myapp.com`:

```json
{
  "cors": {
    "allowed_origins": ["https://myapp.com"],
    "allowed_methods": ["HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"],
    "allowed_headers": ["Content-Type", "Authorization"]
  }
}
```

- `allowed_origins` – A list of origins a cross-domain request can be executed from. Use `*` to allow any origin, or a wildcard such as `http://*.domain.com`. Default is `["*"]`.
- `allowed_methods` – A list of methods the client is allowed to use with cross-domain requests. Default value is simple methods (HEAD, GET, POST).
- `allowed_headers` – A list of headers the client is allowed to use with cross-domain requests. If the special `*` value is present in the list, all headers will be allowed. Default is `[]`.
- `exposed_headers` – A list of headers which are safe to expose to the API of a CORS response.
- `max_age` – A number indicating how long (in seconds) the results of a preflight request can be cached.
- `allow_credentials` – A boolean indicating whether the request can include user credentials such as cookies, HTTP authentication or client side SSL certificates. Defaults to true.

## Reverse Proxy

Up acts as a reverse proxy in front of your server, this is how CORS, redirection, script injection and other middleware style features are provided.

The following settings are available:

- `command` – Command run through the shell to start your server (Default `./server`)
  - When `package.json` is detected `npm start` is used
  - When `app.js` is detected `node app.js` is used
  - When `app.py` is detected `python app.py` is used
- `backoff` – Backoff configuration object described in "Crash Recovery"
- `retry` – Retry idempotent requests upon 5xx or server crashes. (Default `true`)
- `timeout` – Timeout in seconds per request (Default `15`, Max `25`)
- `listen_timeout` – Timeout in seconds Up will wait for your app to boot and listen on `PORT` (Default `15`, Max `25`)
- `shutdown_timeout` – Timeout in seconds Up will wait after sending a SIGINT to your server, before sending a SIGKILL (Default `15`)

```json
{
  "proxy": {
    "command": "node app.js",
    "timeout": 5,
    "listen_timeout": 20,
    "shutdown_timeout": 5
  }
}
```

Lambda's function timeout is implied from the `.proxy.timeout` setting.

### Crash Recovery

Another benefit of using Up as a reverse proxy is performing crash recovery. Up will retry idempotent requests upon failure, and upon crash it will restart your server and re-attempt before responding to the client.

By default the back-off is configured as:

- `min` – Minimum time before retrying (Default `100ms`)
- `max` – Maximum time before retrying (Default `500ms`)
- `factor` – Factor applied to each attempt (Default `2`)
- `attempts` – Attempts made before failing (Default `3`)
- `jitter` – Apply jitter (Default `false`)

A total of 3 consecutive attempts will be made before responding with an error, in the default case this will be a total of 700ms for the three attempts.

Here's an example tweaking the default behaviour:

```json
{
  "proxy": {
    "command": "node app.js",
    "backoff": {
      "min": 500,
      "max": 1500,
      "factor": 1.5,
      "attempts": 5,
      "jitter": true
    }
  }
}
```

Since Up's purpose is to proxy your http traffic, Up will treat network errors as a crash.  When Up detects this, it will allow the server to cleanly close by sending a SIGINT, it the server does not close within `proxy.shutdown_timeout` seconds, it will forcibly close it with a SIGKILL.

## DNS Zones & Records

Up allows you to configure DNS zones and records. One or more zones may be provided as keys in the `dns` object ("myapp.com" here), with a number of records defined within it.

```json
{
  "name": "gh-polls",
  "dns": {
    "gh-polls.com": [
      {
        "name": "app.gh-polls.com",
        "type": "CNAME",
        "value": ["gh-polls.netlify.com"]
      }
    ]
  }
}
```

The record `type` must be one of:

- A
- AAAA
- CNAME
- MX
- NAPTR
- NS
- PTR
- SOA
- SPF
- SRV
- TXT

## Stages & Custom Domains

Up supports per-stage configuration, such as mapping of custom domains.

By defining a stage and its `domain`, Up knows it will need to create a free SSL certificate for `gh-polls.com`, setup the DNS records, and map the domain to API Gateway.

```json
{
  "stages": {
    "production": {
      "domain": "gh-polls.com"
    }
  }
}
```

 Here's another example mapping each stage to a domain, note that Up will create a wildcard certificate, so sub-domains are fine.


```json
{
  "stages": {
    "production": {
      "domain": "api.gh-polls.com"
    },
    "staging": {
      "domain": "stage.gh-polls.com"
    },
    "development": {
      "domain": "dev.gh-polls.com"
    }
  }
}
```

You may also provide an optional base path, for example to prefix your API with `/v1`. Note that currently your application will still receive "/v1" in its request path, for example Node's `req.url` will be  "/v1/users" instead of "/users".

```json
{
  "stages": {
    "production": {
      "domain": "api.gh-polls.com",
      "path": "/v1"
    }
  }
}
```


Plan the changes via `up stack plan` and `up stack apply` to perform the changes. Note that CloudFront can take up to ~40 minutes to distribute this configuration globally, so grab a coffee while these changes are applied.

Custom stages may be supported in the future, for now there are three:

- `development` – The latest deployment
- `staging` – The latest deployment specific to "staging"
- `production` – The latest deployment specific to "production"

If you're familiar with AWS Lambda, "development" is the $LATEST version, while "staging" and "production" are aliases.

You may [purchase domains](#guides.development_to_production_workflow.purchasing_a_domain) from the command-line, or map custom domains from other registrars. Up uses Route53 to purchase domains using your AWS account credit card. See `up help domains`.

## Ignoring Files

Up supports gitignore style pattern matching for omitting files from deployment. The following files are loaded in sequence:

- .gitignore
- .upignore

An example `.upignore` to omit markdown and `.go` source files might look like this:

```
*.md
*.go
```

### Negation

By default dotfiles are ignored, if you wish to include them, you may use `!` to negate a pattern in `.upignore`:

```
!.myfile
```

Another use-case for negation is to ignore everything and explicitly include a number of files instead, to be more specific:

```
*
!app.js
!package.json
!node_modules/**
!src/**
```

### Inspecting

To get a better idea of which files are being filtered or added, use `up -v` when deploying, and you may also find it useful to `grep` in some cases:

```
$ up -v 2>&1 | grep filtered
DEBU filtered .babelrc – 25
DEBU filtered .git – 408
DEBU filtered .gitignore – 13
DEBU filtered node_modules/ansi-regex/readme.md – 1749
DEBU filtered node_modules/ansi-styles/readme.md – 1448
DEBU filtered node_modules/binary-extensions/readme.md – 751
DEBU filtered node_modules/chalk/readme.md – 6136
```

You may also wish to use `up build --size` to view the largest files within the zip.

### Node specifics

When building a Node project, `node_modules` is often added to `.gitignore`, so Up special-cases this scenario so that `node_modules` is not excluded by default.

If you're using a bundler such as Webpack or Browserify instead, you may wish to ignore `node_modules` as they're not required in production, add this to your `./.upignore`:

```
node_modules
```

### Pattern matching

Note that patterns are matched much like `.gitignore`, so if you have the following `.upignore` contents even `node_modules/debug/src/index.js` will be ignored since it contains `src`.

```
src
```

You can be more specific with a leading `./`:

```
./src
```

Files can be matched recursively using `**`, for example ignoring everything except the files in `dist`:

```
*
!dist/**
```
