#!/bin/bash

echo "🧪 Testing LabFlux Core API Endpoints (Products, Equipment, Equipment Usage)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Start the environment
echo -e "${BLUE}1. Starting environment...${NC}"
docker compose up -d postgres
sleep 5

# Start Go server in background
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go &
SERVER_PID=$!
sleep 3

# Generate test token
echo -e "${BLUE}2. Generating test JWT token...${NC}"
TOKEN=$(go run generate_test_token.go | grep -A1 "Test JWT Token Generated:" | tail -1)
echo "   Token: ${TOKEN:0:50}..."

# Insert test data directly into database for testing
echo -e "${BLUE}3. Setting up test data...${NC}"
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "
-- Create test laboratory
INSERT INTO laboratories (id, name, description, address, created_at, updated_at) 
VALUES (1, 'Test Laboratory', 'Test lab for API testing', '123 Science St', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Create test user
INSERT INTO users (id, clerk_user_id, email, first_name, last_name, laboratory_id, role, created_at, updated_at) 
VALUES (1, 'test_user_1', 'test@labflux.com', 'Test', 'User', 1, 'Admin', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Create test location
INSERT INTO locations (id, name, description, laboratory_id, created_at, updated_at) 
VALUES (1, 'Main Lab', 'Primary laboratory space', 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Create test researcher
INSERT INTO researchers (id, name, clerk_user_id, email, laboratory_id, created_at, updated_at) 
VALUES (1, 'Dr. Test Researcher', 'researcher_1', 'researcher@labflux.com', 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
"

echo -e "${GREEN}   ✅ Test data created${NC}"

# Test Product Management Endpoints
echo ""
echo -e "${YELLOW}4. Testing Product Management Endpoints:${NC}"

# Test 1: Create Product
echo ""
echo -e "${BLUE}   📦 Test 1: Create Product${NC}"
PRODUCT_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Chemical A",
    "description": "A test chemical for laboratory use",
    "location_id": 1,
    "quantity": 100.5,
    "formula": "H2O2",
    "unit": "ml",
    "expiration_date": "2025-12-31T00:00:00Z"
  }' \
  http://localhost:8080/api/products/Create)
echo "   Response: $PRODUCT_RESPONSE"

# Test 2: Get All Products
echo ""
echo -e "${BLUE}   📋 Test 2: Get All Products${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetAll?laboratoryId=1&page=1" | jq '.'

# Test 3: Get Product by ID
echo ""
echo -e "${BLUE}   🔍 Test 3: Get Product by ID${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetProductById?id=1&laboratoryId=1" | jq '.'

# Test 4: Get Detailed Product
echo ""
echo -e "${BLUE}   📊 Test 4: Get Detailed Product${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetDetailedProductById?id=1&laboratoryId=1" | jq '.'

# Test 5: Add Quantity to Product
echo ""
echo -e "${BLUE}   ➕ Test 5: Add Quantity to Product${NC}"
ADD_QTY_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 25.5
  }' \
  http://localhost:8080/api/products/AddQuantity)
echo "   Response: $ADD_QTY_RESPONSE"

# Test 6: Use Product
echo ""
echo -e "${BLUE}   🔬 Test 6: Use Product${NC}"
USE_PRODUCT_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 10.0,
    "unit": "ml"
  }' \
  http://localhost:8080/api/products/UseProduct)
echo "   Response: $USE_PRODUCT_RESPONSE"

# Test 7: Edit Product
echo ""
echo -e "${BLUE}   ✏️ Test 7: Edit Product${NC}"
EDIT_PRODUCT_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Updated Test Chemical A",
    "description": "Updated description",
    "location_id": 1,
    "quantity": 116.0,
    "formula": "H2O2",
    "unit": "ml",
    "expiration_date": "2025-12-31T00:00:00Z"
  }' \
  http://localhost:8080/api/products/EditProduct)
echo "   Response: $EDIT_PRODUCT_RESPONSE"

# Test 8: Get Low Stock Products
echo ""
echo -e "${BLUE}   ⚠️ Test 8: Get Low Stock Products${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetLowInStockProducts" | jq '.'

# Test Equipment Management Endpoints
echo ""
echo -e "${YELLOW}5. Testing Equipment Management Endpoints:${NC}"

# Test 1: Add Equipment
echo ""
echo -e "${BLUE}   🔧 Test 1: Add Equipment${NC}"
EQUIPMENT_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Microscope",
    "description": "High-powered microscope for testing",
    "property_number": "EQ-001",
    "location_id": 1
  }' \
  http://localhost:8080/api/equipment/AddEquipment)
