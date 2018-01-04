build:
	go build ./cmd/xenv

test:
	go test ./... -ignore ./vendor
