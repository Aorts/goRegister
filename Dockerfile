FROM golang:alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goRegister

EXPOSE 8080

# Run
CMD ["/goRegister"]