# LabFlux API - Implementation TODO List

## Overview
This document tracks missing features and incomplete implementations identified during codebase analysis. The project is approximately 70% complete with core functionality implemented.

---

## Critical Missing Features (HIGH PRIORITY)

### 1. Location Management Endpoints
**Status:** Model exists, endpoints missing
**Priority:** HIGH - Required by Product and Equipment features
**File to create:** `functions/locations.go`

**Endpoints needed:**
- `GET /api/locations` - List all locations for a laboratory
- `POST /api/locations` - Create new location
- `GET /api/locations/:id` - Get location details
- `PATCH /api/locations/:id` - Edit location
- `DELETE /api/locations/:id` - Delete location
- `GET /api/locations/:id/sublocations` - Get sub-locations (hierarchical support)

**Implementation notes:**
- Location model already exists in `pkg/models/laboratory.go`
- Supports hierarchical locations via `SubLocationID`
- Must filter by `LaboratoryID` for multi-tenant isolation
- Register routes in `internal/router.go`

---

### 2. Researcher Management Endpoints
**Status:** Model exists, endpoints missing
**Priority:** HIGH - Required by Equipment Usage feature
**File to create:** `functions/researchers.go`

**Endpoints needed:**
- `GET /api/researchers` - List all researchers for a laboratory
- `POST /api/researchers` - Create new researcher
- `GET /api/researchers/:id` - Get researcher details
- `PATCH /api/researchers/:id` - Edit researcher
- `DELETE /api/researchers/:id` - Delete researcher

**Implementation notes:**
- Researcher model already exists in `pkg/models/laboratory.go`
- Links to Clerk user via `ClerkUserID`
- Must filter by `LaboratoryID` for multi-tenant isolation
- Currently referenced in Equipment Usage but can't be managed
- Register routes in `internal/router.go`

---

## Model Enhancements (MEDIUM-HIGH PRIORITY)

### 3. Enhance Product Model
**Status:** Partial implementation - missing fields from Postman collection
**Priority:** MEDIUM-HIGH
**File to modify:** `pkg/models/laboratory.go`

**Missing fields to add:**
```go
MinimumQuantity  float64 `gorm:"type:decimal(10,2)" json:"minimum_quantity"`
Manufacturer     string  `gorm:"type:varchar(255)" json:"manufacturer"`
CatalogNumber    string  `gorm:"type:varchar(100)" json:"catalog_number"`
```

**Affected files:**
- `pkg/models/laboratory.go` - Add fields to Product struct
- `functions/products.go` - Update CreateProduct, EditProduct requests/responses
- Database migration will auto-update via GORM

**Benefits:**
- Better inventory management with minimum quantity alerts
- Improved product tracking with manufacturer/catalog info
- Matches Postman collection expectations

---

### 4. Enhance Equipment Model
**Status:** Partial implementation - missing fields from Postman collection
**Priority:** MEDIUM-HIGH
**File to modify:** `pkg/models/laboratory.go`

**Missing fields to add:**
```go
Manufacturer     string    `gorm:"type:varchar(255)" json:"manufacturer"`
Model            string    `gorm:"type:varchar(255)" json:"model"`
SerialNumber     string    `gorm:"type:varchar(100);unique" json:"serial_number"`
AcquisitionDate  time.Time `json:"acquisition_date"`
Status           string    `gorm:"type:varchar(50);default:'available'" json:"status"` // available, in_use, maintenance, broken, retired
```

**Affected files:**
- `pkg/models/laboratory.go` - Add fields to Equipment struct
- `functions/equipment.go` - Update AddEquipment, EditEquipment requests/responses
- Database migration will auto-update via GORM

**Status enum values:**
- `available` - Ready for use
- `in_use` - Currently being used
- `maintenance` - Under maintenance
- `broken` - Needs repair
- `retired` - No longer in use

---

### 5. Enhance Equipment Usage Model
**Status:** Partial implementation - missing fields from Postman collection
**Priority:** MEDIUM
**File to modify:** `pkg/models/laboratory.go`

