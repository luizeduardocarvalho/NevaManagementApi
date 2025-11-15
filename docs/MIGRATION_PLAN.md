# LabFlux Serverless Migration Plan

## 🎯 **Project Overview**

**Goal**: Transform NevaManagement (now LabFlux) from ASP.NET Core + PostgreSQL hosting to serverless architecture for 50 users.

**Current Status**: ✅ Local development environment ready with Go + Chi + PostgreSQL

**Expected Cost Reduction**: From ~$200-300/month to $1-10/month (95%+ savings)

---

## 📋 **Migration Strategy: Gradual Replacement**

### **Phase 1: Foundation & New Features** ✅ COMPLETED
- [x] Set up Go project structure for GCP Cloud Functions
- [x] Configure local development environment with Docker Compose
- [x] Create Chi router setup and middleware
- [x] Implement database connection and GORM models
- [x] Set up JWT authentication middleware
- [x] Implement laboratory onboarding functions
- [x] Test local development environment

**Current State**: 
- Local environment running on port 8080
- PostgreSQL on port 5433 (no conflicts)
- Laboratory management endpoints ready
- Multi-tenant data isolation implemented

### **Phase 2: Authentication System** ✅ COMPLETED
- [x] Create authentication endpoints in Go (login, register)
- [x] Implement password hashing and JWT token generation
- [x] Add user management functions
- [x] Test authentication flow with existing frontend
- [x] Create invitation acceptance endpoint
- [x] Fix JWT token validation issues (secret key synchronization)

### **Phase 3: Core API Migration** ✅ COMPLETED
- [x] Convert high-traffic endpoints (Products, Equipment)
- [x] Migrate Equipment Usage tracking
- [x] Implement Location management
- [x] Implement Researcher management
- [ ] Convert Container management (future)

### **Phase 4: Deployment & Production** 🔄 PENDING
- [ ] Set up Neon database (production)
- [ ] Deploy functions to GCP Cloud Functions
- [ ] Configure custom domain and CORS
- [ ] Update frontend to use new endpoints
- [ ] Performance testing and optimization

### **Phase 5: Full Migration** 🔄 PENDING
- [ ] Migrate remaining endpoints
- [ ] Data migration from current database
- [ ] Switch frontend completely to Go API
- [ ] Deprecate ASP.NET Core API
- [ ] Monitor and optimize costs

---

## 🏗️ **Technical Architecture**

### **Current Setup (Working)**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │───▶│   Go Functions  │───▶│   PostgreSQL    │
│   (existing)    │    │   (Chi + GORM)  │    │   (port 5433)   │
│   Port: 3000    │    │   Port: 8080    │    │   Docker        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Target Production Architecture**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │───▶│  GCP Functions  │───▶│  Neon Database  │
│   (Vercel)      │    │  (Go + Chi)     │    │  (PostgreSQL)   │
│   Free          │    │  $1-3/month     │    │  Free Tier      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 🛠️ **Development Workflow**

### **Local Development Commands**
```bash
# Start local environment
cd neva-functions
docker compose up -d postgres

# Run Go server
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
go run cmd/main.go

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/ping
```

### **Project Structure**
```
labflux-functions/
├── cmd/
│   └── main.go              # Local development server
├── functions/
│   ├── laboratory.go        # ✅ Laboratory management
│   ├── auth.go             # 🔄 Authentication (next)
│   ├── products.go         # 🔄 Product management
│   └── equipment.go        # 🔄 Equipment management
├── pkg/
│   ├── database/           # ✅ Database connection
│   ├── middleware/         # ✅ Auth middleware
│   └── models/             # ✅ GORM models
├── internal/
│   └── router.go           # ✅ Chi router setup
├── docker-compose.yml      # ✅ Local development
└── Dockerfile             # ✅ Container build
```

---

## 📊 **Cost Analysis**

### **Current Hosting Costs (Estimated)**
- **API Hosting**: $50-100/month
- **Database**: $50-150/month
- **Frontend**: $20-50/month
- **Total**: $120-300/month

### **Target Serverless Costs (50 users, ~30k requests/month)**
- **GCP Cloud Functions**: $1-3/month (2M free invocations)
- **Neon Database**: $0/month (free tier: 0.5GB, 100 compute hours)
- **Vercel Frontend**: $0/month (free tier)
- **Total**: $1-3/month

**Savings**: 95-99% cost reduction

---

## 🔧 **Key Technologies**

### **Backend Stack**
- **Language**: Go 1.21+
- **Framework**: Chi (lightweight HTTP router)
- **Database**: PostgreSQL (GORM ORM)
- **Authentication**: JWT tokens
- **Deployment**: GCP Cloud Functions

### **Database Strategy**
- **Development**: Docker PostgreSQL (port 5433)
- **Production**: Neon (serverless PostgreSQL)
- **Migration**: GORM auto-migrate
- **Multi-tenancy**: LaboratoryID filtering

