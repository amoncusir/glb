
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
	go test -v ./...

.PHONY: run
run:
	go run .
