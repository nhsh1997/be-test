FROM golang:1.13.4-alpine3.10 as builder
WORKDIR /app
RUN apk update && apk upgrade && \
    apk add bash git openssh gcc libc-dev
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build  /app/cmd/server


FROM alpine:3.9

RUN apk add --update ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    echo "Asia/Ho_Chi_Minh" > /etc/timezone && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app/server /app/server

WORKDIR /app
CMD ["/app/server"]
