#!/bin/bash

echo "🧪 Testing Clerk Integration"
echo "=============================="

BASE_URL="http://localhost:8080/api"

# Test 1: Webhook endpoint
echo "1. Testing Clerk webhook..."
curl -X POST "$BASE_URL/webhooks/clerk" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user.created",
    "object": "event",
    "data": {
      "id": "user_test123",
      "email_addresses": [{
        "id": "email_123",
        "email_address": "test@example.com",
        "verification": {"status": "verified"}
      }],
      "first_name": "Test",
      "last_name": "User",
      "created_at": 1697000000,
      "updated_at": 1697000000
    }
  }' && echo -e "\n✅ Webhook test completed\n"

# Test 2: Create test laboratory (requires auth bypass for now)
echo "2. Testing basic endpoints..."
curl -X GET "$BASE_URL/ping" && echo -e "\n✅ Ping test completed\n"

# Test 3: Database check
echo "3. Checking database for created user..."
docker exec labflux-functions-postgres-1 psql -U dev -d neva_local -c "SELECT id, clerk_user_id, email, first_name, last_name FROM users WHERE clerk_user_id = 'user_test123';"

echo -e "\n🎉 Clerk integration tests completed!"
echo "Next steps:"
echo "- Set up real Clerk account and get API keys"
echo "- Configure ngrok webhook URL in Clerk dashboard"
echo "- Test with real Clerk JWT tokens"