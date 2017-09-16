Open an issue and discuss changes before spending time on them, unless the change is trivial or an issue already exists.

Use "VERB some thing here. Closes #n" to close the relevant issue, where VERB is one of:

  - add
  - remove
  - change
  - refactor

If the change is documentation related prefix with "docs: ", as these are filtered from the changelog.

  docs: add ~/.aws/config

Omit changes to the up-proxy binary if they are included in your PR. The proxy is re-built before releases.

Run `dep ensure` if you introduce any new `import`'s so they're included in the ./vendor dir.
