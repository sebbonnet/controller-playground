project_dir := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))

VERSION ?= $(shell git rev-parse --short HEAD)
NAMESPACE ?= my-namespace
IMG ?= controller-playground:latest

.PHONY: build
build:
	go build -o bin/manager main.go

.PHONY: run
run:
	go run ./main.go

.PHONY: docker-build
docker-build:
	docker build -t $(IMG):$(VERSION) .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push $(IMG):$(VERSION)

.PHONY: deploy
deploy:
	VERSION=$(VERSION) IMG=$(IMG) NAMESPACE=$(NAMESPACE) envsubst < $(project_dir)/deploy.yml | kubectl -n $(NAMESPACE) apply -f -
