---
title: Installation
slug: setup
teaser: Learn how to install Up locally on your machine and in continuous integration servers.
---

Up is distributed in a binary form and can be installed manually via the [tarball releases](https://github.com/apex/up/releases) or one of the options below. The quickest way to get `up` is to run the following command:

```
$ curl -sf https://up.apex.sh/install | sh
```

By default Up is installed to `/usr/local/bin`, to specify a directory use `BINDIR`, this can be useful in CI where you may not have access to `/usr/local/bin`. Here's an example installing to the current directory:

```
$ curl -sf https://up.apex.sh/install | BINDIR=. sh
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
$ sudo chown -R $(whoami) /usr/local/bin/
```
