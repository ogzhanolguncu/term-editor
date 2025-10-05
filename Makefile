.PHONY: build run clean dev

build:
	go build -o bin/app .

run: build
	./bin/app

dev:
	go run .

clean:
	rm -rf bin/

test:
	go test ./...
