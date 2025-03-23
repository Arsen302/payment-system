FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
# COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/payment-service .

# Expose the application port
EXPOSE 50052

# Run the application
CMD ["./payment-service"] 