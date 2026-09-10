VERSION=`git describe --abbrev=0 --tags`

.PHONY: all build generate copy-assets fmt test run build-all

all: run

run: build
	./mailman

build: generate
	go build -o mailman main.go

generate: copy-assets
	go generate

copy-assets:
	npm run copy-assets

fmt:
	go fmt ./...

test:
	go vet ./... && go test -v ./...

build-all: fmt generate
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-linux-amd64 main.go && \
	GOOS=linux GOARCH=arm go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-linux-arm main.go && \
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-darwin-amd64 main.go && \
	GOOS=darwin GOARCH=arm64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-darwin-arm64 main.go && \
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o out/mailman-windows-amd64.exe main.go
