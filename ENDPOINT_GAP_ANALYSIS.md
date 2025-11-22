# Frontend API Requirements vs Current Implementation - Gap Analysis

## Overview

This document compares the frontend's required endpoints (from `ENDPOINT_SPEC_FROM_FRONTEND.md`) against the current backend implementation.

---

## 1. PRODUCTS ⚠️ **Partial Implementation**

### Frontend Requirements (8 endpoints)
- `GET /api/products` - Paginated list
- `GET /api/products/:id` - Single product
- `GET /api/products/low-stock` - Low stock products
- `POST /api/products` - Create product
- `PUT /api/products/:id` - Update product
- `POST /api/products/:id/add-quantity` - Add quantity
- `POST /api/products/:id/use` - Record usage
- `DELETE /api/products/:id` - Delete product

### Current Implementation
| Frontend Endpoint | Current Endpoint | Status | Notes |
|-------------------|------------------|--------|-------|
| `GET /api/products` | `GET /api/products/GetAll` | ⚠️ **Incompatible URL** | Functionality exists, URL differs |
| `GET /api/products/:id` | `GET /api/products/GetProductById` | ⚠️ **Incompatible URL** | Uses query param instead of path param |
| `GET /api/products/low-stock` | `GET /api/products/GetLowInStockProducts` | ⚠️ **Incompatible URL** | Different URL structure |
| `POST /api/products` | `POST /api/products/Create` | ⚠️ **Incompatible URL** | |
| `PUT /api/products/:id` | `PATCH /api/products/EditProduct` | ⚠️ **Incompatible** | Wrong method + URL |
| `POST /api/products/:id/add-quantity` | `PATCH /api/products/AddQuantity` | ⚠️ **Incompatible** | Wrong method + uses query param |
| `POST /api/products/:id/use` | `PATCH /api/products/UseProduct` | ⚠️ **Incompatible** | Wrong method + uses query param |
| `DELETE /api/products/:id` | - | ❌ **MISSING** | Not implemented |

### Key Differences:
1. **URL Pattern**: Current uses action-based (`/GetAll`, `/Create`), Frontend expects RESTful (`/products`, `/products/:id`)
2. **HTTP Methods**: Current uses `PATCH`, Frontend expects `POST`/`PUT`
3. **Parameter Style**: Current uses query params (`?id=1`), Frontend expects path params (`/:id`)
4. **Response Structure**: May differ (needs verification)
5. **Missing DELETE operation**

---

## 2. EQUIPMENT ⚠️ **Partial Implementation**

### Frontend Requirements (9 endpoints)
- `GET /api/equipment` - Paginated list
- `GET /api/equipment/:id` - Single equipment
- `POST /api/equipment` - Create equipment
- `PUT /api/equipment/:id` - Update equipment
- `DELETE /api/equipment/:id` - Delete equipment
- `POST /api/equipment/:id/schedule` - Schedule usage
- `GET /api/equipment/:id/usage-history` - Usage history
- `GET /api/equipment/:id/calendar` - Calendar view
- `POST /api/equipment/:id/check-overlap` - Check conflicts

### Current Implementation
| Frontend Endpoint | Current Endpoint | Status | Notes |
|-------------------|------------------|--------|-------|
| `GET /api/equipment` | `GET /api/equipment/GetEquipments` | ⚠️ **Incompatible URL** | |
| `GET /api/equipment/:id` | `GET /api/equipment/GetDetailedEquipment` | ⚠️ **Incompatible** | Uses query param |
| `POST /api/equipment` | `POST /api/equipment/AddEquipment` | ⚠️ **Incompatible URL** | |
| `PUT /api/equipment/:id` | `PATCH /api/equipment/EditEquipment` | ⚠️ **Incompatible** | Wrong method |
| `DELETE /api/equipment/:id` | - | ❌ **MISSING** | |
| `POST /api/equipment/:id/schedule` | `POST /api/equipmentusage/UseEquipment` | ⚠️ **Incompatible** | Different URL structure |
| `GET /api/equipment/:id/usage-history` | `GET /api/equipmentusage/GetEquipmentUsageHistory` | ⚠️ **Incompatible** | Different URL |
| `GET /api/equipment/:id/calendar` | `GET /api/equipmentusage/GetEquipmentUsageCalendar` | ⚠️ **Incompatible** | Different URL |
| `POST /api/equipment/:id/check-overlap` | - | ❌ **MISSING** | May be embedded in UseEquipment logic |

