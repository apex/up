
# Run all tests.
test:
	@go test -cover ./...
.PHONY: test

# Test setup for CI.
test.setup:
	@echo "==> Install dep"
	@go get github.com/golang/dep/cmd/dep
	@echo "==> Install dependencies"
	@dep ensure
.PHONY: test.setup

# Generate files.
generate:
	@go generate ./...
.PHONY: generate

cloc:
	@cloc -exclude-dir=vendor,node_modules .
.PHONY: cloc

# Release binaries to GitHub.
release:
	@goreleaser --rm-dist --config .goreleaser.yml
.PHONY: release

# Output to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

# Binary size.
size:
	@go build -o up cmd/up/main.go
	@go build -o up-proxy cmd/up-proxy/main.go
	@du -h up up-proxy
	@$(MAKE) clean
.PHONY: size

# Clean.
clean:
	@rm -f up up-proxy
.PHONY: clean
