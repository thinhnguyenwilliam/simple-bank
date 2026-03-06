# simple-bank/Dockerfile
# docker build -t simple-bank . --no-cache
# docker run -p 8084:8084 simple-bank
# docker logs <container_id>

# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main .

# runtime stage
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8084

CMD ["./main"]