# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN go build -o policy_service .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/policy_service .
EXPOSE 3000
CMD ["./policy_service"]