GO := go

# Do not use trailing slashes here
STATIK_SRCS :=langserver/documentation/functions

.PHONY: build
build: statik
	$(GO) build cmd/*

.PHONY: clean
clean: 
	rm -f $(patsubst %, %-statik/statik.go, $(STATIK_SRCS))
	rm -f $(patsubst cmd/%.go, %, $(wildcard cmd/*))

.PHONY: statik
statik: $(patsubst %, %-statik/statik.go, $(STATIK_SRCS))

%-statik/statik.go: $(wildcard %/*)
	statik -src "$*" -dest $(dir $*) -p $(notdir $*-statik) -f