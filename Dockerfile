# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate proto files if needed
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /myapp /app/myapp

# Copy .env file (if needed)
COPY .env .env

# Expose the correct port
EXPOSE 50051

# Run the application
CMD ["/app/myapp"]
