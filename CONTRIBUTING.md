
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
