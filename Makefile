GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

help:
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: api runners ## build all images and run in docker
	docker run -it --rm -p 8090:8090 -v "/var/run/docker.sock:/var/run/docker.sock" yximages/yxi-api

.PHONY: fmt
fmt: ## formate all go files (use go)
	$(GOFMT) -w $(GOFILES)

.PHONY: dev
dev:fmt vet ## formate vet and compile (use go)
	go build -mod=vendor -ldflags '-w -s' -o main cmd/apiServer/main.go

vet: ## vat all go files (use go)
	go vet $(PACKAGES)

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

pullali: ## pull runner images from aliyun
	cd scripts && ./images.sh -l

pushali: ## push runner images to aliyun
	cd scripts && ./images.sh -a

clean: ## clean local build
	rm ./main

