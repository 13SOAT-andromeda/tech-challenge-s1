FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/api

# Use a minimal image for running the app
FROM alpine:latest

WORKDIR /app

# Install curl and ca-certificates for healthcheck HTTP probing
RUN apk add --no-cache ca-certificates curl

# Copy the built binary from builder stage
COPY --from=builder /app/main .
# Copy generated or provided swagger files so the app can serve them
COPY --from=builder /app/swagger ./swagger

EXPOSE 8080

CMD ["./main"]