MOCK = $(PWD)test/mock/


.PHONY: mock
mock: $(MOCK)instance.go $(MOCK)request_connection.go $(MOCK)net_address.go

$(MOCK)instance.go: $(SRC)pkg/service/instance/instance.go
	$(MOCKGEN) -source $^ -destination $@ -package mock

$(MOCK)request_connection.go: $(SRC)pkg/types/request.go
	$(MOCKGEN) -source $^ -destination $@ -package mock

$(MOCK)net_address.go:
	$(MOCKGEN) -destination $@ -package mock net Addr
