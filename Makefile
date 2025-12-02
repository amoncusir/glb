
PWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
SRC = $(PWD)

#GO ?= docker run --rm -it -v $(PWD):/go \
#	--net=host \
#	golang:1.24.4 go
GO ?= go


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

install:
	$(GO) get $(PWD)

update:
	$(GO) mod tidy

.PHONY: test
test:
	$(GO) test -v $(PWD)...

.PHONY: run
run:
	$(GO) run .
