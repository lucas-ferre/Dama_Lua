.PHONY: build build-cli build-gui run run-cli run-gui test docker-build docker-run docker-test clean

build: build-cli

build-cli:
	go build -o bin/damas ./cmd/damas

build-gui:
	go build -o bin/damas-gui ./cmd/damas-gui

build-all: build-cli build-gui

run: run-cli

run-cli:
	go run ./cmd/damas

run-gui:
	go run ./cmd/damas-gui

test:
	go test -v ./...

docker-build:
	docker build -t damas-go:latest .

docker-run:
	docker run --rm -it damas-go:latest

docker-test:
	docker run --rm damas-go:latest go test -v ./...

clean:
	rm -rf bin/