### **Authentication Flow**
- JWT tokens with laboratory claims
- Middleware-based route protection
- Multi-tenant data isolation
- Secure invitation system (7-day expiry)

---

## 🎯 **Current Implementation Status**

### **✅ Completed Features**
1. **Laboratory Management**
   - `POST /api/laboratories` - Create laboratory
   - `GET /api/laboratories/{id}` - Get laboratory
   - `POST /api/laboratories/{id}/invitations` - Create invitation
   - `GET /api/laboratories/{id}/invitations` - List invitations

2. **Authentication System**
   - `POST /api/auth/register` - User registration
   - `POST /api/auth/login` - User login  
   - `POST /api/auth/invite/{token}` - Accept invitation
   - `GET /api/auth/me` - Get current user
   - JWT token generation/validation (fixed secret key sync)

3. **Infrastructure**
   - Local development environment
   - Database connection and migrations
   - JWT authentication middleware
   - Docker Compose setup

4. **Models**
   - Laboratory entity
   - LaboratoryInvitation entity
   - User entity with laboratory relationship

### **✅ Completed (Phase 3)**
1. **Core API Migration**
   ```go
   // Product Management
   GET  /api/products/GetAll
   GET  /api/products/GetDetailedProductById
   GET  /api/products/GetProductById
   POST /api/products/Create
   PATCH /api/products/AddQuantity
   PATCH /api/products/UseProduct
   PATCH /api/products/EditProduct
   GET  /api/products/GetLowInStockProducts

   // Equipment Management
   GET  /api/equipment/GetEquipments
   POST /api/equipment/AddEquipment
   GET  /api/equipment/GetDetailedEquipment
   PATCH /api/equipment/EditEquipment

   // Equipment Usage Tracking
   GET  /api/equipmentusage/GetEquipmentUsageCalendar
   GET  /api/equipmentusage/GetEquipmentUsageHistory
   POST /api/equipmentusage/UseEquipment
   GET  /api/equipmentusage/GetEquipmentUsage
   ```

2. **Implemented Features**
   - Product catalog management with inventory tracking
   - Equipment management with location tracking
   - Equipment usage scheduling with conflict detection
   - Multi-tenant data isolation for all endpoints
   - GORM models for all core entities

---

## 🚀 **Deployment Instructions**

### **Deploy to GCP Cloud Functions**
```bash
# Deploy single function
gcloud functions deploy laboratory-api \
  --runtime go121 \
  --trigger-http \
  --source . \
  --entry-point HandleRequest \
  --set-env-vars DATABASE_URL=$NEON_URL,JWT_SECRET=$JWT_SECRET

# Deploy all functions (future)
./deploy.sh
```

### **Environment Variables**
```bash
# Local Development
DATABASE_URL=postgres://dev:dev@localhost:5433/neva_local?sslmode=disable
JWT_SECRET=your-local-jwt-secret-key
PORT=8080

# Production (Neon)
DATABASE_URL=postgres://username:password@ep-xxx.us-east-1.aws.neon.tech/dbname?sslmode=require
JWT_SECRET=your-production-secret-256-bits
```

---

## 📝 **Development Notes**

### **Testing Strategy**
- Local development with Docker Compose
- Integration testing with existing frontend
- Gradual endpoint migration (A/B testing)
- Performance monitoring

### **Migration Approach**
1. **Parallel Development**: New Go endpoints alongside existing C# API
2. **Gradual Cutover**: Switch frontend routes one by one
3. **Feature Parity**: Ensure identical functionality before switching
4. **Rollback Plan**: Keep C# API running until full migration

### **Key Decisions Made**
- **Go over C#**: Better serverless performance, lower costs
- **Chi over Gin**: Lightweight, stdlib-compatible
- **Neon over Cloud SQL**: Cost optimization for small scale
- **GCP over AWS**: Simpler deployment, better pricing for small scale

---

## 🎯 **Success Metrics**

### **Technical Goals**
- [ ] API response time < 200ms (95th percentile)
- [ ] 99.9% uptime
- [ ] Zero data loss during migration
- [ ] Full feature parity with C# API

### **Business Goals**
- [ ] 95%+ cost reduction
- [ ] Support 50+ concurrent users
- [ ] Maintain existing user experience
- [ ] Enable faster feature development

---

## 📞 **Next Session Priorities**

1. **Product Management Endpoints** (functions/products.go)
2. **Equipment Management Endpoints** (functions/equipment.go)
3. **Set up Neon Database** (production)
4. **Deploy Core Functions to GCP**
5. **Performance Testing & Optimization**

**Current Repository**: `/Users/luizcarvalho/labflux/NevaManagementApi/`
- **C# API**: `NevaManagement.Api/` (keep running)
- **Go Functions**: `neva-functions/` (active development)

**Status**: ✅ Phase 3 (Core API Migration) complete, ready for Phase 4 deployment
