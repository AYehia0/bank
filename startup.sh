#!/bin/sh

set -e

echo "Running migrations"
ls -la
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Starting the server"
exec "$@"
