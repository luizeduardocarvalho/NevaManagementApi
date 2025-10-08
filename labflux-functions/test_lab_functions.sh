#!/bin/bash

echo "🧪 Testing LabFlux Laboratory Functions"

# Start the environment
echo "1. Starting PostgreSQL and Go server..."
docker compose up -d postgres
sleep 5

# Start Go server in background
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go &
SERVER_PID=$!
sleep 3

echo "2. Server started on http://localhost:8080"

# Create a test JWT token (we'll need this for protected endpoints)
echo "3. Testing public endpoints first..."

# Test health check
echo "   Health check:"
curl -s http://localhost:8080/health | jq '.'

# Test ping
echo "   API ping:"
curl -s http://localhost:8080/api/ping | jq '.'

echo ""
echo "4. For protected endpoints, you need a JWT token."
echo "   Since we don't have auth endpoints yet, here's how to test:"
echo ""
echo "   Option 1: Create a test JWT token manually"
echo "   Option 2: Use the token generator script"
echo ""

# Clean up
echo "5. Cleaning up..."
kill $SERVER_PID 2>/dev/null
docker compose stop

echo "✅ Test complete!"
echo ""
echo "Next steps:"
echo "- Run './generate_test_token.go' to create a JWT token"
echo "- Use the token to test protected laboratory endpoints"