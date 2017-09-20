---
title: Installation
---

Up is distributed in a binary form and can be installed manually via the [tarball releases](https://github.com/apex/up/releases) or one of the options below.

The quickest way to get `up` is to run the following command, which installs to to `/usr/local/bin` by default.

```
$ curl -sfL https://raw.githubusercontent.com/apex/up/master/install.sh | sh
```

NPM's `up` package runs the same script as above.

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
$ sudo chown -R $(whoami) /usr/local/bin/
```
