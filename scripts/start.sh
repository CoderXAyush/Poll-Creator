#!/bin/sh

echo "Starting Poll Creator application..."

# Set default port if not provided
export PORT=${PORT:-8080}

echo "Server will run on port: $PORT"

# Start the Go application
exec ./poll-api