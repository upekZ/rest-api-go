 FROM golang:1.24 AS builder

 WORKDIR /app

 # Copy go.mod and go.sum files
 COPY go.mod go.sum ./

 # Download dependencies
 RUN go mod download

 # Copy the source code
 COPY . .

 # Build the application
 RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

 # Use a smaller base image for the final container
 FROM alpine:latest

 # Set the working directory
 WORKDIR /app

 # Copy the binary from the builder stage
 COPY --from=builder /app/main .

 # Expose the port your REST API uses (replace 8080 with your port)
 EXPOSE 3000

 # Command to run the application
 CMD ["./main"]