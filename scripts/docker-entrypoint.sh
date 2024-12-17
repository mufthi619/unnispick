#!/bin/sh
set -e

# Wait for database to be ready
echo "Waiting for database..."
for i in $(seq 1 30); do
    if pg_isready -h db -U k-style; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database... $i/30"
    sleep 1
done

if [ $i = 30 ]; then
    echo "Database connection timeout"
    exit 1
fi

# Run migrations
echo "Running database migrations..."
./main -command=migrate -db="${DATABASE_URL}" up

# Start the application
echo "Starting application..."
exec ./main -command=api