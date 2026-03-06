# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
# Uncomment if you have a config file
# COPY --from=builder /app/config.yaml ./config.yaml
EXPOSE 8080
CMD ["./api"]
