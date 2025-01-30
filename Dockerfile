# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /myapp /app/myapp

# Copy .env file (if needed)
COPY .env .env

# Expose the correct port (change to 50051 if needed)
EXPOSE 50051

# Run the application
CMD ["/app/myapp"]
