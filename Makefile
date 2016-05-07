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
