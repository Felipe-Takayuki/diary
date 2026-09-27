.PHONY: all build run test test-race test-cover fmt vet clean

all: test build

build:
	@mkdir -p bin
	go build -o bin/diary ./cmd/diary

run:
	go run .

test:
	go test -v ./...

test-race:
	go test -race ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin coverage.out coverage.html
