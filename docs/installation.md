---
title: Installation
---

Up can be installed via pre-compiled binaries, head over to the [Releases](https://github.com/apex/up/releases) page, or use this one-liner which will install `up` to `/usr/local/bin` by default.

```
$ curl -sfL https://raw.githubusercontent.com/apex/up/master/install.sh | sh
```

Or via NPM with:

```
$ npm i -g up
```

Verify installation with:

```
$ up version
```

Later when you want to update `up` to the latest version use the following command:

```
$ up upgrade
```

If you hit permission issues, you may need to run the following, as `up` is installed to `/usr/local/bin/up` by default.

```
$ sudo chown -R $(whoami) /usr/local/
```
