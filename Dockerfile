# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o air-quality-monitor .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/air-quality-monitor .

# Create directory for database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV DEVICE_URL=http://192.168.0.249/json
ENV SERVER_ADDR=:8080

# Run the application
CMD ["./air-quality-monitor", "--server", "${DEVICE_URL}"]
