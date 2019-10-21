GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions

.PHONY: build
build: statik
	$(GO) build cmd/*

.PHONY: clean
clean: 
	rm -f $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))
	rm -f $(patsubst cmd/%.go, %, $(wildcard cmd/*))

.PHONY: statik
statik: $(patsubst %, %_statik/statik.go, $(STATIK_SRCS))

%_statik/statik.go: $(wildcard %/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*_statik) -f