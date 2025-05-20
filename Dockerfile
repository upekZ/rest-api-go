 FROM golang:1.24 AS builder

 WORKDIR /app

 RUN go install github.com/air-verse/air@latest

 COPY go.mod go.sum ./

 RUN go mod download

 COPY . .

 RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

 FROM alpine:latest

 WORKDIR /app

 COPY --from=builder /app/main .

 EXPOSE 3000

 # Command to run the application
 CMD ["./main"]