FROM golang:alpine AS builder
RUN apk add --no-cache git && CGO_ENABLED=0 GOOS=linux go get github.com/faghani/ssl-proxy

FROM scratch
COPY --from=builder /go/bin/ssl-proxy /go/bin/ssl-proxy
ENTRYPOINT ["/go/bin/ssl-proxy"]