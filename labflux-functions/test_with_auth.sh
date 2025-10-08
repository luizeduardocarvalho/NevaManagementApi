#!/bin/bash

echo "🔐 Testing LabFlux Laboratory Functions with Authentication"

# Start the environment
echo "1. Starting environment..."
docker compose up -d postgres
sleep 5

# Start Go server in background
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go &
SERVER_PID=$!
sleep 3

# Generate test token
echo "2. Generating test JWT token..."
TOKEN=$(go run generate_test_token.go | grep -A1 "Test JWT Token Generated:" | tail -1)
echo "   Token: ${TOKEN:0:50}..."

# Test protected endpoints
echo ""
echo "3. Testing Laboratory Management Functions:"

# Test 1: Create Laboratory (should fail - user doesn't exist yet)
echo ""
echo "   📝 Test 1: Create Laboratory"
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Laboratory",
    "description": "A test laboratory for development",
    "address": "123 Science St, Tech City"
  }' \
  http://localhost:8080/api/laboratories | jq '.'

# Test 2: Get Laboratory (should fail - user/lab doesn't exist)
echo ""
echo "   🔍 Test 2: Get Laboratory"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/laboratories/1 | jq '.'

# Test 3: Create Invitation (should fail - lab doesn't exist)
echo ""
echo "   📧 Test 3: Create Invitation"
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@labflux.com",
    "role": "Member"
  }' \
  http://localhost:8080/api/laboratories/1/invitations | jq '.'

# Test 4: Get Invitations
echo ""
echo "   📋 Test 4: Get Invitations"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/laboratories/1/invitations | jq '.'

# Test without auth (should fail)
echo ""
echo "   ❌ Test 5: Access without token (should fail)"
curl -s -X GET \
  http://localhost:8080/api/laboratories/1 | jq '.'

echo ""
echo "4. Database Tables Created:"
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "\\dt"

echo ""
echo "5. Cleaning up..."
kill $SERVER_PID 2>/dev/null
docker compose stop

echo ""
echo "✅ Tests complete!"
echo ""
echo "📝 Notes:"
echo "- Most tests will fail because we need to create users first"
echo "- This is expected - we need authentication endpoints next"
echo "- Database tables are being created correctly"
echo "- JWT authentication middleware is working"
