GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions
STATIK_FILES := $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))

MAIN_GO_FILES := $(wildcard cmd/*)
BINARYS := $(patsubst cmd/%.go, %, $(MAIN_GO_FILES)) 


all: build

# Allows running things such as make <binary_name>
$(BINARYS): build

.PHONY: build
build: $(STATIK_FILES)
	$(GO) build $(MAIN_GO_FILES)

.PHONY: clean
clean: 
	rm -f $(STATIK_FILES)
	rm -f $(BINARYS)

%_statik/statik.go: $(wildcard %/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*_statik) -f