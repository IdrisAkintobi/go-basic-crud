# Build stage
FROM golang:1.23.6-alpine AS builder

# Install required tools for geoip download
RUN apk --no-cache add make curl

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download and tidy dependencies
RUN go mod tidy

# Copy source code
COPY . .

# Download GeoIP data
ARG GEO21P_ACCOUNT_ID
ARG GEO21P_LICENSE_KEY
RUN make geoip.download

# Build the application
RUN go build -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy GeoIP data from builder stage
COPY --from=builder /app/geo2ip-data ./geo2ip-data

# Expose port
EXPOSE 3003

# Run the application
CMD ["./main"]