# Dockerfile.test
FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Install any test dependencies (e.g., testify if used)
RUN apk add --no-cache bash

# Run integration tests
CMD ["sh", "-c", "go test -v . -tags=integration -count=1"]