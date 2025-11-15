# 🧪 LabFlux Functions Testing Guide

## Quick Test Commands

### 1. **Basic Health Check**
```bash
# Start environment
docker compose up -d postgres
go run cmd/main.go

# Test in another terminal
curl http://localhost:8080/health
curl http://localhost:8080/api/ping
```

### 2. **Generate Test JWT Token**
```bash
go run generate_test_token.go
```

### 3. **Test Laboratory Functions**
```bash
./test_with_auth.sh
```

## Manual Testing Steps

### **Step 1: Start Local Environment**
```bash
cd neva-functions
docker compose up -d postgres
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
export JWT_SECRET="your-local-jwt-secret-key"
go run cmd/main.go
```

### **Step 2: Generate JWT Token**
```bash
# In another terminal
go run generate_test_token.go
```
Copy the token from the output.

### **Step 3: Test Protected Endpoints**

#### **Create Laboratory**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Test Lab",
    "description": "Testing laboratory creation",
    "address": "123 Science Avenue"
  }' \
  http://localhost:8080/api/laboratories
```

#### **Get Laboratory**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     http://localhost:8080/api/laboratories/1
```

#### **Create Invitation**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "colleague@example.com",
    "role": "Member"
  }' \
  http://localhost:8080/api/laboratories/1/invitations
```

#### **List Invitations**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     http://localhost:8080/api/laboratories/1/invitations
```

## Expected Results

### **✅ Should Work:**
- Health check: `{"status": "ok"}`
- API ping: `{"message": "pong"}`
- JWT token generation
- Database table creation

### **❌ Expected Failures (Normal):**
- Laboratory creation: "User not found" 
- Most protected endpoints: Need users in database first

### **🔧 Why Some Tests Fail:**
The JWT token contains `user_id: 1` but no users exist in the database yet. This is expected! We need to implement authentication endpoints first.

## Database Inspection

### **Check Created Tables**
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "\\dt"
```

### **View Table Structure**
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "\\d laboratories"
```

### **Insert Test User (Manual)**
```bash
docker exec neva-functions-postgres-1 psql -U dev -d neva_local -c "
INSERT INTO users (email, password_hash, laboratory_id, role, created_at, updated_at) 
VALUES ('test@labflux.com', 'dummy_hash', NULL, 'Admin', NOW(), NOW());
"
```

## Troubleshooting

### **Port 5433 Already in Use**
```bash
docker compose down
docker compose up -d postgres
```

### **Go Build Issues on macOS**
Use Docker instead:
```bash
docker compose up functions
```

### **Database Connection Errors**
```bash
# Check if PostgreSQL is running
docker compose ps
docker compose logs postgres
```

### **JWT Token Issues**
- Ensure JWT_SECRET matches in both generation and server
- Check token expiration (24 hours)
- Verify token format: "Bearer TOKEN_STRING"

## Next Steps

1. **Implement Auth Endpoints**: `functions/auth.go`
2. **Create User Registration**: Allow creating users
3. **Test Full Flow**: Register → Login → Create Lab → Invite
4. **Add More Endpoints**: Products, Equipment, etc.

## Cost Tracking

**Current Local Development**: $0
**Target Production Cost**: $1-3/month for 50 users