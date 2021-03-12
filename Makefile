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
	go run ./server/

test:
	go test -v -count 1 ./server

testcover:
	go test -cover $(shell go list ./server)

up:
	docker-compose up -d

open:
	bash +x ./scripts/openwebpage.sh