all: test

deps:
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

test: deps
	@go list ./... | xargs -n1 go test

.PNONY: all test
