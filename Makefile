VERSION=`git describe --abbrev=0 --tags`

.PHONY: all run

all: run

run: build
	./mailman

build: generate
	go build -o mailman main.go

generate:
	go generate

fmt:
	go fmt ./...

test:
	go vet ./... && go test -v ./...

build-all: fmt generate
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-linux-amd64 main.go && \
	GOOS=linux GOARCH=arm go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-linux-arm main.go && \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-darwin-amd64 main.go && \
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-windows-amd64.exe main.go
