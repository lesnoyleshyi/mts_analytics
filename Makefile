GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_NAME=analytics-service
SERVICE_PATH=g6834/team17/protos


OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

PB_REL=https://github.com/protocolbuffers/protobuf/releases
PB_VER=21.2

ifeq ($(OS_NAME),Darwin)
	PROTO_OS_VERSION=osx-$(OS_ARCH)
endif

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.PHONY: run
run:
	go run cmd/auth/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'

# ----------------------------------------------------------------

.PHONY: build
build: .build

.build:
	go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-v -o ./bin/analytics-service$(shell go env GOEXE) ./cmd/analytics/main.go

.PHONY: deps
deps: protobuf deps-go

protobuf:
	curl -LO $(PB_REL)/download/v$(PB_VER)/protoc-$(PB_VER)-$(PROTO_OS_VERSION).zip
	unzip protoc-$(PB_VER)-$(PROTO_OS_VERSION).zip -d $(GO_BIN)/protoc
	rm -rf protoc-$(PB_VER)-$(PROTO_OS_VERSION).zip
	export PATH='$PATH':$(GO_BIN)/protoc/bin
	chmod u+x $(GO_BIN)/protoc/bin/protoc

.PHONY: deps-go
deps-go:
	#plugin for buffer compiler(protoc) to generate Go code(serialisation)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	#plugin for buffer compiler(protoc) to generate Go code(gRPC server & client)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
