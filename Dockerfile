FROM golang:alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/main.go

FROM alpine:latest

WORKDIR /app

EXPOSE 8080

COPY .env .

COPY image .

COPY --from=build /app/main /app/main

CMD ["./main"]