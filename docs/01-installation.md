---
title: Installation
---

Up is distributed in a binary form and can be installed manually via the [tarball releases](https://github.com/apex/up/releases) or one of the options below.

The quickest way to get `up` is to run the following command, which installs to `/usr/local/bin` by default.

```
$ curl -sf https://up.apex.sh/install | sh
```

To install `up` to a specific directory, use `BINDIR`. Here's an example installing to the current directory:

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