### Key Differences:
1. **URL Pattern**: Same issues as Products (action-based vs RESTful)
2. **Nested Resources**: Equipment usage is separate (`/equipmentusage/*`), Frontend expects nested (`/equipment/:id/*`)
3. **Missing DELETE operation**
4. **Missing explicit check-overlap endpoint**

---

## 3. SAMPLES ❌ **NOT IMPLEMENTED**

### Frontend Requirements (5 endpoints)
- `GET /api/samples` - List all samples
- `GET /api/samples/:id` - Single sample
- `POST /api/samples` - Create sample
- `PUT /api/samples/:id` - Update sample
- `DELETE /api/samples/:id` - Delete sample

### Current Implementation
**❌ COMPLETELY MISSING** - No samples functionality implemented

### Required Actions:
1. Create `pkg/models/sample.go` model
2. Create `functions/samples.go` handler
3. Implement all CRUD operations
4. Add migration for `samples` table
5. Register routes in `internal/router.go`

---

## 4. REPLICAS ❌ **NOT IMPLEMENTED**

### Frontend Requirements (7 endpoints)
- `GET /api/replicas` - List all replicas
- `GET /api/replicas/sample/:sampleId` - Replicas for a sample
- `GET /api/replicas/:id` - Single replica
- `POST /api/replicas` - Create replica
- `PUT /api/replicas/:id` - Update replica
- `DELETE /api/replicas/:id` - Delete replica
- `POST /api/replicas/:id/subculture` - Record subculture

### Current Implementation
**❌ COMPLETELY MISSING** - No replicas functionality implemented

### Required Actions:
1. Create `pkg/models/replica.go` model
2. Create `functions/replicas.go` handler
3. Implement all CRUD + subculture logic
4. Add migration for `replicas` table
5. Register routes in `internal/router.go`

---

## 5. LOCATIONS ❌ **NOT IMPLEMENTED**

### Frontend Requirements (1 endpoint)
- `GET /api/locations` - List all locations

### Current Implementation
**❌ MISSING** - No dedicated locations endpoint

### Note:
- Location model may exist (referenced in products/equipment)
- Need endpoint to list all locations for a laboratory

### Required Actions:
1. Verify if `Location` model exists in `pkg/models/`
2. Create `functions/locations.go` handler
3. Implement GET endpoint
4. Register route in `internal/router.go`

---

## 6. RESEARCHERS ❌ **NOT IMPLEMENTED**

### Frontend Requirements (1 endpoint)
- `GET /api/researchers` - List all researchers

### Current Implementation
**❌ MISSING** - No dedicated researchers endpoint

### Note:
- We have authentication/user management
- Need endpoint to list users/researchers for a laboratory

### Required Actions:
1. Create endpoint in `functions/laboratory.go` or separate handler
2. Return users associated with a laboratory
3. Include role information

---

## 7. ROUTINES ❌ **NOT IMPLEMENTED**

### Frontend Requirements (13 endpoints)
- `GET /api/routines` - List routines
- `GET /api/routines/:id` - Single routine
- `GET /api/routines/upcoming` - Upcoming routine instances
- `POST /api/routines` - Create routine
- `PUT /api/routines/:id` - Update routine
- `DELETE /api/routines/:id` - Delete routine
- `POST /api/routines/:id/execute` - Start execution
- `GET /api/routines/executions/:executionId` - Execution details
- `PUT /api/routines/executions/:executionId/step` - Update step status
- `POST /api/routines/executions/:executionId/complete` - Complete execution
- `POST /api/routines/executions/:executionId/cancel` - Cancel execution
- `GET /api/routines/:id/executions` - Execution history

### Current Implementation
**❌ COMPLETELY MISSING** - No routines functionality implemented

### Required Actions:
1. Create `pkg/models/routine.go`, `routine_execution.go`, `routine_step.go` models
2. Create `functions/routines.go` handler
3. Implement complex scheduling logic (recurring, one-time)
4. Implement execution tracking
5. Add material deduction on completion
6. Add migrations for all routine-related tables
7. Register routes in `internal/router.go`

---

## Summary

### Implementation Status

| Domain | Endpoints Needed | Endpoints Exist | Status | Completion % |
|--------|-----------------|-----------------|--------|--------------|
| Products | 8 | 7 | ⚠️ Partial | ~70% |
| Equipment | 9 | 6 | ⚠️ Partial | ~60% |
| Samples | 5 | 0 | ❌ Missing | 0% |
| Replicas | 7 | 0 | ❌ Missing | 0% |
| Locations | 1 | 0 | ❌ Missing | 0% |
| Researchers | 1 | 0 | ❌ Missing | 0% |
| Routines | 13 | 0 | ❌ Missing | 0% |
| **TOTAL** | **44** | **13** | - | **~30%** |

