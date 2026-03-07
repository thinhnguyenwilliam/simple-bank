# simple-bank/Dockerfile
# docker build -t simple-bank . --no-cache
# docker run -p 8084:8084 simple-bank
# docker logs 6ceb32bee33c

# build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache curl

# install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
    | tar xvz && mv migrate /usr/local/bin/

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o /app/main .
RUN go build -o /app/worker ./cmd/worker

# runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates netcat-openbsd

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/worker .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

COPY db/migrations db/migrations
COPY start.sh .
COPY app.env .

RUN chmod +x start.sh

EXPOSE 8084

ENTRYPOINT ["./start.sh"]