FROM golang:alpine AS builder
RUN apk add --no-cache git && CGO_ENABLED=0 GOOS=linux go get github.com/faghani/ssl-proxy

FROM alpine:latest
RUN apk --no-cache add curl nano ca-certificates git
COPY --from=builder /go/bin/ssl-proxy /go/bin/ssl-proxy
RUN mkdir /go/bin/ssl-proxy/certs
ENTRYPOINT ["/go/bin/ssl-proxy"]