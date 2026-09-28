.PHONY: all build run web run-web test test-race test-cover fmt vet clean install

all: test build

build:
	@mkdir -p bin
	go build -o bin/diary ./cmd/diary

run:
	go run .

web:
	go run . --web

run-web:
	go run . --web

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

install: build
	@mkdir -p $(HOME)/.local/bin $(HOME)/.local/share/applications $(HOME)/.local/share/icons/hicolor/scalable/apps $(HOME)/.local/share/icons/hicolor/64x64/apps
	install -m 755 bin/diary $(HOME)/.local/bin/diary
	cp assets/com.omarchy.diary.desktop $(HOME)/.local/share/applications/
	cp assets/icon.svg $(HOME)/.local/share/icons/hicolor/scalable/apps/diary.svg
	cp assets/icon.png $(HOME)/.local/share/icons/hicolor/64x64/apps/diary.png
	@echo "Diary successfully installed to $(HOME)/.local/bin/diary!"
