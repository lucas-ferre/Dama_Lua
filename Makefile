.PHONY: build run test docker-build docker-run docker-test clean

build:
	go build -o bin/damas ./cmd/damas

run:
	go run ./cmd/damas

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
