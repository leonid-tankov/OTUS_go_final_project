BIN := "./bin/app"
DOCKER_IMG="banner-rotation:develop"
DOCKER_COMPOSE_DIR="./deployments/docker-compose.yaml"

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		banner_rotation.proto

go-build:
	go build -v -o $(BIN) ./cmd

build:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build
	docker run $(DOCKER_IMG)

run: build
	docker-compose -f $(DOCKER_COMPOSE_DIR) up

test:
	go test  -race -count 100 ./internal/...

integration-tests: run
	go test -tags integration -race ./internal/...
	docker-compose -f $(DOCKER_COMPOSE_DIR) down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

down:
	docker-compose -f $(DOCKER_COMPOSE_DIR) down

.PHONY: generate go-build build run-img run test integration-tests lint down
