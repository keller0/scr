GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: test
test:
	go test

.PHONY: build
build:fmt
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o main cmd/apiServer/main.go

vet:
	go vet $(PACKAGES)

buildapi:
	docker build -t keller0/yxi-api .

dbuild:
	docker run -it --rm -v `pwd`:/go/src/github.com/keller0/yxi.io \
	-w /go/src/github.com/keller0/yxi.io golang:1.12 \
	go build -ldflags '-w -s' -o main cmd/apiServer/main.go

buildimages:
	cd scripts && ./images.sh -b

push2ali:
	cd scripts && ./images.sh -a

push2dh:
	cd scripts && ./images.sh -d

clean:
	rm ./main

