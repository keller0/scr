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
build:deps fmt
	govendor sync
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o main

vet:
	go vet $(PACKAGES)

deps:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)


tDB:
	docker run --rm -d -p 3306:3306 -v `pwd`/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=111 mariadb:10.3
clean:
	rm ./main

