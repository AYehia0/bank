#!/bin/sh

set -e

echo "Running migrations"
#source /app/config.env
echo $DB_SOURCE

/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Starting the server"
exec "$@"
