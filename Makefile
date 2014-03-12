all:
	@mkdir -p bin/
	@bash --norc -i ./scripts/build.sh

test: deps
	go list ./... | xargs -n1 go test

.PNONY: all test