**Missing fields to add:**
```go
Purpose  string `gorm:"type:text" json:"purpose"`
Notes    string `gorm:"type:text" json:"notes"`
```

**Affected files:**
- `pkg/models/laboratory.go` - Add fields to EquipmentUsage struct
- `functions/equipment_usage.go` - Update UseEquipment request/response
- Database migration will auto-update via GORM

**Benefits:**
- Better documentation of equipment usage
- Track experiment details and observations
- Matches Postman collection expectations

---

## Feature Completion (MEDIUM PRIORITY)

### 6. Implement Product Usage Analytics
**Status:** TODO comment in code
**Priority:** MEDIUM
**Location:** `functions/products.go:170`

**Current code:**
```go
QuantityUsedInTheLastThreeMonths: 0, // TODO: Calculate from usage data
```

**Implementation approach:**
1. Create ProductUsage model (see #7)
2. Query usage records from last 3 months
3. Sum quantities used
4. Update `GetDetailedProductByID` to return calculated value

**SQL query needed:**
```sql
SELECT SUM(quantity) FROM product_usages
WHERE product_id = ?
  AND laboratory_id = ?
  AND created_at >= NOW() - INTERVAL '3 months'
```

---

### 7. Create ProductUsage Model and Tracking
**Status:** Not implemented
**Priority:** MEDIUM
**Files to create/modify:** `pkg/models/laboratory.go`, `functions/product_usage.go`

**New model needed:**
```go
type ProductUsage struct {
    gorm.Model
    ProductID    uint      `gorm:"not null" json:"product_id"`
    Product      Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    Quantity     float64   `gorm:"type:decimal(10,2);not null" json:"quantity"`
    Unit         string    `gorm:"type:varchar(50)" json:"unit"`
    Purpose      string    `gorm:"type:text" json:"purpose"`
    Notes        string    `gorm:"type:text" json:"notes"`
    UsedBy       string    `gorm:"type:varchar(255)" json:"used_by"` // User/Researcher name
    LaboratoryID uint      `gorm:"not null;index" json:"laboratory_id"`
    Laboratory   Laboratory `gorm:"foreignKey:LaboratoryID" json:"laboratory,omitempty"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Endpoints needed:**
- `GET /api/products/:id/usage-history` - Get usage history for a product
- `GET /api/products/usage-analytics` - Get analytics/dashboard data
- Modify `UseProduct` endpoint to create ProductUsage records

**Benefits:**
- Track who used what, when, and why
- Enable analytics and reporting
- Support audit trails for laboratory compliance

---

## Integration & Testing (ONGOING)

### 8. Update All Affected Endpoints
**Status:** Pending model enhancements
**Priority:** HIGH (after model changes)

**Files to update:**
- `functions/products.go` - Update all Product-related endpoints
- `functions/equipment.go` - Update all Equipment-related endpoints
- `functions/equipment_usage.go` - Update Equipment Usage endpoints
- Ensure all request/response structs match enhanced models
- Update validation rules as needed

**Testing checklist:**
- [ ] Verify all CRUD operations work with new fields
- [ ] Test multi-tenant isolation still works
- [ ] Verify Clerk authentication on protected routes
- [ ] Test pagination on list endpoints
- [ ] Verify soft deletes work correctly

---

### 9. Local Testing
**Status:** Pending implementation
**Priority:** HIGH (before deployment)

**Testing steps:**
```bash
# Start local environment
docker compose up -d postgres
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"
go run cmd/main.go

# Test with curl or Postman collection
curl http://localhost:8080/health
curl http://localhost:8080/api/ping

# Run through all endpoints in Postman collection
```

**Test scenarios:**
- [ ] Create laboratory and invite users
- [ ] Create locations (including sub-locations)
- [ ] Create researchers
- [ ] Create products with new fields
- [ ] Create equipment with new fields
- [ ] Schedule equipment usage with purpose/notes
- [ ] Test product usage tracking and analytics
- [ ] Verify low stock alerts use minimum quantity
- [ ] Test all GET endpoints with pagination
- [ ] Verify multi-tenant data isolation

---

### 10. Update Postman Collection
**Status:** Needs sync with implementation
**Priority:** MEDIUM
**File to update:** `docs/LabFlux_API.postman_collection.json`

**Updates needed:**
- Add Location management endpoints
- Add Researcher management endpoints
- Update Product endpoints with new fields (minimumQuantity, manufacturer, catalogNumber)
- Update Equipment endpoints with new fields (manufacturer, model, serialNumber, etc.)
- Update Equipment Usage with purpose/notes fields
- Add Product Usage History endpoints
- Ensure all examples match current implementation

---

## Future Features (LOW PRIORITY)

### Container Management
**Status:** Not started - marked as "future" in migration plan
**Priority:** LOW
**Mentioned in:** `docs/MIGRATION_PLAN.md:43`

**Design considerations:**
- Containers can hold products (freezers, refrigerators, storage boxes)
- Hierarchical structure (container -> shelf -> box -> product)
- Temperature tracking for climate-controlled containers
- Capacity management

**Potential model:**
```go
type Container struct {
    gorm.Model
    Name             string  `gorm:"type:varchar(255);not null" json:"name"`
    Description      string  `gorm:"type:text" json:"description"`
    ContainerType    string  `gorm:"type:varchar(100)" json:"container_type"` // freezer, refrigerator, cabinet, etc.
    LocationID       uint    `json:"location_id"`
    Location         Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
    Temperature      float64 `gorm:"type:decimal(5,2)" json:"temperature"` // Current temperature
    ParentContainerID *uint  `json:"parent_container_id"` // For nested containers
    LaboratoryID     uint    `gorm:"not null;index" json:"laboratory_id"`
    Laboratory       Laboratory `gorm:"foreignKey:LaboratoryID" json:"laboratory,omitempty"`
    DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

**Endpoints needed:**
- `GET /api/containers` - List containers
- `POST /api/containers` - Create container
- `GET /api/containers/:id` - Get container details
- `PATCH /api/containers/:id` - Edit container
- `DELETE /api/containers/:id` - Delete container
- `GET /api/containers/:id/contents` - Get products in container

---

## Implementation Roadmap

### Phase 1: Critical Dependencies (Week 1)
1. Implement Location Management endpoints
2. Implement Researcher Management endpoints
3. Test basic CRUD operations

### Phase 2: Model Enhancements (Week 1-2)
4. Enhance Product model with missing fields
5. Enhance Equipment model with missing fields
6. Enhance Equipment Usage model with Purpose/Notes
7. Update all affected endpoints

### Phase 3: Analytics & Tracking (Week 2)
8. Create ProductUsage model
9. Implement usage tracking endpoints
10. Implement 3-month usage analytics

### Phase 4: Testing & Documentation (Week 2-3)
11. Comprehensive local testing
12. Update Postman collection
13. Update documentation

### Phase 5: Deployment (Week 3)
14. Deploy to GCP Cloud Functions via Pulumi
15. Configure production database (Neon)
16. Production testing

---

## Current Status Summary

**Completion:** ~70% of core functionality

### ✅ Completed
- Laboratory management with invitations
- Clerk authentication integration
- Product management (CRUD + inventory)
- Equipment management (CRUD)
- Equipment usage scheduling
- Multi-tenant data isolation
- Local development environment

### 🔄 In Progress
- Location management (model exists, endpoints needed)
- Researcher management (model exists, endpoints needed)
- Model enhancements (missing fields)
- Usage analytics

### ⏳ Not Started
- Product usage history tracking
- Container management (future)

---

## Notes

- All new endpoints must include multi-tenant isolation via `LaboratoryID`
- All protected routes must use `AuthMiddleware` from `pkg/middleware/auth.go`
- Use consistent error handling patterns from existing endpoints
- Follow naming conventions from existing code (PascalCase for endpoint paths)
- Ensure GORM auto-migration works by updating models in `pkg/models/laboratory.go`
- Update router registration in `internal/router.go` for all new route groups

---

**Last Updated:** 2025-01-15
**Next Review:** After Phase 1 completion
