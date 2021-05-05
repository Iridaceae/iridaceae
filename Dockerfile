# Build stage
FROM golang:alpine AS builder

COPY . /iris/src
WORKDIR /iris/src/cmd/iridaceae-server

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/iridaceae-server -v

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY --from=builder /iris/src/internal internal/
COPY --from=builder /go/bin/iridaceae-server .

CMD exec /app/iridaceae-server
