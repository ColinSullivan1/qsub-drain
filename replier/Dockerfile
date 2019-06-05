FROM golang:1.9.0-alpine3.6 AS builder
WORKDIR $GOPATH/src/github.com/nats-io/demos/reqrepl-graph/replier
COPY . .
RUN apk add --update git
RUN CGO_ENABLED=0 go build -installsuffix cgo -o /replier main.go

FROM alpine:3.6
COPY --from=builder /replier /replier
ENTRYPOINT ["/replier"]
CMD ["-h"]
