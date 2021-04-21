# Build stage
FROM golang:alpine AS builder

COPY . /iris/src
WORKDIR /iris/src/cmd/tensrose

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/tensrose -v

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY --from=builder /iris/src/internal internal/
COPY --from=builder /go/bin/tensrose .
COPY .env ./

CMD exec /app/tensrose
