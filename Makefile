VERSION := v0.3.0

all: clean build-linux build-mac build-windows

build-linux:
	env GOOS=linux GOARCH=amd64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/linux/tatu-amd cmd/*.go
	env GOOS=linux GOARCH=arm64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/linux/tatu-arm cmd/*.go

build-mac:
	env GOOS=darwin GOARCH=amd64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/mac/tatu-intel cmd/*.go
	env GOOS=darwin GOARCH=arm64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/mac/tatu-apple cmd/*.go

build-windows:
	env GOOS=windows GOARCH=amd64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/windows/tatu-amd cmd/*.go
	env GOOS=windows GOARCH=arm64 go build -v -ldflags '-s -w -X main.version=$(VERSION)' -o dist/windows/tatu-arm cmd/*.go

clean:
	rm -rf dist/

.PHONY: build-linux build-mac build-windows clean
