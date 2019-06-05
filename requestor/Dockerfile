FROM golang:1.9.0-alpine3.6 AS builder
WORKDIR $GOPATH/src/github.com/nats-io/demos/reqrepl-graph/requestor
COPY . .
RUN apk add --update git
RUN CGO_ENABLED=0 go build -installsuffix cgo -o /requestor main.go

FROM alpine:3.6
COPY --from=builder /requestor /requestor
ENTRYPOINT ["/requestor"]
CMD ["-h"]
