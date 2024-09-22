FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build cmd/main.go

FROM alpine:latest

WORKDIR /app

EXPOSE 8080

COPY .env .

COPY --from=build /app/main /app/main

CMD ["./main"]