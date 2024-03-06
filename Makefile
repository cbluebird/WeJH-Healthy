build:
	go build -o healthy

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GIN_MODE=release go build -o healthy

.PHONY: build build-linux