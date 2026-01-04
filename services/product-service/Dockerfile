# =========================
# 1. Build stage
# =========================
FROM golang:1.24-alpine AS builder

# Cài git (chỉ cần ở build stage)
RUN apk add --no-cache git

WORKDIR /app

# Copy module files trước để cache dependency
COPY go.mod go.sum ./

# Download dependencies (cache layer này)
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server

# =========================
# 2. Run stage
# =========================
FROM scratch

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
