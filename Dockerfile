FROM golang:1.8 as builder

WORKDIR /go/src/github.com/keller0/yxi.io/
COPY . .

RUN go get -d -v ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags '-w -s' -o main cmd/apiServer/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/keller0/yxi.io/main ./main
CMD ["./main"]
