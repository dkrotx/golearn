.PHONY: mocks
.PHONY: test

mocks:
	@mkdir -p mocks
	mockgen -destination=mocks/mock_keyvalue.go  golearn/mocking/keyvalue KeyValue

test:
	go test -v
