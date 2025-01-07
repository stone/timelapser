# Variables
BINARY_NAME=timelapser
DOCKER_REPO=stone
IMAGE_NAME=$(DOCKER_REPO)/$(BINARY_NAME)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
UUID := $(shell uuidgen)
DEVCONTAINER=ttl.sh/$(BINARY_NAME)-$(UUID):1h
VERSION?=1.0.0

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean docker-build docker-push test

all: clean build

# Build the application locally
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Build docker image
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--label org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		-t $(IMAGE_NAME):$(VERSION) \
		-t $(IMAGE_NAME):latest \
		.

# Push devcontainer to ttl.sh
docker-push-devcontainer:
	docker tag $(IMAGE_NAME):$(VERSION) $(DEVCONTAINER)
	docker push $(DEVCONTAINER)

# Push docker image
docker-push:
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest

# Development target for running locally
dev:
	go run $(LDFLAGS) .
