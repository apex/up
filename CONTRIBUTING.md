
# Contributing

Before contributing to Up you'll need a few things:

- Install [Golang 1.9](https://golang.org/dl/) for that Go thing if you don't have it
- Install [GIT LFS](https://git-lfs.github.com/) for large file management
- Install [golang/dep](https://github.com/golang/dep) for dependency management

The following are optional:

- Install [pointlander/peg](https://github.com/pointlander/peg) if you're working on the log grammar
- Install [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata) if you need to bake `up-proxy` into `up`

## Setup

Grab Up:

```
$ go get github.com/apex/up
```

Change into the project:

```
$ cd $GOPATH/src/github.com/apex/up
```

Grab the dependencies:

```
$ dep ensure
```

## Testing

```
$ make test
```

## Layout

Although Up is not provided as a library, it is structured primarily as if it was for organizational purposes. The project layout is loosely:

- *.go – Primary API
- [reporter](reporter) – Event based CLI reporting
- [platform](platform) – Platform specifics (AWS Lambda, Azure, Google, etc)
- [internal](internal) – Internal utilities and lower level tooling
- [http](http) – HTTP middleware for up-proxy
- [handler](handler) – HTTP middleware aggregate, effectively the entire proxy
- [docs](docs) – Documentation used to generate the static site
- [config](config) – Configuration structures and validation for up.json
- [cmd](cmd) – Commands, where `up` is the CLI and `up-proxy` is serving requests in production

Note that this is just a first past, and the code / layout will be refactored.