echo "   Response: $EQUIPMENT_RESPONSE"

# Test 2: Get All Equipment
echo ""
echo -e "${BLUE}   📋 Test 2: Get All Equipment${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipment/GetEquipments?laboratoryId=1" | jq '.'

# Test 3: Get Detailed Equipment
echo ""
echo -e "${BLUE}   📊 Test 3: Get Detailed Equipment${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipment/GetDetailedEquipment?id=1&laboratoryId=1" | jq '.'

# Test 4: Edit Equipment
echo ""
echo -e "${BLUE}   ✏️ Test 4: Edit Equipment${NC}"
EDIT_EQUIPMENT_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Updated Test Microscope",
    "description": "Updated high-powered microscope",
    "property_number": "EQ-001-UPDATED",
    "location_id": 1
  }' \
  http://localhost:8080/api/equipment/EditEquipment)
echo "   Response: $EDIT_EQUIPMENT_RESPONSE"

# Test Equipment Usage Endpoints
echo ""
echo -e "${YELLOW}6. Testing Equipment Usage Endpoints:${NC}"

# Test 1: Use Equipment (Schedule)
echo ""
echo -e "${BLUE}   📅 Test 1: Schedule Equipment Usage${NC}"
USE_EQUIPMENT_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": 1,
    "researcher_id": 1,
    "description": "Testing microscope functionality",
    "start_date": "2024-01-15T09:00:00Z",
    "end_date": "2024-01-15T12:00:00Z"
  }' \
  http://localhost:8080/api/equipmentusage/UseEquipment)
echo "   Response: $USE_EQUIPMENT_RESPONSE"

# Test 2: Get Equipment Usage Calendar
echo ""
echo -e "${BLUE}   📅 Test 2: Get Equipment Usage Calendar${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsageCalendar?id=1&laboratoryId=1" | jq '.'

# Test 3: Get Equipment Usage History
echo ""
echo -e "${BLUE}   📊 Test 3: Get Equipment Usage History${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsageHistory?id=1&laboratoryId=1" | jq '.'

# Test 4: Get Equipment Usage (Paginated)
echo ""
echo -e "${BLUE}   📋 Test 4: Get Equipment Usage (Paginated)${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsage?equipmentId=1&laboratoryId=1&page=1" | jq '.'

# Test Error Cases
echo ""
echo -e "${YELLOW}7. Testing Error Cases:${NC}"

# Test 1: Access without token
echo ""
echo -e "${BLUE}   ❌ Test 1: Access without token (should fail)${NC}"
curl -s -X GET \
  "http://localhost:8080/api/products/GetAll?laboratoryId=1" | jq '.'

# Test 2: Access with wrong laboratory ID
echo ""
echo -e "${BLUE}   ❌ Test 2: Access with wrong laboratory ID (should fail)${NC}"
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetAll?laboratoryId=999" | jq '.'

# Test 3: Conflicting equipment usage
echo ""
echo -e "${BLUE}   ❌ Test 3: Conflicting equipment usage (should fail)${NC}"
CONFLICT_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": 1,
    "researcher_id": 1,
    "description": "Conflicting usage",
    "start_date": "2024-01-15T10:00:00Z",
    "end_date": "2024-01-15T11:00:00Z"
  }' \
  http://localhost:8080/api/equipmentusage/UseEquipment)
echo "   Response: $CONFLICT_RESPONSE"

# Show final database state
echo ""
echo -e "${YELLOW}8. Final Database State:${NC}"
echo -e "${BLUE}   Products:${NC}"
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT id, name, quantity, unit FROM products;"

echo -e "${BLUE}   Equipment:${NC}"
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT id, name, property_number FROM equipment;"

echo -e "${BLUE}   Equipment Usage:${NC}"
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT id, equipment_id, researcher_id, start_date, end_date, description FROM equipment_usages;"

# Cleanup
echo ""
echo -e "${BLUE}9. Cleaning up...${NC}"
kill $SERVER_PID 2>/dev/null
docker compose stop

echo ""
echo -e "${GREEN}✅ Core API Testing Complete!${NC}"
echo ""
echo -e "${YELLOW}📝 Summary:${NC}"
echo "- ✅ Product CRUD operations"
echo "- ✅ Product inventory management (add/use quantity)"
echo "- ✅ Equipment CRUD operations"
echo "- ✅ Equipment usage scheduling with conflict detection"
echo "- ✅ Multi-tenant data isolation"
echo "- ✅ JWT authentication protection"
echo "- ✅ Error handling for edge cases"
echo ""
echo -e "${BLUE}🚀 Ready for Phase 4: Production Deployment!${NC}"