
GO ?= go

# Build all files.
build:
	@echo "==> Building"
	@rm -f internal/proxy/bin/bin_assets.go
	@$(GO) generate ./...
.PHONY: build

# Install from source.
install:
	@echo "==> Installing up ${GOPATH}/bin/up"
	@$(GO) install ./...
.PHONY: install

# Install the development dependencies.
install.deps:
	@echo "==> Installing dev dependencies"
	@$(GO) get -u rsc.io/gt
	@$(GO) get -u github.com/jteeuwen/go-bindata/...
	@$(GO) get -u github.com/pointlander/peg/...
.PHONY: install.deps

# Run all tests.
test:
	@$(GO) test -timeout 2m ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test

# Run all tests in CI.
test.ci:
	@$(GO) test -timeout 2m -race ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test.ci

# Show source statistics.
cloc:
	@cloc -exclude-dir=vendor,node_modules .
.PHONY: cloc

# Release binaries to GitHub.
release:
	@echo "==> Releasing"
	@goreleaser -p 1 --rm-dist -config .goreleaser.yml
	@echo "==> Complete"
.PHONY: release

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

# Show size of imports.
size:
	@curl -sL https://gist.githubusercontent.com/tj/04e0965e23da00ca33f101e5b2ed4ed4/raw/9aa16698b2bc606cf911219ea540972edef05c4b/gistfile1.txt | bash
.PHONY: size

# Clean.
clean:
	@rm -fr dist
.PHONY: clean
