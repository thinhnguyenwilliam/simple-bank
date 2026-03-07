#!/bin/sh
set -e

echo "🚀 Starting simple-bank..."

echo "⏳ Waiting for database..."
while ! nc -z postgres 5432; do
  sleep 1
done

echo "✅ Database is ready"

echo "📦 Running database migrations..."
migrate \
  -path db/migrations \
  -database "$DB_SOURCE" \
  up

echo "▶️ Starting service: $1"

exec "$@"