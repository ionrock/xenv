dep:
	dep ensure

test: dep
	go test ./...

build: dep
	go build ./cmd/xenv


example: build
	./xenv --config examples/config.yml echo 'Hi'
