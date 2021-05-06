# Build stage
FROM golang:alpine AS builder
COPY . /iris/src
WORKDIR /iris/src

FROM builder as concertina-builder
WORKDIR /iris/src/cmd/concertina-test
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/concertina-test -v

FROM alpine:latest as concertina-runner
WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY --from=builder /iris/src/internal internal/
COPY --from=concertina-builder /go/bin/concertina-test .

CMD exec /app/concertina-test

FROM builder as iridaceae-builder
WORKDIR /iris/src/cmd/iridaceae-server
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/iridaceae-server -v

FROM alpine:latest as iridaceae-runner
WORKDIR /app

RUN apk --no-cache add git ca-certificates

COPY --from=builder /iris/src/internal internal/
COPY --from=iridaceae-builder /go/bin/iridaceae-server .

CMD exec /app/iridaceae-server
