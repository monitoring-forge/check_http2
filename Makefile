VERSION=0.0.23
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} "

all: check_http2

.PHONY: check_http2

check_http2: writer.go checker.go main.go
	go build $(LDFLAGS) -o check_http2 writer.go checker.go main.go

linux: writer.go checker.go main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check_http2 writer.go checker.go main.go

check:
	go test -v ./...

fmt:
	go fmt ./...
