BIN          = terraform-provider-ceph
GOFMT_FILES ?= $$(find . -name '*.go')
GO_ARGS     ?=

all: build

$(BIN): ceph main.go go.mod go.sum
	go build $(GO_ARGS) -o $@

fmt:
	go generate
	gofmt -s -w $(GOFMT_FILES)

build: $(BIN)

debug: GO_ARGS += -gcflags=all="-N -l"
debug: $(BIN)

.PHONY: all build fmt debug
