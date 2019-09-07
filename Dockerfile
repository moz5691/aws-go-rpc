FROM golang:alpine as builder

RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh

LABEL maintainer="Chan Ahn<moz5691@gmail.com>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main -v cmd/moviesserver/main.go
EXPOSE 8080

CMD ["./main"]