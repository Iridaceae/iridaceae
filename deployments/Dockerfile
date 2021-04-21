# Build stage
FROM golang:alpine AS builder
ENV GO111MODULE=on

COPY . /iris/src
WORKDIR /iris/src

RUN apk add git ca-certificates
RUN go mod download

WORKDIR /iris/src/cmd/tensrose

# now we build
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/tensrose -v


FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/tensrose ./
COPY ./internal ./internal
COPY ./.env ./

ENTRYPOINT ["./tensrose"]
