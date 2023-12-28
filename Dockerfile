FROM golang:alpine as builder

WORKDIR /app

RUN apk add git

COPY go.mod .
COPY go.sum .

COPY . .

RUN apk update

RUN apk add gcc libc-dev make


FROM alpine:latest as release



COPY --from=builder /app/main /app/cmd/

RUN chmod +x /app/cmd/main

WORKDIR /app

EXPOSE 8080

CMD ["cmd/main"]

