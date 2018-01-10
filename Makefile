dep:
	dep ensure

test: dep
	go test ./...

build: dep
	go build ./cmd/xenv
