.PHONY: all run

all: run

run: build
	./mailman

build: format
	go build

format:
	go fmt ./...

test:
	go test ./...

# user: Current not implemented on linux/arm
gox:
	gox -osarch="linux/amd64 linux/386 linux/arm64 linux/arm darwin/amd64 darwin/386"
