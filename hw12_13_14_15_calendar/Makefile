.EXPORT_ALL_VARIABLES:
BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
BIN_CLIENT := "./bin/client"
DOCKER_IMG="calendar:develop"
CONN_STR = "postgresql://localhost/postgres?user=postgres&password=sqlSync24&sslmode=disable"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

RABBIT_URL := amqp://guest:guest@localhost:5672/

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.toml

client:
	go run ./cmd/client

scheduler:
	go run ./cmd/scheduler -config ./configs/scheduler_config.toml

sender:
	go run ./cmd/sender -config ./configs/sender_config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile_calendar .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CALENDAR) version
	$(BIN_SCHEDULER) version
	$(BIN_SENDER) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/EventService.proto

timeout5s:
	sleep 10

goose_up:
	goose -dir ./migrations postgres $(CONN_STR) up
goose_down:
	goose -dir ./migrations postgres $(CONN_STR) down

build_client:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t "integration-tests" \
		-f Dockerfile_client .

docker_build: build
	docker-compose -f docker-compose.yaml build
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t "integration-tests" \
		-f Dockerfile_client .

docker_up:
	docker-compose -f docker-compose.yaml up -d

docker_down:
	docker-compose -f docker-compose.yaml down -v --remove-orphans

int_test: build_client
	docker run --network host integration-tests && echo PASS || echo FAIL

up: docker_up
	
down: docker_down

integration-tests: docker_up timeout5s int_test docker_down
