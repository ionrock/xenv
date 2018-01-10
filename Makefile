test:
	dep ensure
	go test ./...

build:
	go build ./cmd/xenv
