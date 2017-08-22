
v0.1.11 / 2017-08-22
====================

  * add support for regions defined in `~/.aws/config`
  * add `up stack plan` and `up stack apply` support. Closes #115
  * add environment variables to hooks when performing builds etc
  * fix support for implicit `app.js` when `package.json` is present without a `start` script defined

v0.1.10 / 2017-08-15
====================

  * add default of ./server back for when source is omitted (main.go for example)
  * add `**` .upignore support
  * add forced inclusion of Up's required files
  * add support for omitting `node_modules` when using Browserify or Webpack
  * update go-archive for gitignore parity improvements

v0.1.9 / 2017-08-14
===================

  * add -modtime 0
  * add smaller progress bar for initial stack
  * revert "add error when a dir does not look like a valid project. Closes #197"
    * caused an issue if you ignore *.go for example, not robust enough

v0.1.8 / 2017-08-14
===================

  * add error when a dir does not look like a valid project. Closes #197
  * add convenience make targets `install` and `install.deps`
  * add note about AWS_PROFILE in getting started. Closes #230
  * add python projects with a requirements.txt
  * add install.sh
  * fix greedy default error page, add option to explicitly enable. Closes #233
  * fix exec bit on windows. Closes #225
  * fix python overriding of custom command
  * remove default of ./server
  * remove "-api" suffix from IAM role (breaking change)
  * refactor NewLogs() to properly delegate the error instead of panic

v0.1.7 / 2017-08-12
===================

  * add size of code/zip before attempting deploy. Closes #222
  * add better description for --force
  * change default timeout to 15s from 5s
  * change default memory from 128 to 512 (Node.js require() is slow)
  * fix relay timeout (lack of an error)
