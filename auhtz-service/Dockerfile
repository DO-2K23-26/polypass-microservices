# Stage 1: Build the application
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o authz-service ./main.go

# Stage 2: Create a runtime image with glibc support
FROM debian:bookworm


# Install required libraries (glibc)
RUN apt-get update && apt-get install -y libc6 && apt-get clean


# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/authz-service  .
COPY --from=builder /app/schema.zed  .


# Command to run the application
CMD ["./authz-service"]
