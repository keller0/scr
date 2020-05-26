GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

help:
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: api runners ## build all images and run in docker
	docker run -it --rm -p 8090:8090 -v "/var/run/docker.sock:/var/run/docker.sock" yximages/yxi-api

.PHONY: fmt
fmt: ## format all go files (use go)
	$(GOFMT) -w $(GOFILES)

vet: ## vat all go files (use go)
	go vet $(PACKAGES)

.PHONY: dev
dev:fmt vet ## format vet and compile (use go)
	go build -mod=vendor -ldflags '-w -s' -o main cmd/apiServer/main.go

test: ## run test (use go)
	go test -v -mod=vendor ./...

api: ## build api image
	docker build -t yximages/yxi-api .

drun: api ## build api image and run it
	docker run -it --rm -p 8090:8090 -v "/var/run/docker.sock:/var/run/docker.sock" yximages/yxi-api

runners: ## build runner images
	cd scripts && ./images.sh -b

push: ## push runner images to docker hub
	cd scripts && ./images.sh -d

pull: ## pull runner images from docker hub
	cd scripts && ./images.sh -p

clean: ## clean local build
	rm ./main

