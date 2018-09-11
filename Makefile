all: test

deps:
	@which godep > /dev/null || go get github.com/tools/godep
	@export CGO_CFLAGS_ALLOW=-L.*; godep go install

test: deps
	@export CGO_CFLAGS_ALLOW=-L.*; go list ./... | grep -v /vendor/ | xargs -n1 godep go test

bench: deps
	@export CGO_CFLAGS_ALLOW=-L.*; go list ./... | grep -v /vendor/ | xargs -n1 godep go test -run=XXX -benchtime=1s -bench=.

.PNONY: all deps test
