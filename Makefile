
VERSION := "v0.0.0"
BUILD := $(shell git rev-parse HEAD)
BINARY_NAME="service"
LDFLAGS=-ldflags "-X=$(BINARY_NAME).Version=$(VERSION) -X=$(BINARY_NAME).Build=$(BUILD)"
BASE_NAME := go-share-secret
DOCKER_IMAGE_NAME := $(BASE_NAME):$(BUILD)
NETWORK_NAME := $(BASE_NAME)-$(BUILD)
DOCKER_LINT_IMAGE_NAME := golangci-lint:$(BUILD)

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
	docker-compose up -d

install:
	go install

build:
	 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BINARY_NAME) $(LDFLAGS) && chmod 777 $(BINARY_NAME)

docker-network:
	docker network create $(NETWORK_NAME) || true	

docker-image-build: docker-network
	docker build --build-arg NETRC -t $(DOCKER_IMAGE_NAME) --no-cache .

docker-network-clean:
	docker network rm $(NETWORK_NAME) || true

docker-clean:
	docker rmi --force $(DOCKER_IMAGE_NAME) || true

docker-clean-volumes:
	docker volume prune -f || true

docker-test:
	docker run --rm $(DOCKER_IMAGE_NAME) make test	

docker-build-lint:
	docker build -f lint.Dockerfile --build-arg NETRC -t $(DOCKER_LINT_IMAGE_NAME) .
	
docker-lint:docker-build-lint
	docker run --rm -v $(PWD):/app -w /app $(DOCKER_LINT_IMAGE_NAME) golangci-lint run --timeout 5m	