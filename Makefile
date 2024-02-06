# Make does not offer a recursive wildcard function, so here's one:
# from https://github.com/jenkins-x-plugins/jx-gitops/blob/main/Makefile
rwildcard=$(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))

SHELL := /bin/bash
BUILD_TARGET = build
MAIN_SRC_FILE=cmd/kite/main.go

GO := go
# short commit id
REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
ORG := inhere
REPO := kite-go
# exe name
NAME := kite

ORG_REPO := $(ORG)/$(REPO)
RELEASE_ORG_REPO := $(ORG)/$(NAME)
ROOT_PACKAGE := github.com/$(ORG_REPO)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
GO_DEPENDENCIES := $(call rwildcard,pkg/,*.go) $(call rwildcard,cmd/,*.go)

BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y/%m/%d-%H:%M:%S)
CGO_ENABLED = 0

REPORTS_DIR=$(BUILD_TARGET)/reports

GOTEST := $(GO) test

# set dev version unless VERSION is explicitly set via environment
# manual set: make VERSION=1.2.3
VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' | echo 'main' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

# Build flags for setting build-specific configuration at build time - defaults to empty
#BUILD_TIME_CONFIG_FLAGS ?= ""

# Full build flags used when building binaries. Not used for test compilation/execution.
BUILDFLAGS := -ldflags \
  " -s -w -X $(ROOT_PACKAGE).Version=$(VERSION)\
		-X $(ROOT_PACKAGE).Revision=$(REV)\
		-X $(ROOT_PACKAGE).Branch=$(BRANCH)\
		-X $(ROOT_PACKAGE).BuildDate=$(BUILD_DATE)\
		-X $(ROOT_PACKAGE).GoVersion=$(GO_VERSION)\
		$(BUILD_TIME_CONFIG_FLAGS)"

# Some tests expect default values for version.*, so just use the config package values there.
TEST_BUILDFLAGS :=  -ldflags "$(BUILD_TIME_CONFIG_FLAGS)"

ifdef DEBUG
BUILDFLAGS := -gcflags "all=-N -l" $(BUILDFLAGS)
endif

ifdef PARALLEL_BUILDS
BUILDFLAGS += -p $(PARALLEL_BUILDS)
GOTEST += -p $(PARALLEL_BUILDS)
else
# -p 4 seems to work well for people
GOTEST += -p 4
endif

ifdef DISABLE_TEST_CACHING
GOTEST += -count=1
endif

TEST_PACKAGE ?= ./...
COVER_OUT:=$(REPORTS_DIR)/cover.out
COVERFLAGS=-coverprofile=$(COVER_OUT) --covermode=count --coverpkg=./...

.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo -e "Provide some quick usage commands\n"
	@echo "Commands:"
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | sort

full: check ## Build and run the tests
check: build test ## Build and run the tests
get-test-deps: ## Install test dependencies
	get install github.com/axw/gocov/gocov
	get install gopkg.in/matm/v1/gocov-html

print-version: ## Print version
	@echo $(VERSION)

build: $(GO_DEPENDENCIES) clean ## Build jx-labs binary for current OS
	go mod download
	CGO_ENABLED=$(CGO_ENABLED) $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME) $(MAIN_SRC_FILE)

install: $(GO_DEPENDENCIES) darwin ## Install the kite binary to gopath/bin
	cp -f build/kite-darwin-amd64 ${GOPATH}/bin/kite
	@ls -alh ${GOPATH}/bin/kite

install2: $(GO_DEPENDENCIES) ## Install the kite to gopath/bin, without upx compress
	go build $(BUILDFLAGS) -o $(GOPATH)/bin/kite ./cmd/kite
	@ls -alh ${GOPATH}/bin/kite

install3: install win linux cp-build-to-win ## Build for local and Linux and Windows then copy to Windows(Local dev)

cp-build-to-win: ## Cleans up dependencies
	cp -f build/kite-windows-amd64.exe /Volumes/inhere-win/tools/bin/kite.exe
	cp -f build/kite-linux-amd64 /Volumes/inhere-win/tools/bin/kite
	cp -f build/kite-darwin-amd64 /Volumes/inhere-win/tools/bin/kite-darwin-amd64

pprof-cli: ## generate pprof file and start an web-ui
	go run ./_examples/pprof-cli.go
	#go tool pprof rux_prof_data.prof
	go tool pprof -http=:8080 rux_prof_data.prof

build-all:linux linux-arm win win-arm darwin darwin-arm ## Build for Linux,OSX,Windows platform

linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-linux-amd64 $(MAIN_SRC_FILE)
	upx -6 --no-progress build/$(NAME)-linux-amd64
	chmod +x build/$(NAME)-linux-amd64

linux-arm: ## Build for Linux ARM64
	GOOS=linux GOARCH=arm $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-linux-arm $(MAIN_SRC_FILE)
	upx -6 --no-progress build/$(NAME)-linux-arm
	chmod +x build/$(NAME)-linux-arm

win: ## Build for Windows
	GOOS=windows GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-windows-amd64.exe $(MAIN_SRC_FILE)
	upx -6 --no-progress build/$(NAME)-windows-amd64.exe

win-arm: ## Build for Windows arm64
	GOOS=windows GOARCH=arm64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-windows-arm64.exe $(MAIN_SRC_FILE)
	upx -6 --no-progress build/$(NAME)-windows-arm64.exe

darwin: ## Build for OSX AMD
	GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-darwin-amd64 $(MAIN_SRC_FILE)
	#upx -6 --no-progress build/$(NAME)-darwin-amd64 # upx has bug for macos 12+
	chmod +x build/$(NAME)-darwin-amd64

darwin-arm: ## Build for OSX ARM64
	GOOS=darwin GOARCH=arm64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o build/$(NAME)-darwin-arm64 $(MAIN_SRC_FILE)
	#upx -6 --no-progress build/$(NAME)-darwin-arm64
	chmod +x build/$(NAME)-darwin-arm64

.PHONY: release
release: clean linux test

release-all: release linux win darwin

promoter:
	cd promote && go build main.go

.PHONY: goreleaser
goreleaser:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --go-version=$(GO_VERSION) --root-package=$(ROOT_PACKAGE) --version=$(VERSION) --timeout 200m

.PHONY: clean
clean: ## Clean the generated artifacts
	rm -rf build release dist

get-fmt-deps: ## Install test dependencies
	get install golang.org/x/tools/cmd/goimports

.PHONY: fmt
fmt: importfmt ## Format the code
	$(eval FORMATTED = $(shell $(GO) fmt ./...))
	@if [ "$(FORMATTED)" == "" ]; \
      	then \
      	    echo "All Go files properly formatted"; \
      	else \
      		echo "Fixed formatting for: $(FORMATTED)"; \
      	fi

.PHONY: importfmt
importfmt: get-fmt-deps
	@echo "Formatting the imports..."
	goimports -w $(GO_DEPENDENCIES)

.PHONY: lint
lint: ## Lint the code
	./hack/gofmt.sh
	./hack/linter.sh
	./hack/generate.sh
