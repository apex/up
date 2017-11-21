
v0.3.8 / 2017-11-30
===================

  * add {pre,post}{build,deploy} hooks
  * add flushing of logs after [re]start. See #359
  * add "w" for week to `ParseDuration()`
  * refactor: fix Map for now
  * refactor: use effective domain for CFN id
  * refactor: add test for existing zone and apex domain
  * refactor: add test for existing zone
  * refactor: add test coverage for CFN resources
  * fix hosted zones for sub-domains. Closes #447
  * fix `.type` precedence when runtime files are detected. Closes #436

v0.3.7 / 2017-11-24
===================

  * add date formatting for older logs
  * remove project init from `up account login`
  * fix timestamps for lambda plain text logs

v0.3.6 / 2017-11-22
===================

  * fix subscription without coupon

v0.3.5 / 2017-11-21
===================

  * add `stage` field to all log contexts (fixes log filtering against `production`)
  * fix DNS record logical id collision. Closes #420
  * refactor `up stack` output

v0.1.6-pro / 2017-11-20
=======================

  * add initial alerting support

v0.1.5-pro / 2017-11-20
=======================

  * fix "development" env support for `up env`
  * Rebase

v0.3.4 / 2017-11-20
===================

  * add `up accounts ci` and --copy to help with setting up UP_CONFIG for CI
  * fix domain verification for ssl certificates. Closes #425
  * update tj/kingpin for arg output formatting fix

v0.1.4-pro / 2017-11-18
=======================

  * Rebase

v0.3.3 / 2017-11-18
===================

  * fix zip paths on Windows. Closes #418

v0.1.3-pro / 2017-11-18
=======================

  * Rebase

v0.3.2 / 2017-11-18
===================

  * add support for UP_CONFIG from environment
  * add `up docs` command back for opening documentation in the browser
  * change logs `--since` default to 1 day
  * fix intermittent metrics failure. Closes #414

v0.3.1 / 2017-11-15
===================

  * add `up account` and sub-commands
  * add extended duration parsing for `--since` flags. Closes #401
  * add log expansion. Closes #399
  * add Content-Length request header
  * add request logs
  * add pom.xml and build.grade to whitelist which cannot be ignored
  * change metrics `--since` default to 1 month
  * refactor: remove .size defaulting of 0
  * refactor progress bar with diffing, making it more responsive
  * fix missing logs when json does not take the shape of a log. Closes #411

v0.1.2-pro / 2017-11-15
=======================

  * fix missing logs when json does not take the shape of a log. Closes #411

v0.1.0-pro / 2017-11-15
=======================

  * add `env` command

v0.3.0 / 2017-10-19
===================

  * add listing of NS records in `up stack` output
  * add changelog exclusion of docs: for goreleaser
  * add nicer domain registration form
  * update tj/survey for color changes
  * update dependencies
  * refactor: add more properties to deploy track
  * refactor: tweak cert email output
  * refactor: exclude Makefile from todo target (#382)
  * refactor: add stack to ResourceType mapping
  * refactor reporting for aws types
  * fix install.deps target
  * fix case where improper cert is created due to second-level domain (.co.uk). Closes #350
  * fix hosted zone regression introduced by e8a33a3
  * fix permission issues for static file serving. Closes #385
  * docs: add domains command
  * docs: move policy behind a details element for collapsing
  * docs: tweak for domain changes

v0.2.10 / 2017-10-13
====================

  * add flushing of proxy logs after response. Closes #370
  * add periodic flushing of proxy logs for `up start`. Closes #369
  * add internal text handler to `up start`

v0.2.9 / 2017-10-10
===================

  * fix: disable relay keep alive conns, they interact poorly with suspension (#365)

v0.2.8 / 2017-10-09
===================

  * fix missing body regression

v0.2.7 / 2017-10-09
===================

  * update go-apex dep
  * update lambda shim with concurrency support
  * fix: implement proxy GetBody to allow for re-reading request bodies. Closes #363
  * remove .lambda.timeout, replace with .proxy.timeout

v0.2.6 / 2017-09-29
===================

  * add `proxy.retry` option defaulting to `true`
  * add UP_STAGE to `up start`
  * add stage `.path` basepath support
  * fix install script for Yosemite. Closes #345

v0.2.5 / 2017-09-20
===================

  * add more relay logs
  * docs: refactor
  * add .proxy.timeout for requests and retries. Closes #335
  * refactor: remove a duplicate test
  * add retrying of 5xx errors for idempotent requests. Closes #214
  * docs: change chown to bin only. Closes #337
  * docs: add deletion info
  * docs: add more stage info
  * docs: add guide for full app
  * docs: add note about CF provisioning
  * docs: add stage section
  * docs: refactor dns section
  * docs: remove references to `certs`
  * docs: remove "coming soon"
  * docs: tweak faq
  * docs: add vendor mention
  * update Bowery/prompt dep and fix spacing
  * fix 404 checksum not found (#331)
  * docs: add missing package comments
  * docs: add missing package comments
  * docs: add note about omitting proxy bin changes

v0.2.4 / 2017-09-15
===================

  * add custom domain support
  * add Up version to the -v debug output
  * add support for JSON log lines, captured and translated to the internal format
  * add support for indented log lines to be captured as a single message
  * add sub-process cleanup and grace period. Closes #311
  * add `ssm:GetParametersByPath` to the function policy
  * add UP_STAGE env var. Closes #200
  * change default `proxy.listen_timeout` to 15
  * fix gzip handling when previously compressed. Closes #328
  * fix ignoring of .pypath

v0.2.3 / 2017-09-05
===================

  * fix rewrite content-type. Closes #304

v0.2.2 / 2017-09-05
===================

  * add logging of log query for debugging
  * add stage shorthands to log grammar. Closes #286
  * add bytes / duration units to logging grammar. Closes #283
  * add humanization of .size field in logs. Closes #252
  * add support for checking domain availability and registration. Closes #159
  * add support for multiple hook commands with arrays. Closes #127
  * add forced inclusion of ./server
  * add eu-west-2 to the regions list. Closes #280
  * fix ignoring of node_modules dotfiles (removed .bin by accident etc)
  * fix stage validation, move before building zip
  * fix support for other authentication schemes. Closes #287
  * fix dns record .ttl default
  * rename .proxy.timeout to .proxy.listen_timeout (BREAKING)
  * remove `docs` command
  * remove omission of stage from logs

v0.2.1 / 2017-08-25
===================

  * fix missing param in Infof log call, outputting `MISSING`

v0.2.0 / 2017-08-25
===================

  * add hiding of cursor for stack delete and apply
  * add support for configuring proxy timeout (#273)
  * add cost to metrics output. Closes #204
  * add: ignore dotfiles by default
  * add nicer formatting for numeric metrics
  * add build command. Closes #257
  * add validation of stage name to `url` and `deploy`. Closes #261
  * remove .npmignore support. Closes #270

v0.1.12 / 2017-08-23
=====================

  * add some basic formatting to `up stack plan`
  * rename `up stack show` to `up stack status`
  * fix hard-coded versions for stack updates

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
