GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

all: apiimage runnerimages
	docker run -it --rm -p 8090:8090 yximages/yxi-api

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: dev
dev:fmt vet
	GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags '-w -s' -o main cmd/apiServer/main.go

vet:
	go vet $(PACKAGES)

test:
	go test -v -mod=vendor ./...

apiimage:
	docker build -t yximages/yxi-api .

dbuild: apiimage
	docker run -it --rm -p 8090:8090 yximages/yxi-api

runnerimages:
	cd scripts && ./images.sh -b

push2ali:
	cd scripts && ./images.sh -a

push2dh:
	cd scripts && ./images.sh -d

pullimages:
	cd scripts && ./images.sh -p

pullali:
	cd scripts && ./images.sh -l

clean:
	rm ./main

