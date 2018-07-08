default: test

.PHONY: dev test build run push

dev:
	@docker-compose up

test:
	@CONFIG=config/test.config.json GOCACHE=off go test ./...

build:
	@docker build -t followboard-api .

run:
	@go run main.go -logtostderr=true

push:
	@now
