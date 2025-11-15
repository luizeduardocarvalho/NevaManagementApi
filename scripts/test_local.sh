#!/bin/bash

echo "🚀 Testing LabFlux Functions Local Environment"

# Start database
echo "1. Starting PostgreSQL..."
docker compose up -d postgres

# Wait for database
echo "2. Waiting for database to be ready..."
sleep 10

# Test database connection
echo "3. Testing database connection..."
docker exec neva-functions-postgres-1 pg_isready -U dev -d neva_local

# Build the application
echo "4. Building Go application..."
go build -o labflux-server cmd/main.go

# Start server in background
echo "5. Starting LabFlux server..."
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
./labflux-server &
SERVER_PID=$!

# Wait for server to start
sleep 5

# Test endpoints
echo "6. Testing API endpoints..."

echo "   Testing health endpoint..."
curl -s http://localhost:8080/health | jq '.'

echo "   Testing ping endpoint..."
curl -s http://localhost:8080/api/ping | jq '.'

# Clean up
echo "7. Cleaning up..."
kill $SERVER_PID
rm labflux-server

echo "✅ Local environment test complete!"
