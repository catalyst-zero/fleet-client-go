PROJECT=fleet-client-go

PARENT_PACKAGE=github.com/catalyst-zero

PROJECT_PACKAGE=$(PARENT_PACKAGE)/$(PROJECT)

BUILD_PATH := $(shell pwd)/.gobuild

D0_PATH := "$(BUILD_PATH)/src/$(PARENT_PACKAGE)"

.PHONY=clean run-test get-deps update-deps fmt run-tests examples

GOPATH := $(BUILD_PATH)

SOURCE=$(shell find -L $(BUILD_PATH) -name '*.go' | grep -v examples)

OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

# === Top Level Targets ==============================

all: get-deps examples

clean:
	rm -rf $(BUILD_PATH) examples/bin/

run-tests:
	GOPATH=$(GOPATH) go test ./...

fmt:
	gofmt -l -w .

# === CODE ============================================
get-deps: .gobuild

.gobuild:
	mkdir -p $(D0_PATH)
	cd "$(D0_PATH)" && ln -s ../../../.. $(PROJECT)

	#
	# Fetch public dependencies via `go get`
	GOPATH=$(GOPATH) go get -d -v $(PROJECT_PACKAGE)

## -- examples ---------------------------------------------
examples: examples/bin/cli examples/bin/status

examples/bin/cli: .gobuild examples/cli/cli.go $(SOURCE)
	GOPATH=$(GOPATH) go build -o examples/bin/cli examples/cli/cli.go

examples/bin/status: .gobuild examples/status/status.go $(SOURCE)
	GOPATH=$(GOPATH) go build -o examples/bin/status examples/status/status.go
