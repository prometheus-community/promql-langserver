GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions
STATIK_FILES := $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))

MAIN_GO_FILES := $(wildcard cmd/*)
BINARYS := $(patsubst cmd/%.go, %, $(MAIN_GO_FILES)) 


all: build

generated: $(STATIK_FILES)

# Allows running things such as make <binary_name>
$(BINARYS): build

.PHONY: build
build: $(STATIK_FILES)
	$(GO) build $(MAIN_GO_FILES)

.PHONY: clean
clean: 
	rm -f $(STATIK_FILES)
	rm -f $(BINARYS)

%_statik/statik.go: $(wildcard $*/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*_statik) -f -m

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: update_internal_packages
update_internal_packages:
	for dir in `ls -d vendored/*/`;                                           \
	do                                                                        \
	    echo "Updating:";                                                     \
	    NAME=`basename $$dir`;                                                \
	    echo "Name: " $$NAME;                                                 \
	    REPO=`cat vendored/$$NAME.repo`;                                      \
	    echo "Repo: " $$REPO;                                                 \
	    FCMD=`cat vendored/$$NAME.cmd`;                                       \
	    echo "File cmd: " $$FCMD;                                             \
	    DIRS=`cat vendored/$$NAME.dirs`;                                      \
	    echo "Directories: " $$DIRS;                                          \
	    VERSION=`cat vendored/$$NAME.version`;                                \
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
	        cp -r $$TMPDIR/internal/$$subdir $$dir$$subdir;                   \
	    done;                                                                 \
	    for file in `find $$dir -type f`;                                     \
	    do                                                                    \
	        CMD=`echo $$FCMD $$file`;                                         \
	        echo $$CMD;                                                       \
	        bash -c "$$CMD";                                                  \
	    done;                                                                 \
	    rm -rf $$TMPDIR;                                                      \
	    make fmt;                                                             \
	done