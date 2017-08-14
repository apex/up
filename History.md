
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
