MOCK = $(PWD)test/mock/


.PHONY: mock
mock: $(MOCK)instance.go $(MOCK)request_connection.go $(MOCK)net_address.go
	@echo "exec mocks"

$(MOCK)instance.go: $(SRC)pkg/service/instance/instance.go
	mockgen -source $^ -destination $@ -package mock

$(MOCK)request_connection.go: $(SRC)pkg/types/request.go
	mockgen -source $^ -destination $@ -package mock

$(MOCK)net_address.go:
	mockgen -destination $@ -package mock net Addr