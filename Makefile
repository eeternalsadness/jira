# --- Variables ---
APP_NAME         := jira
PKG_PREFIX       := github.com/eeternalsadness/jira/internal/util
GO_PATH          := $(shell go env GOPATH)

# Get the latest tag (falls back to "dev" if none exists)
VERSION          := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Git commit SHA
GITHUB_SHA       := $(shell git rev-parse --short HEAD)

LDFLAGS := -X $(PKG_PREFIX).Version=$(VERSION) \
           -X $(PKG_PREFIX).GitCommitSHA=$(GITHUB_SHA)

# --- Targets ---

.PHONY: build build-multiarch clean install

build:
	go build -v -ldflags "$(LDFLAGS)" -o $(APP_NAME)

build-multiarch:
	[[ ! -d dist ]] && mkdir dist; \
	platforms="linux/amd64 linux/arm64 darwin/arm64"; \
	for plat in $$platforms; do \
		go_os=$${plat%/*}; \
		go_arch=$${plat#*/}; \
		output="jira-$${go_os}-$${go_arch}.tar.gz"; \
		GOOS=$${go_os} GOARCH=$${go_arch} go build -v -ldflags "$(LDFLAGS)" -o "dist/$(APP_NAME)"; \
		tar -czvf "dist/$$output" -C dist "$(APP_NAME)" && rm "dist/$(APP_NAME)"; \
	done

clean:
	rm dist/*

install: build
	mv $(APP_NAME) "$(GO_PATH)/bin"
