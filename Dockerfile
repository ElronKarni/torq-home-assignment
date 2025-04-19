FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ip2country-api ./cmd/main.go

# Create final lightweight image
FROM alpine:latest

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/ip2country-api /app/ip2country-api

# Copy data directory
COPY --from=builder /app/data /app/data

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Set environment variables with sensible defaults for production
ENV IP2COUNTRY_DB_TYPE=csv \
    RATE_LIMIT=100 \
    PORT=8080 \
    CSV_DATA_PATH=/app/data/ip2country.csv

# Expose port
EXPOSE 8080

# Start the application
CMD ["/app/ip2country-api"] 