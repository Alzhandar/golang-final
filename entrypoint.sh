#!/bin/sh
# filepath: /Users/lzandaribaev/Desktop/go final/golang-final/entrypoint.sh

set -e

# Ждем запуска PostgreSQL - увеличиваем таймаут и добавляем больше информации
echo "Waiting for PostgreSQL to start..."
for i in $(seq 1 30); do
  echo "Checking PostgreSQL connection... (attempt $i)"
  if nc -z postgres 5432; then
    echo "PostgreSQL is up and running!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "Timeout waiting for PostgreSQL"
    exit 1
  fi
  sleep 1
done

# Запускаем миграции
echo "Running migrations..."
# Здесь используем опцию -verbose, чтобы видеть больше информации
migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/restaurant_db?sslmode=disable" up -verbose || echo "Migration failed but continuing..."

# Запускаем приложение
echo "Starting application..."
exec "$@"