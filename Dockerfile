FROM golang:1.18 as builder

WORKDIR /root/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -ldflags '-w -s' -o main cmd/apiServer/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /root/main ./main
CMD ["./main"]
