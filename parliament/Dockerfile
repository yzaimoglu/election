# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server server.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/server .

CMD [ "/app/server" ]