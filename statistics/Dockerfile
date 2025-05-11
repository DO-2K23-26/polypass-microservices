# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o statistics-service

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/statistics-service .

# Copy configuration files if needed
COPY --from=builder /app/config.json .

# Expose the port the service listens on
EXPOSE 8080

# Run the service
CMD ["./statistics-service"]