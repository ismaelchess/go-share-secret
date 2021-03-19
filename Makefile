
WORKING_DIR := $(shell pwd)
REDIS_DIR := ${HOME}/data/redis

FORCE:

clean: down

down:
	docker-compose down --remove-orphans

dev:get up log

get:
	go get -t ./...

lint:
	golangci-lint run --timeout 5m

log:
	docker-compose logs -f

run:
	go run ./

test:
	go test -v -count 1 .

testcover:
	go test -cover $(shell go list .)

up:
	REDIS_DIR=$(REDIS_DIR)  WORKING_DIR=$(WORKING_DIR) docker-compose up -d