# Contributing

Before contributing to Up you'll need a few things:

- Install [Golang 1.9](https://golang.org/dl/) for that Go thing if you don't have it
- Install [GIT LFS](https://git-lfs.github.com/) for large file management
- Install [golang/dep](https://github.com/golang/dep) for dependency management

The following are optional:

- Install [pointlander/peg](https://github.com/pointlander/peg) if you're working on the log grammar
- Install [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata) if you need to bake `up-proxy` into `up`
- Or run `make install.deps` to grab these tools

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

Although Up is not provided as a library it is structured as if it was, for organizational purposes. The project layout is loosely:

- *.go – Primary API
- [reporter](reporter) – Event based CLI reporting
- [platform](platform) – Platform specifics (AWS Lambda, Azure, Google, etc)
- [internal](internal) – Internal utilities and lower level tooling
- [http](http) – HTTP middleware for up-proxy
- [handler](handler) – HTTP middleware aggregate, effectively the entire proxy
- [docs](docs) – Documentation used to generate the static site
- [config](config) – Configuration structures and validation for `up.json`
- [cmd](cmd) – Commands, where `up` is the CLI and `up-proxy` is serving requests in production

Note that this is just a first past, and the code / layout will be refactored. View [Godoc](http://godoc.org/github.com/apex/up) for more details of the internals.

## Proxy

One oddity is that the `up-proxy` is baked into `up`. Yes there's a binary within the binary :) – this is so `up` can inject the proxy before deploying your function to Lambda.

The proxy accepts AWS Lambda events from API Gateway, translates them to HTTP, and sends a request to your application, then translates it back to an event that API Gateway understands.

Reverse proxy features such as URL rewriting, gzip compression, script injection, error pages and others are also provided in `up-proxy`.

## Roadmap

Up uses GitHub issue tracking and milestones for its loose roadmap. I highly recommend installing Zenhub (https://www.zenhub.com/) as well, however I primarily organize by milestones and labels for now.

## Financial contributions

We also welcome financial contributions in full transparency on our [open collective](https://opencollective.com/apex-up).
Anyone can file an expense. If the expense makes sense for the development of the community, it will be "merged" in the ledger of our open collective by the core contributors and the person who filed the expense will be reimbursed.


## Credits

### Contributors

Thank you to all the people who have already contributed to apex-up!
<a href="graphs/contributors"><img src="https://opencollective.com/apex-up/contributors.svg?width=890" /></a>


### Backers

Thank you to all our backers! [[Become a backer](https://opencollective.com/apex-up#backer)]

<a href="https://opencollective.com/apex-up#backers" target="_blank"><img src="https://opencollective.com/apex-up/backers.svg?width=890"></a>


### Sponsors

Thank you to all our sponsors! (please ask your company to also support this open source project by [becoming a sponsor](https://opencollective.com/apex-up#sponsor))

<a href="https://opencollective.com/apex-up/sponsor/0/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/1/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/2/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/3/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/4/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/5/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/6/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/7/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/8/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/9/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/9/avatar.svg"></a>
