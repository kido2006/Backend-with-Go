# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go app
RUN go build -o backend ./cmd/api

# build seed
RUN go build -o seed ./cmd/migrate/seed
# Run stage
FROM debian:bullseye-slim

WORKDIR /app

# Copy binary từ builder stage
COPY --from=builder /app/backend .
COPY --from=builder /app/seed .


EXPOSE 8081

CMD ["./backend"]
