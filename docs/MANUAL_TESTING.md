# 🧪 Manual Testing Guide - Core API Endpoints

## Prerequisites

### 1. Start the Environment
```bash
cd /Users/luizcarvalho/labflux/NevaManagementApi/labflux-functions
docker compose up -d postgres
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go
```

### 2. Generate JWT Token (in another terminal)
```bash
cd /Users/luizcarvalho/labflux/NevaManagementApi/labflux-functions
TOKEN=$(go run generate_test_token.go | grep -A1 "Test JWT Token Generated:" | tail -1)
echo "Token: $TOKEN"
```

### 3. Setup Test Data
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "
INSERT INTO laboratories (id, name, description, address, created_at, updated_at) 
VALUES (1, 'Test Laboratory', 'Test lab for API testing', '123 Science St', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, clerk_user_id, email, first_name, last_name, laboratory_id, role, created_at, updated_at) 
VALUES (1, 'test_user_1', 'test@labflux.com', 'Test', 'User', 1, 'Admin', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO locations (id, name, description, laboratory_id, created_at, updated_at) 
VALUES (1, 'Main Lab', 'Primary laboratory space', 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO researchers (id, name, clerk_user_id, email, laboratory_id, created_at, updated_at) 
VALUES (1, 'Dr. Test Researcher', 'researcher_1', 'researcher@labflux.com', 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
"
```

---

## 📦 Product Management Testing

### Create Product
```bash
curl -X POST \
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
  http://localhost:8080/api/products/Create
```

### Get All Products
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetAll?laboratoryId=1&page=1"
```

### Get Product by ID
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetProductById?id=1&laboratoryId=1"
```

### Get Detailed Product
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetDetailedProductById?id=1&laboratoryId=1"
```

### Add Quantity to Product
```bash
curl -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 25.5
  }' \
  http://localhost:8080/api/products/AddQuantity
```

### Use Product (Consume Inventory)
```bash
curl -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 10.0,
    "unit": "ml"
  }' \
  http://localhost:8080/api/products/UseProduct
```

### Edit Product
```bash
curl -X PATCH \
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
  http://localhost:8080/api/products/EditProduct
```

### Get Low Stock Products
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetLowInStockProducts"
```

---

## 🔧 Equipment Management Testing

### Add Equipment
```bash
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Microscope",
    "description": "High-powered microscope for testing",
    "property_number": "EQ-001",
    "location_id": 1
  }' \
  http://localhost:8080/api/equipment/AddEquipment
```

### Get All Equipment
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipment/GetEquipments?laboratoryId=1"
```

### Get Detailed Equipment
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipment/GetDetailedEquipment?id=1&laboratoryId=1"
```

### Edit Equipment
```bash
curl -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Updated Test Microscope",
    "description": "Updated high-powered microscope",
    "property_number": "EQ-001-UPDATED",
    "location_id": 1
  }' \
  http://localhost:8080/api/equipment/EditEquipment
```

---

## 📅 Equipment Usage Testing

### Schedule Equipment Usage
```bash
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": 1,
    "researcher_id": 1,
    "description": "Testing microscope functionality",
    "start_date": "2024-01-15T09:00:00Z",
    "end_date": "2024-01-15T12:00:00Z"
  }' \
  http://localhost:8080/api/equipmentusage/UseEquipment
```

### Get Equipment Usage Calendar
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsageCalendar?id=1&laboratoryId=1"
```

### Get Equipment Usage History
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsageHistory?id=1&laboratoryId=1"
```

### Get Equipment Usage (Paginated)
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/equipmentusage/GetEquipmentUsage?equipmentId=1&laboratoryId=1&page=1"
```

### Test Conflicting Usage (Should Fail)
```bash
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": 1,
    "researcher_id": 1,
    "description": "Conflicting usage",
    "start_date": "2024-01-15T10:00:00Z",
    "end_date": "2024-01-15T11:00:00Z"
  }' \
  http://localhost:8080/api/equipmentusage/UseEquipment
```

---

## 🚨 Error Testing

### Access Without Token (Should Fail)
```bash
curl "http://localhost:8080/api/products/GetAll?laboratoryId=1"
```

### Access Wrong Laboratory (Should Fail)
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetAll?laboratoryId=999"
```

### Invalid Product ID (Should Fail)
```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/products/GetProductById?id=999&laboratoryId=1"
```

---

## 🔍 Database Inspection

### Check All Tables
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "\\dt"
```

### View Products
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT * FROM products;"
```

### View Equipment
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT * FROM equipment;"
```

### View Equipment Usage
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT * FROM equipment_usages;"
```

### View Locations
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT * FROM locations;"
```

### View Researchers
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "SELECT * FROM researchers;"
```

---

## 📊 Expected Results

### ✅ Success Cases
- **Product Creation**: "Successfully created Test Chemical A."
- **Product List**: JSON array with product objects
- **Quantity Add**: "Successfully added 25.50 to product."
- **Product Use**: "Successfully used 10.00ml."
- **Equipment Creation**: "Test Microscope was created successfully."
- **Usage Scheduling**: "Equipment was used successfully."

### ❌ Expected Failures
- **No Token**: `{"error": "Authorization header required"}`
- **Wrong Lab ID**: Empty arrays or "not found" errors
- **Conflicting Usage**: "Equipment is already in use during the specified time period"
- **Insufficient Quantity**: "Insufficient quantity available"

---

## 🧹 Cleanup
```bash
# Stop server (Ctrl+C in server terminal)
docker compose stop
docker compose down
```

---

## 🚀 Automated Testing

### Run Full Test Suite
```bash
./test_core_api.sh
```

### Run Basic Health Checks
```bash
./test_with_auth.sh
```

This comprehensive testing approach validates all the new Core API endpoints work correctly with proper authentication, multi-tenancy, and error handling!