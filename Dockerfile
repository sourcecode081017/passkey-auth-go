# Use the official Golang image for building the application
FROM golang:1.24.0 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/

# Use a minimal image with a newer glibc version for the final container
FROM debian:bookworm-slim

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the backend port
EXPOSE 8080

# Run the compiled Go binary
CMD ["./main"]