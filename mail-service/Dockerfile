FROM golang:latest as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

RUN chmod +x /app/mailerApp

# build tiny image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailerApp /app

CMD ["/app/mailerApp"]
