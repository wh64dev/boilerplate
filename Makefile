GO = go
SRCS = *.go **/*.go

MAIN = .
FRONTEND_DIR = frontend/
OUTPUT = sample-app

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | cut -d' ' -f3)

LD_FLAGS := -X config.version=$(FULL_VERSION)	\
            -X config.branch=$(BRANCH)		\
            -X config.hash=$(HASH)		\
            -X config.buildTime=$(BUILD_TIME)	\
            -X config.goVersion=$(GO_VERSION)

BUILD_FLAGS := -ldflags "$(LDFLAGS)"

.PHONY: all build-frontend build clean

all: build-frontend build

build-frontend: $(FRONTEND_DIR)
	@rm -rf public/
	@make -C $(FRONTEND_DIR)

build: $(SRCS)
	@go build $(BUILD_FLAGS) -o $(OUTPUT) $(MAIN)

clean:
	@make -C $(FRONTEND_DIR) clean
	@rm -rf $(OUTPUT)
	@rm -rf public/

