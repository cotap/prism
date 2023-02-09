all: test

deps:
	@export CGO_CFLAGS_ALLOW=-L.*; go install

test: deps
	@export CGO_CFLAGS_ALLOW=-L.*; go list ./... | grep -v /vendor/ | xargs -n1 go test

bench: deps
	@export CGO_CFLAGS_ALLOW=-L.*; go list ./... | grep -v /vendor/ | xargs -n1 go test -run=XXX -benchtime=1s -bench=.

.PNONY: all deps test
