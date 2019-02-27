GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")

all: build

.PHONY: clean
fmt:
	$(GOFMT) -w $(GOFILES)

build:fmt
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o run

dbuild:
	docker run -it --rm -v `pwd`:/go/src/github.com/keller0/yxi.io/cmd/ric \
	 -w /go/src/github.com/keller0/yxi.io/cmd/ric golang:1.8  go build -ldflags '-w -s' -o run
.PHONY: gotest
gotest:
	go test
	
vet:
	go vet $(PACKAGES)

clean:
	rm ./run
