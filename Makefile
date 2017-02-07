all: test

deps:
	@which godep > /dev/null || go get github.com/tools/godep
	@godep go install

test: deps
	@go list ./... | grep -v /vendor/ | xargs -n1 godep go test

bench: deps
	@go list ./... | grep -v /vendor/ | xargs -n1 godep go test -run=XXX -benchtime=1s -bench=.

.PNONY: all deps test
