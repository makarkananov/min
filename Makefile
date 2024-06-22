default: all

all: build run

build:
	docker compose build

run:
	docker compose up -d

lint:
	golangci-lint run -v