FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install timezone data
RUN apk add --no-cache tzdata

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/server

# Create a minimal image from alpine
FROM alpine:latest

# Copy CA certificates for HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/bin/sh", "-c", "/app/server"]