### Critical Issues

#### 1. **URL Structure Incompatibility** (BREAKING)
- **Current**: Action-based URLs (`/api/products/GetAll`, `/api/products/Create`)
- **Required**: RESTful URLs (`/api/products`, `/api/products/:id`)
- **Impact**: Frontend cannot call any existing endpoints without changes

#### 2. **HTTP Method Mismatch** (BREAKING)
- **Current**: Uses `PATCH` for updates
- **Required**: Uses `PUT` and `POST`
- **Impact**: Frontend requests will fail with 405 Method Not Allowed

#### 3. **Parameter Style Mismatch** (BREAKING)
- **Current**: Uses query parameters (`?id=1&laboratoryId=1`)
- **Required**: Uses path parameters (`/api/products/:id?laboratory_id=1`)
- **Impact**: Routing won't match frontend requests

#### 4. **Missing Domains** (BLOCKER)
- Samples (5 endpoints)
- Replicas (7 endpoints)
- Locations (1 endpoint)
- Researchers (1 endpoint)
- Routines (13 endpoints)
- **Total**: 27 endpoints (61% of requirements)

#### 5. **Missing DELETE Operations** (HIGH)
- Products: Cannot delete products
- Equipment: Cannot delete equipment
- All other domains: N/A (not implemented)

---

## Recommended Action Plan

### Phase 1: Fix Existing Endpoints (Products & Equipment)
**Goal**: Make current functionality compatible with frontend

1. **Refactor URL Structure**
   - Change from action-based to RESTful paths
   - Use path parameters instead of query parameters
   - Update Chi router configuration

2. **Fix HTTP Methods**
   - Change `PATCH /EditProduct` → `PUT /products/:id`
   - Change `PATCH /AddQuantity` → `POST /products/:id/add-quantity`
   - Ensure proper method usage

3. **Add DELETE Endpoints**
   - `DELETE /api/products/:id`
   - `DELETE /api/equipment/:id`
   - Implement cascade delete logic

4. **Verify Response Structures**
   - Compare actual responses with frontend spec
   - Add missing fields (e.g., `location` object nesting)
   - Ensure consistent pagination format

**Estimated Effort**: 2-3 days

### Phase 2: Implement Locations & Researchers
**Goal**: Add simple list endpoints

1. **Locations**
   - Create handler
   - Implement `GET /api/locations`

2. **Researchers**
   - Create handler
   - Implement `GET /api/researchers`

**Estimated Effort**: 1 day

### Phase 3: Implement Samples & Replicas
**Goal**: Add biological sample management

1. **Samples**
   - Create models and migrations
   - Implement all CRUD endpoints
   - Add replica creation on sample creation

2. **Replicas**
   - Create models and migrations
   - Implement all CRUD endpoints
   - Implement subculture tracking

**Estimated Effort**: 3-4 days

### Phase 4: Implement Routines System
**Goal**: Add laboratory routine management

1. **Core Routines**
   - Create models for routines, executions, steps
   - Implement CRUD endpoints
   - Implement scheduling logic (recurring, one-time)

2. **Execution System**
   - Implement execution tracking
   - Implement step completion
   - Implement material deduction

**Estimated Effort**: 5-7 days

---

## Breaking Changes Summary

⚠️ **The current API is NOT compatible with the frontend requirements**

### What Breaks:
1. All URL paths are different
2. All HTTP methods for mutations are different
3. Parameter passing mechanism is different
4. 61% of required functionality doesn't exist

### Migration Strategy:
**Option A: Big Bang (Recommended for Greenfield)**
- Refactor all at once
- Deploy with breaking changes
- Frontend updates simultaneously

**Option B: Dual Support (Complex)**
- Support both old and new URLs temporarily
- Deprecation warnings
- Gradual migration
- Higher maintenance burden

**Option C: Version Prefix**
- Keep current as `/v1/products/GetAll`
- New endpoints as `/v2/products` or `/api/products`
- Clean separation
- Eventually deprecate v1

---

## Notes

1. **Response Structure**: Detailed response comparison needed for each endpoint
2. **Authentication**: Frontend spec mentions JWT but doesn't detail auth endpoints - current implementation seems adequate
3. **Error Handling**: Frontend spec defines standard error format - verify current implementation matches
4. **Multi-Tenancy**: Current implementation has strict `laboratory_id` filtering - ensure all new endpoints maintain this
5. **Pagination**: Verify pagination format matches frontend expectations (`nextPage`, `totalCount`)

---

**Generated**: 2025-11-16
**Status**: Current as of latest codebase exploration
