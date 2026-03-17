.PHONY: build run clean install

BINARY = pulsereader
VERSION = 0.2.0

build:
	go build -ldflags="-s -w" -o bin/$(BINARY) ./cmd/pulsereader

run:
	go run ./cmd/pulsereader $(ARGS)

install:
	go install ./cmd/pulsereader

clean:
	rm -rf bin/

# Cross-compilation
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY)-linux-amd64 ./cmd/pulsereader

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY)-windows-amd64.exe ./cmd/pulsereader

build-mac:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/$(BINARY)-darwin-arm64 ./cmd/pulsereader

build-all: build-linux build-windows build-mac
