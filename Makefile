.PHONY: build run run-build test clean

build:
	go build -o ./bin/api

run-build: build
	./bin/api

test:
	go test ./... -v

destroy:
	docker compose down -v

up-build:
	docker-compose up -d --build

up:
	docker-compose up

stop:
	docker-compose stop api

rebuild:
	docker-compose down api
	docker-compose up api -d --build

ingest_files:
	./scripts/ingest_files.sh 16

