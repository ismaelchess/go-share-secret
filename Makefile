
VERSION := "v0.0.0"
BASE_NAME := go-share-secret
BINARY_NAME="service"

BUILD := $(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X=$(BINARY_NAME).Version=$(VERSION) -X=$(BINARY_NAME).Build=$(BUILD)"
DOCKER_IMAGE_NAME := $(BASE_NAME):$(BUILD)
NETWORK_NAME := $(BASE_NAME)-$(BUILD)
DOCKER_LINT_IMAGE_NAME := golangci-lint:$(BUILD)

FORCE:

clean: down

build:
	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) $(LDFLAGS) && chmod 777 $(BINARY_NAME)

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
	docker-compose up -d

install:
	go install


# Jekins using Dockers
jenkins: jk-docker-image-build jk-docker-test jk-docker-lint jk-docker-clean-all

jk-docker-build-lint:
	docker build -f lint.Dockerfile -t $(DOCKER_LINT_IMAGE_NAME) .

jk-docker-clean:
	docker rmi --force $(DOCKER_IMAGE_NAME) || true

jk-docker-clean-volumes:
	docker volume prune -f || true

jk-docker-lint:jk-docker-build-lint
	docker run --rm -v $(PWD):/app -w /app $(DOCKER_LINT_IMAGE_NAME) golangci-lint run --timeout 5m	

jk-docker-image-build: jk-docker-network
	docker build -t $(DOCKER_IMAGE_NAME) --no-cache .

jk-docker-network:
	docker network create $(NETWORK_NAME) || true	

jk-docker-network-clean:
	docker network rm $(NETWORK_NAME) || true

jk-docker-test:
	docker run --rm $(DOCKER_IMAGE_NAME) make test	

jk-docker-clean-all:jk-docker-clean jk-docker-clean-volumes jk-docker-network-clean


