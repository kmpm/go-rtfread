GOTEST=go test -v -timeout=5s -cover -coverprofile=coverage.out -covermode=atomic

test:
	$(GOTEST) ./...
	