export CGO_ENABLED ?= 0
export GOFLAGS ?= -mod=readonly -trimpath
export GO111MODULE = on
CMD ?= $(filter-out OWNERS, $(notdir $(wildcard ./cmd/*)))
GOBUILDFLAGS ?= -v
GIT_HEAD ?= $(shell git log -1 --format=%H)
GIT_VERSION = $(shell git describe --tags --always)
LDFLAGS += -extldflags '-static'
BUILD_DEST ?= _build
GOTOOLFLAGS ?= $(GOBUILDFLAGS) -ldflags '-w $(LDFLAGS)' $(GOTOOLFLAGS_EXTRA)

.PHONY: all
all: build

.PHONY: build
build: $(CMD)

.PHONY: $(CMD)
$(CMD): %: $(BUILD_DEST)/%

$(BUILD_DEST)/%: cmd/%
	go build $(GOTOOLFLAGS) -o $@ ./cmd/$*

.PHONY: clean
clean:
	rm -rf $(BUILD_DEST)
	@echo "Cleaned $(BUILD_DEST)"

.PHONY: run
run:
	$(BUILD_DEST)/testapp1 \
	  --kubeconfig kind.kubeconfig \
	  --kcp-kubeconfig kcp.kubeconfig \
	  -v 6
