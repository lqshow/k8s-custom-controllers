SHELL := /bin/bash

BUILD_IMAGE := golang:alpine
REGISTRY := docker.io/lqshow
IMAGE_NAME := $(REGISTRY)/foobar-code-generator
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)
MODULE := github.com/lqshow/k8s-custom-controllers/foobar-code-generator
VERSION:= $(shell git rev-parse --short HEAD)

all: build

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: clean-bin
clean-bin:
	rm -rf bin .go

.PHONY: build
build:
	@$(MAKE) $(SOURCES)
	$(shell go env -w GOPRIVATE=github.com/lqshow/*)
	GOBIN=$(shell mkdir -p ./bin && cd ./bin && pwd) go install -ldflags "-X $(MODULE)/pkg/version.VERSION=$(VERSION)" ./...

.PHONY: clean
clean:
	$(MAKE) clean-bin

.PHONY: run
run:
	@bin/foobar

.PHONY: build-image
build-image:
	@docker build -t $(IMAGE_NAME):$(VERSION) -f Dockerfile --build-arg BUILD_IMAGE=$(BUILD_IMAGE) .