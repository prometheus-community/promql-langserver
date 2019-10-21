GO := go 

.PHONY: build
build: 
	$(GO) build cmd/*
.PHONY: clean
clean: 
	git clean -f .

