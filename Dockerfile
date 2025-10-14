FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/api

# Use a minimal image for running the app
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]