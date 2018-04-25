
v0.6.2-pro / 2018-04-25
=======================

  * Rebase

v0.6.2 / 2018-04-25
===================

  * add up prune `--stage` flag. Closes #647
  * add `up` to ignore whitelist by default
  * remove retries on 5xx. Closes #485
  * fix login bug preventing `--email` from overriding the active team email

v0.6.1-pro / 2018-04-16
=======================

  * Rebase

v0.6.1 / 2018-04-16
===================

  * add guard against `up stack plan` before `up`
  * add `prune` command to remove old releases from S3. Closes #322

v0.6.0-pro / 2018-04-10
=======================

  * Rebase

v0.6.0 / 2018-04-10
===================

  * add annual plan subscription option

v0.5.17-pro / 2018-04-09
========================

  * Rebase

v0.5.14 / 2018-04-09
====================

  * add start command --stage flag. Closes #639
  * fix scenario where JSON logs have invalid .level values
  * refactor: add note about running `up upgrade` after subscribing

v0.5.16-pro / 2018-04-07
========================

  * Rebase

v0.5.13 / 2018-04-07
====================

  * fix "Error: fetching git commit: " error when Git is missing from the system

v0.5.15-pro / 2018-04-03
========================

  * Rebase

v0.5.12 / 2018-04-03
====================

  * add support for defining `lambda.runtime`
  * add robots middleware (#627)
  * change default runtime to nodejs 8.10
  * refactor: remove redundant wrapping of "deploying" message

v0.5.14-pro / 2018-03-22
========================

  * fix env var case with '=' in the value

v0.5.13-pro / 2018-03-19
========================

  * Rebase

v0.5.11 / 2018-03-19
====================

  * fix: update tj/go for Git signer fix

v0.5.12-pro / 2018-03-19
========================

  * Rebase

v0.5.10 / 2018-03-19
====================

  * fix: update tj/go for Git subject fix

v0.5.11-pro / 2018-03-16
========================

  * refactor: add mapping of Alarm and Subscription for `up stack plan` output
  * refactor: add .duration to Deploys track call
  * Rebase

v0.5.9 / 2018-03-16
===================

  * add support for serving static files with dynamic applications. Closes #174

v0.5.10-pro / 2018-03-15
===================

  * add nicer `up rollback` failure message when version does not exist
  * add git sha and tag support to `up rollback`
  * add `up deploys` for listing deployments and versions
  * fix log filter relational and equality operators with strings

v0.5.8 / 2018-03-15
===================

  * fix log filter relational and equality operators with strings

v0.5.7 / 2018-03-15
===================

  * add git versioning, used for Pro rollbacks and deployment changelog. Closes #100

v0.5.9-pro / 2018-03-09
=======================

  * add stage overrides for lambda warming. Closes #615

v0.5.8-pro / 2018-03-05
=======================

  * Rebase

v0.5.6 / 2018-03-05
===================

  * add support for upgrading in-place up(1). Closes #607
  * add CI specific upgrade to avoid progress bar
  * fix: remove IsNotFound error check, masks the real issue

v0.5.7-pro / 2018-03-03
=======================

  * Rebase

v0.5.5 / 2018-03-03
===================

  * fix: improve idempotency of stack deletion
  * docs: add sns to policy (necesary for Pro's alerting)

v0.5.6-pro / 2018-03-02
=======================

  * add support for `=` delimited env vars ("FOO=bar")
  * add support for passing multiple env vars to `up env set`
  * add support for overriding envs for `up start` (`$ URL=xxx up start`)

v0.5.5-pro / 2018-03-01
=======================

  * Rebase

v0.5.4 / 2018-03-01
===================

  * add default `up start` command for Go and Crystal. Closes #581
  * add log stage field to all logs, not just request-level
  * add owner to `up team` output
  * fix `up metrics` output, should be stage-specific, not global
  * refactor: add humanized error when the stack (app) does not exist
  * refactor: add stage name to beginning of log line instead of as a field
  * refactor: add os/arch to debug logs to aid in support
  * refactor: add alias upserts when updating (merged from pro)
  * refactor: remove a redundant "deploying" error wrap
  * refactor: tweak some error messages
  * refactor: change perms of up.json to 0644. Closes #601

v0.5.4-pro / 2018-02-23
=======================

  * Rebase

v0.5.3 / 2018-02-23
===================

  * fix log flushing, make it synchronous. Closes #545
  * docs: add changelog link
  * docs: add mention of BINDIR

v0.5.3-pro / 2018-02-22
=======================

  * add 1s sleep to /_ping endpoint for improved warming concurrency accuracy  
  * add `up env get` command for fetching a value
  * Rebase

v0.5.2 / 2018-02-22
===================

  * remove unsetting of `AWS_*` vars for now, reverts #590 fix

v0.5.1 / 2018-02-22
===================

  * add function version to `up stack` output
  * change `up team ci` to output base64 encoded config
  * change UP_CONFIG to attempt base64-decode when not JSON (#594)
  * fix proxy.command overrides. Closes #597
  * fix .profile precedence. Closes #590

v0.5.2-pro / 2018-02-12
=======================

  * add active warming support
  * Rebase

v0.5.1-pro / 2018-02-08
=======================

  * add `up env` --decrypt flag for emergencies when you need to list

v0.5.0-pro / 2018-02-08
=======================

  * add nicer env var logging with masking
  * add custom stage support to `up env`
  * add message for `up env` when no vars are defined
  * fix rollbacks using -previous aliases
  * Rebase

v0.5.0 / 2018-02-08
===================

  * add custom stage support. Closes #326
  * add customer feedback option when unsubscribing
  * add `up team card change` command for updating the CC
  * remove sourcing of .gitignore. Closes #557
  * remove development as a remote stage (now local only). Closes #563
  * refactor: add separator to make log message more obvious
  * refactor: add hiding of cursor when verifying email
  * refactor retry labels below s3 uploads (improves performance)
  * refactor: add nicer output when using `up url -c`

v0.4.12-pro / 2018-02-01
========================

  * Rebase

v0.4.12 / 2018-02-01
====================

  * add -o, --open to `up start` for opening in the browser
  * add `logs.{stdout,stderr}` for configuring log levels. Closes #565
  * add `-c, --command` flag to `up start`. Closes #564
  * fix panic when .domain is missing from a stage, as it is now optional. Closes #567
  * docs: add example .upignore for static sites
  * docs: fix team members rm example. Closes #562
  * docs: add "Unable to associate certificate error" to troubleshooting
  * docs: add gin example

v0.4.11-pro / 2018-01-29
========================

  * Rebase

v0.4.11 / 2018-01-29
====================

  * add development config overrides to `up start`
  * add the ability to override .proxy.command at the stage level
  * docs: mention that the WHOIS contact emails are used
  * docs: fix link for acm validation
  * docs: tweak
  * docs: add guide for hot reloading
  * docs: remove old "Local Environment Variables" guide section
  * docs: add gin example for dev command

v0.4.10-pro / 2018-01-25
========================

  * Rebase

v0.4.10 / 2018-01-25
====================

  * refactor to use a single account/region level S3 bucket, not per-project. Closes #550
  * fix base64 encoded json when params are provided

v0.4.9-pro / 2018-01-24
=======================

  * Rebase

v0.4.9 / 2018-01-24
===================

  * revert tj/go-update, causing permission issues

v0.4.8-pro / 2018-01-24
=======================

  * fix validating after overrides

v0.4.8 / 2018-01-24
===================

  * update tj/go-update for copy instead of rename. Closes #329
  * update api client for RemoveMember() json body change
  * docs: add missing ssm to policy
  * docs: add note about 404s

v0.4.7-pro / 2018-01-19
=======================

  * add rollback support
  * fix upgrade deduplication due to version having -pro suffix

v0.4.7 / 2018-01-19
===================

  * add optimization of ACM certificate creation. Closes #452
  * add `development` Lambda alias. Closes #542
  * add start of stage overrides for config. Closes #314
  * add support for upgrading to a specific version of Up. Closes #387
  * update go-cli-analytics for disabled segment cli logging
  * refactor handler.New() to accept an http.Handler
  * refactor logging configuration, delegate isatty check etc
  * refactor: move internal logs to tj/aws
  * refactor platform integration quickly

v0.4.6-pro / 2018-01-03
=======================

  * add rollback support

v0.4.5-pro / 2018-01-03
=======================

  * add s3 acceleration
  * fix a log call in runtime

v0.4.6 / 2018-01-03
===================

  * add support for Clojure with Leiningen (#522)
  * add coupon price adjustment to `up team` output. Closes #516
  * add support for overriding NODE_ENV. Closes #505
  * add error for multiple regions, until the feature is complete
  * add Paris region
  * change `error_pages` to be enabled by default for text/html requests
  * refactor `handler.New()` to accept config
  * refactor signal handling
  * refactor: update api client
  * refactor: remove unnecessary code (#517)
  * refactor login and provide a non-error when you are already signed in
  * fix s3 buckets, should be scoped to region
  * fix output flickering before build output
  * fix: add a ! in front of build.gradle for forced inclusion (#518)

v0.4.4-pro / 2017-12-22
=======================

  * Rebase

v0.4.5 / 2017-12-22
===================

  * add new subscribe workflow
  * add team CRUD and rename `up account` to `up team`. Re #410
  * refactor: replace `kingpin.CmdClause` with `kingpin.Cmd`
  * refactor: use `time.Since` for time difference (#509)
  * refactor: add "ci" to stats so we can see how often CI is used
  * refactor: simplify start of plain reporter (#508)
  * refactor: a typo fix in http/relay (#507)
  * refactor: drop unnecessary `fmt.Sprintf` in reporter/text (#506)
  * refactor: simplify personal team check (#500)

v0.4.3-pro / 2017-12-19
=======================

  * Rebase

v0.4.4 / 2017-12-19
===================

  * fix `up stack status` scenario before a domain is mapped
  * refactor: config, simplify unmarshal json of dns. Closes #497

v0.4.2-pro / 2017-12-19
=======================

  * Rebase

v0.4.3 / 2017-12-19
===================

  * refactor: shorten s3 bucket name

v0.4.1-pro / 2017-12-19
=======================

  * remove 0.0.0 hack for pro upgrade
  * Rebase

v0.4.2 / 2017-12-19
===================

  * change to disallow uppercase characters in .name. Closes #498
  * refactor: add humanized string for the current version
  * refactor: add config/backoff.go

v0.4.1 / 2017-12-18
===================

  * fix upgrades to pro when version matches

v0.4.0-pro / 2017-12-18
=======================

  * add slack `gif` option
  * add slack alert support
  * add initialization of env vars for builds. Closes #458
  * add initialization of env vars for deployments. Closes #458
  * add initialization of env vars for `up start`. Closes #458
  * add `{alerts,actions}_count` to Deploy track
  * change missing default to `notBreaching`
  * refactor: add title casing to `up env` output

v0.4.0 / 2017-12-18
===================

  * add unquoted string literals for log queries
  * add log string sans-quote literal. Closes #461
  * add log message field equality short-hand. Closes #372
  * add CI=true check for plain text output. Re #422
  * add --format=plain for CI. Closes #422
  * add setup workflow for creating up.json and doing the initial deploy. Closes #482, #386
  * add `NODE_ENV` population by default
  * add env vars to `up start`
  * add s3 deployments. Closes #272
  * add cloudfront endpoint to `up stack` output. Closes #459
  * change logs to purple (match everything else)
  * change how expanded log mode looks
  * remove `--region` flag
  * fix upgrade messages for OSS -> Pro
  * fix clearing state in text reporter

v0.3.0-pro / 2017-12-03
=======================

  * add sms alerting support

v0.2.0-pro / 2017-12-03
=======================

  * add hosted email alerting for nicer formatting
  * change alert default `period` to 1m

v0.1.11-pro / 2017-11-30
========================

  * add support for listing secrets without last modified user name
  * fix secrets listing when user ARN is not present. Closes #433
  * refactor alerting into new resources sub-pkg
  * Rebase

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

v0.1.10-pro / 2017-11-23
========================

  * add support for listing secrets without last modified user name
  * fix secrets listing when user ARN is not present. Closes #433
  * Rebase

v0.3.6 / 2017-11-22
===================

  * fix subscription without coupon

v0.1.9-pro / 2017-11-21
=======================

  * Rebase

v0.3.5 / 2017-11-21
===================

  * add `stage` field to all log contexts (fixes log filtering against `production`)
  * fix DNS record logical id collision. Closes #420
  * refactor `up stack` output

v0.1.8-pro / 2017-11-20
=======================

  * add TreatMissingData as ignore by default

v0.1.7-pro / 2017-11-20
=======================

  * fix email alerting

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
