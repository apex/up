---
title: Guides
---

## Logging

This description describes how you can log from you application in a way that Up will recognize. In the future Up will support forwarding your logs to services such as Loggly, Papertrail or ELK.

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
console.log(`{ "level": "info", "message": "User signin" }`)
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

Here's a simple JavaScript logger for reference, all you need to do is output some JSON to stdout and Up will handle the rest!

```js
function log(level, message, fields = {}) {
  const entry = { level, message, fields }
  console.log(JSON.stringify(entry))
}
```

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

For example filtering on slow production responses:

```
production duration >= 1s
```

Is the same as:

```
stage = "production" duration >= 1s
```

### Severity Levels

Up provides request level logging with severity levels applied automatically, for example a 5xx response is an ERROR level, while 4xx is a WARN, and 3xx or 2xx are the INFO level.

This means that instead of using the following for showing production errors:

```
production status >= 500
```

You may use:

```
production error
```

### In Operator

The `in` operator checks for the presence of a field within the set provided. For example showing only POST, PUT and PATCH requests:

```
method in ("POST", "PUT", "PATCH")
```

### Not Operator

The `not` operator is a low-precedence negation operator, for example excluding requests with the method POST, PUT, or PATCH:

```
not method in ("POST", "PUT", "PATCH")
```

Since it is the lowest precedence operator, the following will show messages that are not "user login" or "user logout":

```
not message = "user login" or message = "user logout"
```

Effectively compiling to:

```
!(message = "user login" or message = "user logout")
```

### Units

The log grammar supports units for bytes and durations, for example showing responses larger than 56kb:

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
