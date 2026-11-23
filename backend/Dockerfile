# Multi-stage build for backend service
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Final stage
FROM alpine:latest

# Create non-root user
RUN adduser -D -u 1000 appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy configuration files
COPY .env .

# Change ownership to non-root user
RUN chown -R appuser:appuser .
USER appuser

# Expose port
EXPOSE 7080

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:7080/health || exit 1

# Run the binary
CMD ["./server"]