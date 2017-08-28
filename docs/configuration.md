---
title: Configuration
---

Up configuration lives in the `up.json` file, and allows you review changes made to the infrastructure before applying them. This follows industry best practices of "infrastructure as code" contrasting CRUD style sub-commands for manipulating resources.

## Name

The name of the application, which is used to name resources such as the Lambda function or API Gateway. By default this is the directory name.

```json
{
  "name": "api"
}
```

## Profile

The `profile` property is equivalent to setting `AWS_PROFILE` for referencing AWS credentials in the `~/.aws` directory. This method is ideal since it prevents accidents with environment variables.

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

- **us-east** – US East (Ohio)
- **us-east** – US East (N. Virginia)
- **us-west** – US West (N. California)
- **us-west** – US West (Oregon)
- **ap-northeast-2** – Asia Pacific (Seoul)
- **ap-south-1** – Asia Pacific (Mumbai)
- **ap-southeast** – Asia Pacific (Singapore)
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
- `timeout` – Function timeout in seconds (Default `15`, Min `1`, Max `300`)

For example:

```json
{
  "name": "api",
  "lambda": {
    "memory": 512,
    "timeout": 2
  }
}
```

View the [Lambda Pricing](https://aws.amazon.com/lambda/pricing/) page for more information regarding the `memory` setting.

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

## Static File Serving

Up ships with a robust static file server, to enable it specify the app `type` as `"static"`.

```json
{
  "type": "static"
}
```

By default the current directory (`.`) is served, however you can change this using the `dir` setting. The following configuration restricts only serving of files in `./public/*`, any attempts to read files from outside of this root directory will fail.

```son
{
  "name": "app",
  "type": "static",
  "static": {
    "dir": "public"
  }
}

```

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

```json
{
  "proxy": {
    "command": "node app.js"
  }
}
```
- `gzipCompression` - By default the proxy has gzip compression enabled.  
If your server has already a gzip compression the proxy is going to compress a
compressed response leading to content decoding errors.  
So you may disable this compression by setting this flag to `false`.  
```json
{
  "proxy": {
    "command": "node app.js",
    "gzipCompression": false
  }
}
```

## Crash Recovery

Another benefit of using Up as a reverse proxy is performing crash recovery, restarting your server upon crashing and re-attempting the request before responding to the client with an error.

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

## DNS Zones & Records

Up allows you to configure DNS zones and records. One or more zones may be provided as keys in the `dns` object ("myapp.com" here), with a number of records defined within it.

```
{
  "name": "app",
  "dns": {
    "myapp.com": [
      {
        "name": "myapp.com",
        "type": "A",
        "ttl": 300,
        "value": ["35.161.83.243"]
      },
      {
        "name": "blog.myapp.com",
        "type": "CNAME",
        "ttl": 300,
        "value": ["34.209.172.67"]
      },
      {
        "name": "api.myapp.com",
        "type": "A",
        "ttl": 300,
        "value": ["54.187.185.18"]
      }
    ]
  }
}
```

## SSL Certificates

Up can provision free SSL certificates for your domains by providing a `certs` array, each with a domain and optional additional names such as `*.myapp.com` here or `api.myapp.com`.

```son
{
  "name": "app",
  "certs": [
    {
      "domains": ["myapp.com", "*.myapp.com"]
    }
  ]
}
```

Note that upon provisioning free SSL certificates AWS will send you a verification email to ensure that you own the domain name.

You can read more about this process on the [Validate Domain Ownership](http://docs.aws.amazon.com/acm/latest/userguide/gs-acm-validate.html) page.

## Custom Domains

Custom domain mapping is not _yet_ supported. For now you can visit `/apigateway` in the AWS console, click "Custom Domain Names", and map your stage(s). Note that Cloudfront can take up to 30 minutes to provision (unfortunately). Better solutions are coming.

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
