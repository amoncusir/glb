
PWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
SRC = $(PWD)

#GO ?= docker run --rm -it -v $(PWD):/go \
#	--net=host \
#	golang:1.24.4 go
GO ?= go
GO_PATH ?= $(shell go env GOPATH)/bin/
MOCKGEN ?= $(GO_PATH)mockgen


define DEV_TOOLS
	go.uber.org/mock/mockgen@latest
endef

# Add .env environment variables unless is running inside of CI pipeline
ifeq ($(CI), false)
	-include $(PWD)/.env
endif

export

-include *.mk


## Target convention naming:
## <Action>[-<Identifier>] :: Examples:
### `install` -> Just the action because is a generic task may implies other tasks.
### `install-poetry` -> The action first, the name after.
### `build-docker` -> Action and identifier.
### `rm-build-docker` -> Action taken for a action result.
## Why?
# Because helps to find the correct targets using the Shell AutoCompletion.

install: | install-mod install-dev

install-mod:
	$(GO) get $(PWD)

install-dev:
	$(foreach dep,$(DEV_TOOLS), $(GO) install $(dep);)

update:
	$(GO) mod tidy

lint:
	$(GO) fmt

.PHONY: test
test: mock
	# -count=1 flag disable the test cache and forces to run at least one every test again
	$(GO) test -v -count=1 $(PWD)...

.PHONY: run
run:
	$(GO) run .

server:
	docker run --rm -t -e HTTP_PORT=9000 -p 9000:9000 mendhak/http-https-echo

client: URI ?= /hi
client:
	curl -v http://localhost:9090$(URI)
