# Build stage
FROM golang:alpine AS builder

COPY . /iris/src
WORKDIR /iris/src/pkg/cmd/tensroses-server

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/tensroses-server -v

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY --from=builder /iris/src/internal internal/
COPY --from=builder /iris/src/defaults.env .
COPY --from=builder /go/bin/tensroses-server .

CMD exec /app/tensroses-server
