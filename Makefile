BIN := "smgo"
DOCKER_IMG="smgo:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build run version build-img run-img test
build: build-lin-amd64

run: build-lin-amd64
	./bin/$(BIN)

version: build-lin-amd64
	$(BIN) version

build-img:
	docker build \
		--network=host \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -race -count 100 ./internal/... ./pkg/...

.PHONY: lint install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: generate
generate:
	rm -f ./pkg/smgo/api/*
	protoc ./api/v1/SmgoService.proto \
		-I ./api/ \
		--go_out=./pkg/ --go-grpc_out=./pkg/ \
		--grpc-gateway_out ./pkg/

.PHONY: build-all build-lin-amd64 build-lin-386 build-mac-arm64 build-win-amd64
build-all: build-lin-amd64 build-lin-386 build-mac-arm64 build-win-amd64

build-lin-amd64:
	rm -f ./bin/$(BIN)
	rm -f ./bin/$(BIN)-client
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BIN) -ldflags "$(LDFLAGS)" ./cmd/smgo
	GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BIN)-client -ldflags "$(LDFLAGS)" ./cmd/smgo-client

build-lin-386:
	rm -f ./bin/lin386/$(BIN)
	rm -f ./bin/lin386/$(BIN)-client
	GOOS=linux GOARCH=386 go build -v -o ./bin/lin386/$(BIN) -ldflags "$(LDFLAGS)" ./cmd/smgo
	GOOS=linux GOARCH=386 go build -v -o ./bin/lin386/$(BIN)-client -ldflags "$(LDFLAGS)" ./cmd/smgo-client

build-mac-arm64:
	rm -f ./bin/mac/$(BIN)
	rm -f ./bin/mac/$(BIN)-client
	GOOS=darwin GOARCH=arm64 go build -v -o ./bin/mac/$(BIN) -ldflags "$(LDFLAGS)" ./cmd/smgo
	GOOS=darwin GOARCH=arm64 go build -v -o ./bin/mac/$(BIN)-client -ldflags "$(LDFLAGS)" ./cmd/smgo-client

build-win-amd64:
	rm -f ./bin/win/$(BIN).exe
	rm -f ./bin/win/$(BIN)-client.exe
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/win/$(BIN).exe -ldflags "$(LDFLAGS)" ./cmd/smgo
	GOOS=windows GOARCH=amd64 go build -v -o ./bin/win/$(BIN)-client.exe -ldflags "$(LDFLAGS)" ./cmd/smgo-client

.PHONY: up stop
up:
	docker-compose up -d --build

stop:
	docker-compose stop
