FROM golang:1.17-alpine3.14 AS builder

RUN apk add --update --no-cache ca-certificates make git curl

WORKDIR /usr/local/src/kube-curiesync-injector

ARG GOPROXY

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/webhook .


FROM alpine:3.15.0

RUN apk add --update --no-cache ca-certificates tzdata bash curl

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/bin/* /usr/local/bin/

CMD webhook
