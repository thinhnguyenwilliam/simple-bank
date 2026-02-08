# =========================
# Build stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (needed for go mod)
RUN apk add --no-cache git

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app main.go

# =========================
# Runtime stage (minimal)
# =========================
FROM scratch

WORKDIR /app

# Copy binary
COPY --from=builder /app/app .

# Copy config if needed
COPY app.env .

# Expose port
EXPOSE 8080

# Run
ENTRYPOINT ["/app/app"]
