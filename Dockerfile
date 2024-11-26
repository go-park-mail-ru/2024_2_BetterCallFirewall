FROM golang:alpine AS build

WORKDIR /test

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod vendor

COPY . .

RUN go build cmd/main.go

FROM alpine:latest

WORKDIR /test

COPY .env .

COPY --from=build /test/main /test/main

CMD ["./main"]