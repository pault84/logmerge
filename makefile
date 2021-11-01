export GO111MODULE=on
export GOFLAGS = -mod=vendor

ifndef PKGS
	PKGS := $(shell GOFLAGS=-mod=vendor go list ./... 2>&1 | grep -v 'go: ' | grep -v 'github.com/pault84/logmerge/vendor' | grep -v versioned)
endif

GO_FILES := $(shell find . -name '*.go' | grep -v 'vendor' | \
                                   grep -v '\.pb\.go' | \
                                   grep -v '\.pb\.gw\.go' | \
                                   grep -v 'externalversions' | \
                                   grep -v 'versioned' | \
                                   grep -v 'generated')

.DEFAULT_GOAL: all
.PHONY: build

all: pretest build

pretest: check-fmt lint vet errcheck staticcheck

lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	for file in $(GO_FILES); do \
        golint $${file}; \
        if [ -n "$$(golint $${file})" ]; then \
            exit 1; \
        fi; \
        done

vet:
	go vet $(PKGS)
staticcheck:
	GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck $(PKGS)

errcheck:
	GO111MODULE=off go get -u github.com/kisielk/errcheck
	errcheck -ignoregenerated -ignorepkg fmt -verbose -blank $(PKGS)
check-fmt:
	bash -c "diff -u <(echo -n) <(gofmt -l -d -s -e $(GO_FILES))"

do-fmt:
	gofmt -s -w $(GO_FILES)

vendor-sync:
	go mod tidy
	go mod vendor

build:
	@echo "Build logmerge"
	@go build -o bin/logmerge .