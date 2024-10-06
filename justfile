default:
    @just --list

test:
    go test ./...

test-race:
    go test -race ./...

build-windows:
    GOOS=windows GOARCH=amd64 go build -o out/tks.exe ./tks

run:
    go run ./tks

clean:
    rm -rf out/*

fmt:
    go fmt ./...

vet:
    go vet ./...