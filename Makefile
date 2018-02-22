SOURCES=$(shell find . -name '*.go')

dep:
	dep ensure

test: dep
	go test ./...

build: dep $(SOURCES)
	go build ./cmd/xenv

example: build
	./xenv --debug --config examples/config.yml examples/web-server

svcs: $(SOURCES)
	go build ./cmd/svcs
