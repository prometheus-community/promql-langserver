GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions
STATIK_FILES := $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))

BINARYS := $(patsubst cmd/%.go, %, $(MAIN_GO_FILES)) 

GOLANGCI_LINT_VERSION ?= v2.6.0

all: build test 

generated: $(STATIK_FILES)


.PHONY: install
install: $(STATIK_FILES)
	$(GO) get ./cmd/...

.PHONY: build
build:
	$(GO) build ./cmd/...

.PHONY: clean
clean: 
	rm -f $(STATIK_FILES)
	rm -f $(BINARYS)

%_statik/statik.go: $(wildcard $*/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*_statik) -f -m
	gofmt -w $@

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -race -v -cover ./...

.PHONY: update_internal_packages
update_internal_packages:
	for dir in `ls -d internal/vendored/*/`;                                  \
	do                                                                        \
	    echo "Updating:";                                                     \
	    NAME=`basename $$dir`;                                                \
	    echo "Name: " $$NAME;                                                 \
	    REPO=`cat internal/vendored/$$NAME.repo`;                             \
	    echo "Repo: " $$REPO;                                                 \
	    FCMD=`cat internal/vendored/$$NAME.cmd`;                              \
	    echo "File cmd: " $$FCMD;                                             \
	    DIRS=`cat internal/vendored/$$NAME.dirs`;                             \
	    echo "Directories: " $$DIRS;                                          \
	    VERSION=`cat internal/vendored/$$NAME.version`;                       \
	    echo "Version: " $$VERSION;                                           \
	    echo "Cleaning up";                                                   \
	    rm -rf $$dir*;                                                        \
	    TMPDIR=`mktemp -d`;                                                   \
	    echo "Temp dir: " $$TMPDIR;                                           \
	    git clone $$REPO $$TMPDIR;                                            \
	    git --git-dir=$$TMPDIR/.git --work-tree=$$TMPDIR checkout $$VERSION;  \
	    echo "Copying Files";                                                 \
	    for subdir in $$DIRS;                                                 \
	    do                                                                    \
	        echo mkdir -p `dirname $$dir$$subdir`;                            \
	        mkdir -p `dirname $$dir$$subdir`;                                 \
	        cp -r $$TMPDIR/internal/$$subdir $$dir$$subdir;                   \
	    done;                                                                 \
	    cp $$TMPDIR/LICENSE $$dir;                                            \
	    for file in `find $$dir -type f`;                                     \
	    do                                                                    \
	        CMD=`echo $$FCMD $$file`;                                         \
	        echo $$CMD;                                                       \
	        bash -c "$$CMD";                                                  \
	    done;                                                                 \
	    rm -rf $$TMPDIR;                                                      \
	    make fmt;                                                             \
	done

.PHONY: htmlcover
htmlcover:
	go test -coverprofile=coverage.out ./langserver/...
	go tool cover -html=coverage.out

.PHONY: crossbuild
crossbuild:
	goreleaser build --snapshot

.PHONY: release
release:
	goreleaser release

.PHONY: golangci-lint-version
golangci-lint-version:
	@echo $(GOLANGCI_LINT_VERSION)
