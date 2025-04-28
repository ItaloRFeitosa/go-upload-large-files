# Build stage
FROM golang:1.24.2-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache build-base git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . ./

# Build the application
RUN go build -o /out/api ./cmd/api

# Final stage
FROM alpine:3.21

# Install only necessary CA certificates (for HTTPS)
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /out/api .

# Expose the port (optional, depends on your app)
# EXPOSE 8080

# Run the binary
ENTRYPOINT ["./api"]
