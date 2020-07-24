# Golang base image
FROM golang:1.14 as go_builder

LABEL stage="ps_builder"

WORKDIR /go/src/github.com/rl404/point-system

COPY . .

WORKDIR /go/src/github.com/rl404/point-system/cmd/point

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o point


# New stage from scratch
FROM alpine:3.10

LABEL maintainer="axel.rl.404@gmail.com"

WORKDIR /docker/bin

COPY --from=go_builder /go/src/github.com/rl404/point-system/cmd/point/point point

CMD ["/docker/bin/point"]