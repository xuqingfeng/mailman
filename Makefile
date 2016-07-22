.PHONY: all run

all: run

run: build
	./mailman

build: fmt
	go build

fmt:
	go fmt ./...

test:
	go test ./...

gox:
	gox -osarch="linux/amd64 linux/386 linux/arm darwin/amd64"
