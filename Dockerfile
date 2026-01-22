# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build the application from cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-server ./cmd/server

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/bin/api-server /api-server
COPY --from=builder /app/migrations /migrations

EXPOSE 8080
CMD ["/api-server"]