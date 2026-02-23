# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN go build -o auth_service .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth_service .
EXPOSE 3000
CMD ["./auth_service"]