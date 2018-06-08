FROM golang:1.8 as builder

WORKDIR /go/src/github.com/keller0/yxi-back/
COPY . .

RUN go get -d -v ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o main

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/keller0/yxi-back/main ./main
ENV YXI_BACK_PORT=":8090"
ENV GIN_MODE="debug"
ENV GIN_LOG_PATH="/var/log/yxi/api.log"
ENV YXI_BACK_KEY="secretkey"
ENV YXI_BACK_MYSQL_ADDR="mariadb:3306"
ENV YXI_BACK_MYSQL_NAME="yxi"
ENV YXI_BACK_MYSQL_USER="root"
ENV YXI_BACK_MYSQL_PASS="111"
CMD ["./main"]
