---
title: Runtimes
---

Up supports a number of interpreted langauges, and virtually any language which can be compiled to a binary such as Golang. Up does its best to provide idiomatic and useful out-of-the-box experiences tailored to each language. Currently first-class support is provided for:

- Python
- Golang
- Node.js
- Crystal
- Static sites

## Node.js

When a `package.json` file is detected, Node.js 6.10 is the assumed runtime.

The `build` hook becomes:

```
$ npm run build
```

The server run by the proxy becomes:

```
$ npm start
```

## Golang

When a `main.go` file is detected, Golang is the assumed runtime.

The `build` hook becomes:

```
$ GOOS=linux GOARCH=amd64 go build -o server *.go
```

The `clean` hook becomes:

```
$ rm server
```

## Crystal

When a `main.cr` file is detected, Crystal is the assumed runtime. Note that this runtime requires Docker to be installed.

The `build` hook becomes:

```
$ docker run --rm -v $(PWD):/src -w /src tjholowaychuk/up-crystal crystal build --link-flags -static -o server main.cr
```

The `clean` hook becomes:

```
$ rm server
```
