# base image
FROM golang:latest as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker ./cmd/api

RUN chmod +x /app/broker

# build tiny image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/broker /app

CMD ["/app/broker"]
