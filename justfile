default:
    @just --list

test:
    go test ./...

test-race:
    go test -race ./...

build-windows:
    GOOS=windows GOARCH=amd64 go build -o out/tksgo.exe ./tks

build-mac:
    GOOS=darwin GOARCH=amd64 go build -o out/tksgo ./tks

build-all: build-windows build-mac

run:
    go run ./tks

clean:
    rm -rf out/*

fmt:
    go fmt ./...

vet:
    go vet ./...