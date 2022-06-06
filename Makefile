.PHONY: build start unit-test start-dep stop-dep

BIN=./bin/app

build:
	go build -o $(BIN) ./transactions

copy-env:
	cp .env.example .env

unit-test:
	go test ./...

start-dependencies:
	docker-compose up -d

stop-dependencies:
	docker-compose stop

webhooks-events-consumer:
	$(BIN) acquirer-webhooks-consumer

migrations:
	$(BIN) migrations

seeds:
	$(BIN) seeds