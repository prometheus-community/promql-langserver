GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions
MAIN_GO_FILES := $(wildcard cmd/*)
BINARYS := $(patsubst cmd/%.go, %, $(MAIN_GO_FILES)) 

all: build

# Allows running things such as make <binary_name>
$(BINARYS): build

.PHONY: build
build: statik
	$(GO) build $(MAIN_GO_FILES)

.PHONY: clean
clean: 
	rm -f $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))
	rm -f $(BINARYS)

.PHONY: statik
statik: $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))

%_statik/statik.go: $(wildcard %/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*_statik) -f