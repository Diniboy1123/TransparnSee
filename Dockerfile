FROM golang:alpine AS builder

COPY . /src

WORKDIR /src

RUN go build -o /bin/transparnsee -ldflags "-s -w" .

# we cannot use scratch, because we need a ca cert bundle
FROM alpine:latest

COPY config/config.json /app/config/config.json
COPY --from=builder /bin/transparnsee /bin/transparnsee

WORKDIR /app

ENTRYPOINT ["/bin/transparnsee"]