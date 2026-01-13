# Estágio de build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/server

# Estágio de execução
FROM alpine:latest
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
