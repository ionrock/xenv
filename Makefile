SOURCES=$(shell find . -name '*.go')
PROTO=$(shell find . -name '*.proto')

dep:
	dep ensure

test: dep
	go test ./...

build: dep $(SOURCES)
	go build ./cmd/xenv

example: build
	./xenv --debug --config examples/config.yml examples/web-server

svcs: $(SOURCES) proto
	go build ./cmd/svcs

proto: $(PROTO)
	protoc -I manager/ manager/manager_api.proto --go_out=plugins=grpc:manager
