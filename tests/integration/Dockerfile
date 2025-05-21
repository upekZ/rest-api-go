FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache bash

CMD ["sh", "-c", "go test -v . -tags=integration -count=1"]