.PHONY: all run

all: run

run: build
	./mailman

build: fmt
	go-bindata ui/... && go build

fmt:
	go fmt ./...

test:
	go test ./...

bin:
	gox -osarch="linux/amd64 linux/386 linux/arm darwin/amd64 windows/386 windows/amd64"
