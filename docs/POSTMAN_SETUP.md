# 📮 Postman Setup Guide for LabFlux API

## 🚀 Quick Start

### 1. Import Collection
1. Open Postman
2. Click **Import** button
3. Drag and drop `LabFlux_Core_API.postman_collection.json` file
4. Collection will appear in your sidebar

### 2. Start Local Environment
```bash
cd /Users/luizcarvalho/labflux/NevaManagementApi/labflux-functions
docker compose up -d postgres
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go
```

### 3. Generate JWT Token
```bash
# In another terminal
cd /Users/luizcarvalho/labflux/NevaManagementApi/labflux-functions
go run generate_test_token.go
```
Copy the generated token.

### 4. Configure Postman Variables
In Postman, click on the collection name → **Variables** tab:

| Variable | Value | Description |
|----------|-------|-------------|
| `base_url` | `http://localhost:8080` | Local server URL |
| `laboratory_id` | `1` | Test laboratory ID |
| `jwt_token` | `your_copied_token` | JWT token from step 3 |
| `invitation_token` | `invitation_uuid` | Invitation token (generated after creating invitation) |

---

## 📋 Collection Structure

### 1. **Health Check**
- **Health Check** - Basic server health
- **API Ping** - API endpoint verification

### 2. **Authentication**
- **Get Current User** - Get authenticated user details
- **Accept Invitation** - Accept laboratory invitation
- **Clerk Webhook** - Webhook for user sync (development)

### 3. **Laboratory Management**
- **Create Laboratory** - Create new laboratory
- **Get Laboratory** - Retrieve laboratory details
- **Create Invitation** - Invite users to laboratory
- **Get Invitations** - List all laboratory invitations

### 4. **Product Management**
- **Create Product** - Add new products to inventory
- **Get All Products** - List products with pagination
- **Get Product by ID** - Get basic product details
- **Get Detailed Product** - Get product with location info
- **Add Quantity to Product** - Increase inventory
- **Use Product** - Consume inventory
- **Edit Product** - Update product details
- **Get Low Stock Products** - Products with low inventory

### 5. **Equipment Management**
- **Add Equipment** - Register new equipment
- **Get All Equipment** - List all laboratory equipment
- **Get Detailed Equipment** - Equipment with usage history
- **Edit Equipment** - Update equipment details

### 6. **Equipment Usage**
- **Schedule Equipment Usage** - Book equipment time slots
- **Get Equipment Usage Calendar** - View all bookings
- **Get Equipment Usage History** - Past usage records
- **Get Equipment Usage (Paginated)** - Paginated usage list
- **Test Conflicting Usage** - Error case testing

### 7. **Error Testing**
- **Access Without Token** - Test authentication requirement
- **Access Wrong Laboratory** - Test multi-tenancy
- **Invalid Product ID** - Test error handling

---

## 🧪 Testing Workflow

### Basic Flow
1. **Health Check** → Verify server is running
2. **Create Laboratory** → Set up test lab
3. **Create Product** → Add test inventory
4. **Add Equipment** → Register test equipment
5. **Schedule Equipment Usage** → Test booking system

### Setup Test Data
Before testing core endpoints, run this SQL to create test data:

```sql
-- Copy this into: docker exec neva-functions-postgres-1 psql -U dev -d neva_local

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
```

---

## 🔧 Environment Variables

### Collection Variables
| Variable | Description | How to Get |
|----------|-------------|------------|
| `base_url` | Server URL | Use `http://localhost:8080` for local |
| `laboratory_id` | Lab ID for testing | Use `1` for test lab |
| `jwt_token` | Authentication token | Run `go run generate_test_token.go` |
| `invitation_token` | Invitation UUID | Create invitation first, copy from response |

### Authorization
The collection uses **Bearer Token** authentication. Most endpoints require authentication except:
- Health checks
- Clerk webhook
- Accept invitation

---

## 📊 Expected Responses

### ✅ Success Examples

**Create Product:**
```
Successfully created Test Chemical A.
```

**Get Products:**
```json
[
  {
    "id": 1,
    "name": "Test Chemical A",
    "quantity": 100.5,
    "unit": "ml"
  }
]
```

**Equipment Usage:**
```
Equipment was used successfully.
```

### ❌ Error Examples

**No Authentication:**
```json
{
  "error": "Authorization header required"
}
```

**Wrong Laboratory:**
```json
[]
```

**Conflicting Usage:**
```
Equipment is already in use during the specified time period
```

---

## 🧹 Cleanup

After testing:
```bash
# Stop server (Ctrl+C)
docker compose stop
docker compose down
```

---

## 🚀 Production Usage

For production testing, update variables:
- `base_url`: Your production URL
- `jwt_token`: Real authentication token
- `laboratory_id`: Actual laboratory ID

The collection is designed to work with both local development and production environments!