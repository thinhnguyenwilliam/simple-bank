# simple-bank/Dockerfile

# =========================
# Build Stage
# =========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple-bank main.go


# =========================
# Run Stage
# =========================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/simple-bank .
COPY app.env .

EXPOSE 8084

CMD ["./simple-bank"]
