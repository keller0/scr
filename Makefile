GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

help:
	@echo "vet/test/fmt"
	@echo "push/push pullali/pushali"
	@echo "all: build all images and run in docker"
	@echo "api: build yxi-api image"
	@echo "drun: run yxi-api in docker"
	@echo "runners: build all runner images"

all: api runners
	docker run -it --rm -p 8090:8090 -v "/var/run/docker.sock:/var/run/docker.sock" yximages/yxi-api

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: dev
dev:fmt vet
	go build -mod=vendor -ldflags '-w -s' -o main cmd/apiServer/main.go

vet:
	go vet $(PACKAGES)

test:
	go test -v -mod=vendor ./...

api:
	docker build -t yximages/yxi-api .

drun:
	docker run -it --rm -p 8090:8090 -v "/var/run/docker.sock:/var/run/docker.sock" yximages/yxi-api

runners:
	cd scripts && ./images.sh -b


push:
	cd scripts && ./images.sh -d

pull:
	cd scripts && ./images.sh -p

pullali:
	cd scripts && ./images.sh -l

pushali:
	cd scripts && ./images.sh -a

clean:
	rm ./main

