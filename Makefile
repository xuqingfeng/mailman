.PHONY: all run

all: run

run: build
	./mailman

build: build-template
	go-bindata ui/... && go build

build-template: fmt # TODO: change package name
	cd mail && go-bindata -o bindata.go && cd -

fmt:
	go fmt ./...

test:
	go test ./...

bin:
	gox -osarch="linux/amd64 linux/386 linux/arm darwin/amd64 windows/386 windows/amd64"
