
PWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
SRC = $(PWD)
MOCK = $(PWD)test/mock/

## Target convention naming:
## <Action>[-<Identifier>] :: Examples:
### `install` -> Just the action because is a generic task may implies other tasks.
### `install-poetry` -> The action first, the name after.
### `build-docker` -> Action and identifier.
### `rm-build-docker` -> Action taken for a action result.
## Why?
# Because helps to find the correct targets using the Shell AutoCompletion.

.PHONY: test
test:
	go test -v $(PWD)...

.PHONY: run
run:
	go run .

.PHONY: mock
mock: $(MOCK)instance.go $(MOCK)request_connection.go $(MOCK)net_address.go
	@echo "exec mocks"

$(MOCK)instance.go: $(SRC)pkg/service/instance/instance.go
	mockgen -source $^ -destination $@ -package mock

$(MOCK)request_connection.go: $(SRC)pkg/types/request.go
	mockgen -source $^ -destination $@ -package mock

$(MOCK)net_address.go:
	mockgen -destination $@ -package mock net Addr