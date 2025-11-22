# LabFlux API - Available Endpoints

**Base URL (Production):** `https://us-central1-labflux-475521.cloudfunctions.net/labflux-api`
**Base URL (Local):** `http://localhost:8080`

---

## Authentication

All protected endpoints require a Clerk authentication token in the `Authorization` header:
```
Authorization: Bearer <clerk_token>
```

---

## Public Endpoints (No Authentication Required)

### Authentication

#### Sign Up (Register)
```http
POST /api/auth/signup
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "clerk_user_id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "laboratory_id": null,
    "role": "user",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Description:**
- Creates a new user account using Clerk Backend API
- Password must be at least 8 characters
- Returns JWT token for immediate use
- User can accept laboratory invitations after signup

**Validations:**
- Email must be valid and unique
- Password must meet Clerk's security requirements
- Clerk will handle email verification

#### Sign In (Login)
```http
POST /api/auth/signin
POST /api/auth/login
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "clerk_user_id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "laboratory_id": 5,
    "role": "researcher",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Description:**
- Authenticates user with Clerk Backend API
- Returns JWT token for API access
- Token should be included in subsequent requests

**Error Responses:**
- `401 Unauthorized` - Invalid email or password
- `500 Internal Server Error` - Server configuration error

**Note:** Both `/signin` and `/login` work identically (alias).

---

### Health & Status

#### Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "ok"
}
```

#### Ping API
```http
GET /api/ping
```
**Response:**
```json
{
  "message": "pong"
}
```

---

### Webhooks

#### Clerk Webhook Handler
```http
POST /api/webhooks/clerk
```
**Headers:**
```
Content-Type: application/json
svix-id: <webhook-id>
svix-timestamp: <timestamp>
svix-signature: <signature>
```
**Description:** Handles Clerk user lifecycle events (user.created, user.updated, user.deleted)

---

### Invitations

#### Accept Invitation
```http
POST /api/auth/invite/:token
```
**Parameters:**
- `token` (path) - Invitation token from email

**Request Body:**
```json
{
  "clerkUserId": "user_xxx"
}
```

**Response:**
```json
{
  "message": "Invitation accepted successfully",
  "laboratory": {
    "id": 1,
    "name": "Research Lab Alpha"
  }
}
```

---

## Protected Endpoints (Authentication Required)

### Authentication

#### Get Current User
```http
GET /api/auth/me
```
**Response:**
```json
{
  "id": 1,
  "clerkUserId": "user_xxx",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "laboratoryId": 1,
  "role": "admin"
}
```

---

### Laboratory Management

#### Create Laboratory
```http
POST /api/laboratories
```
**Request Body:**
```json
{
  "name": "Research Lab Alpha",
  "description": "Main research laboratory for biology experiments",
  "address": "123 Science Ave, Research City"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Research Lab Alpha",
  "description": "Main research laboratory for biology experiments",
  "address": "123 Science Ave, Research City",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

#### Get Laboratory
```http
GET /api/laboratories/:laboratoryId
```
**Parameters:**
- `laboratoryId` (path) - Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "Research Lab Alpha",
  "description": "Main research laboratory for biology experiments",
  "address": "123 Science Ave, Research City",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

---

### Invitation Management

#### Create Invitation
```http
POST /api/laboratories/:laboratoryId/invitations
```
**Parameters:**
- `laboratoryId` (path) - Laboratory ID

**Request Body:**
```json
{
  "email": "newuser@example.com",
  "role": "researcher"
}
```

**Response:**
```json
{
  "id": 1,
  "email": "newuser@example.com",
  "invitationToken": "abc123xyz",
  "role": "researcher",
  "expiresAt": "2024-01-22T10:30:00Z",
  "isAccepted": false
}
```

**Notes:**
- Invitations expire after 7 days
- Send the invitation URL to user: `https://your-frontend.com/accept-invitation?token={invitationToken}`

#### Get Invitations
```http
GET /api/laboratories/:laboratoryId/invitations
```
**Parameters:**
- `laboratoryId` (path) - Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "email": "newuser@example.com",
    "invitationToken": "abc123xyz",
    "role": "researcher",
    "expiresAt": "2024-01-22T10:30:00Z",
    "isAccepted": false,
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

---

### Product Management

#### Get All Products
```http
GET /api/products/GetAll?laboratoryId={id}&page={page}
```
**Query Parameters:**
- `laboratoryId` (required) - Laboratory ID
- `page` (optional) - Page number (default: 1, 10 items per page)

**Response:**
```json
[
  {
    "id": 1,
    "name": "Ethanol 99.9%",
    "description": "High purity ethanol for laboratory use",
    "quantity": 10.5,
    "formula": "C2H5OH",
    "unit": "L",
    "expirationDate": "2025-12-31T00:00:00Z",
    "locationId": 1
  }
]
```

#### Get Product By ID
```http
GET /api/products/GetProductById?id={id}&laboratoryId={labId}
```
**Query Parameters:**
- `id` (required) - Product ID
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "Ethanol 99.9%",
  "description": "High purity ethanol for laboratory use",
  "quantity": 10.5,
  "formula": "C2H5OH",
  "unit": "L",
  "expirationDate": "2025-12-31T00:00:00Z",
  "locationId": 1
}
```

#### Get Detailed Product By ID
```http
GET /api/products/GetDetailedProductById?id={id}&laboratoryId={labId}
```
**Query Parameters:**
- `id` (required) - Product ID
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "Ethanol 99.9%",
  "description": "High purity ethanol for laboratory use",
  "quantity": 10.5,
  "quantity_used_in_the_last_three_months": 0,
  "formula": "C2H5OH",
  "unit": "L",
  "expirationDate": "2025-12-31T00:00:00Z",
  "location": {
    "id": 1,
    "name": "Shelf A3",
    "description": "Chemical storage shelf",
    "sub_location_id": null
  }
}
```

**Note:** `quantity_used_in_the_last_three_months` currently returns 0 (see TODO.md)

#### Get Low In Stock Products
```http
GET /api/products/GetLowInStockProducts
```
**Description:** Returns products with quantity <= 10

**Response:**
```json
[
  {
    "id": 1,
    "name": "Ethanol 99.9%",
    "description": "High purity ethanol for laboratory use",
    "quantity": 5.0,
    "formula": "C2H5OH",
    "unit": "L",
    "expirationDate": "2025-12-31T00:00:00Z",
    "locationId": 1
  }
]
```

#### Create Product
```http
POST /api/products/Create
```
**Request Body:**
```json
{
  "name": "Ethanol 99.9%",
  "description": "High purity ethanol for laboratory use",
  "location_id": 1,
  "quantity": 10.5,
  "formula": "C2H5OH",
  "unit": "L",
  "expiration_date": "2025-12-31"
}
```

**Response:**
```
Successfully created Ethanol 99.9%.
```

**Notes:**
- `location_id` is required
- `expiration_date` can be "2025-12-31" or "2025-12-31T00:00:00Z" format
- Laboratory ID is extracted from authentication token

#### Edit Product
```http
PATCH /api/products/EditProduct
```
**Request Body:**
```json
{
  "id": 1,
  "name": "Ethanol 99.9% - Updated",
  "description": "High purity ethanol for laboratory use - Updated",
  "location_id": 1,
  "quantity": 15.0,
  "formula": "C2H5OH",
  "unit": "L",
  "expiration_date": "2025-12-31"
}
```

**Response:**
```
Successfully edited Ethanol 99.9% - Updated.
```

#### Add Quantity to Product
```http
PATCH /api/products/AddQuantity
```
**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 5.0
}
```

**Response:**
```
Successfully added 5.00 to product.
```

**Description:** Increases product quantity (e.g., restocking)

#### Use Product
```http
PATCH /api/products/UseProduct
```
**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2.5,
  "unit": "L"
}
```

**Response:**
```
Successfully used 2.50L.
```

**Description:** Decreases product quantity (e.g., consumption)
**Validation:** Returns 400 error if insufficient quantity available

---

### Equipment Management

#### Get All Equipment
```http
GET /api/equipment/GetEquipments?laboratoryId={id}
```
**Query Parameters:**
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "Centrifuge Model X1000",
    "description": "High-speed centrifuge for sample preparation",
    "property_number": "EQ-2024-001",
    "location_id": 1
  }
]
```

#### Get Detailed Equipment
```http
GET /api/equipment/GetDetailedEquipment?id={id}&laboratoryId={labId}
```
**Query Parameters:**
- `id` (required) - Equipment ID
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "Centrifuge Model X1000",
  "description": "High-speed centrifuge for sample preparation",
  "property_number": "EQ-2024-001",
  "location": {
    "id": 1,
    "name": "Lab Room 201",
    "description": "Main laboratory room",
    "sub_location_id": null
  },
  "usage_list": [
    {
      "id": 1,
      "researcher_id": 1,
      "researcher_name": "Dr. Jane Smith",
      "start_date": "2024-01-15T09:00:00Z",
      "end_date": "2024-01-15T12:00:00Z",
      "description": "Sample centrifugation"
    }
  ]
}
```

#### Add Equipment
```http
POST /api/equipment/AddEquipment
```
**Request Body:**
```json
{
  "name": "Centrifuge Model X1000",
  "description": "High-speed centrifuge for sample preparation",
  "property_number": "EQ-2024-001",
  "location_id": 1
}
```

**Response:**
```
Centrifuge Model X1000 was created successfully.
```

#### Edit Equipment
```http
PATCH /api/equipment/EditEquipment
```
**Request Body:**
```json
{
  "id": 1,
  "name": "Centrifuge Model X1000 - Updated",
  "description": "High-speed centrifuge - maintenance completed",
  "property_number": "EQ-2024-001",
  "location_id": 2
}
```

**Response:**
```
Centrifuge Model X1000 - Updated was updated successfully.
```

---

### Equipment Usage Management

#### Get Equipment Usage Calendar
```http
GET /api/equipmentusage/GetEquipmentUsageCalendar?id={equipmentId}&laboratoryId={labId}
```
**Query Parameters:**
- `id` (required) - Equipment ID
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "equipment_id": 1,
    "researcher_id": 1,
    "researcher_name": "Dr. Jane Smith",
    "description": "Sample centrifugation",
    "start_date": "2024-01-15T09:00:00Z",
    "end_date": "2024-01-15T12:00:00Z"
  }
]
```

**Description:** Returns all scheduled usage for an equipment (past and future)

#### Get Equipment Usage History
```http
GET /api/equipmentusage/GetEquipmentUsageHistory?id={equipmentId}&laboratoryId={labId}
```
**Query Parameters:**
- `id` (required) - Equipment ID
- `laboratoryId` (required) - Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "equipment_id": 1,
    "researcher_id": 1,
    "researcher_name": "Dr. Jane Smith",
    "description": "Sample centrifugation",
    "start_date": "2024-01-10T09:00:00Z",
    "end_date": "2024-01-10T12:00:00Z"
  }
]
```

**Description:** Returns only past usage (end_date < now), ordered by end_date desc

#### Get Equipment Usage (Paginated)
```http
GET /api/equipmentusage/GetEquipmentUsage?equipmentId={id}&laboratoryId={labId}&page={page}
```
**Query Parameters:**
- `equipmentId` (required) - Equipment ID
- `laboratoryId` (required) - Laboratory ID
- `page` (optional) - Page number (default: 1, 10 items per page)

**Response:**
```json
[
  {
    "id": 1,
    "equipment_id": 1,
    "researcher_id": 1,
    "researcher_name": "Dr. Jane Smith",
    "description": "Sample centrifugation",
    "start_date": "2024-01-15T09:00:00Z",
    "end_date": "2024-01-15T12:00:00Z"
  }
]
```

**Description:** Returns paginated usage records, ordered by start_date desc

#### Schedule/Record Equipment Usage
```http
POST /api/equipmentusage/UseEquipment
```
**Request Body:**
```json
{
  "equipment_id": 1,
  "researcher_id": 1,
  "description": "Sample centrifugation for experiment ABC",
  "start_date": "2024-01-16T09:00:00Z",
  "end_date": "2024-01-16T12:00:00Z"
}
```

**Response:**
```
Equipment was used successfully.
```

**Validations:**
- End date must be after start date
- Equipment must exist and belong to the laboratory
- Researcher must exist and belong to the laboratory
- Returns 409 Conflict if equipment is already scheduled during that time period

**Date formats accepted:**
- `2024-01-16T09:00:00Z`
- `2024-01-16T09:00:00-03:00`

---

## Error Responses

All endpoints return standard HTTP status codes:

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```
Or plain text error message

### 401 Unauthorized
```json
{
  "error": "Laboratory ID not found in token"
}
```

### 404 Not Found
```json
{
  "error": "Product not found"
}
```

### 409 Conflict
```json
{
  "error": "Equipment is already in use during the specified time period"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to fetch products"
}
```

---

## Multi-Tenant Architecture

All protected endpoints automatically filter data by the user's `LaboratoryID` extracted from the Clerk authentication token. This ensures complete data isolation between laboratories.

**Important:**
- Users can only access data from their assigned laboratory
- Some endpoints require `laboratoryId` query parameter for additional validation
- Laboratory ID in token must match the requested laboratory ID

---

## Data Models

### Laboratory
- `id` - uint
- `name` - string
- `description` - string
- `address` - string
- `createdAt` - timestamp
- `updatedAt` - timestamp

### User
- `id` - uint
- `clerkUserId` - string (unique)
- `email` - string (unique)
- `firstName` - string
- `lastName` - string
- `laboratoryId` - uint (nullable)
- `role` - string

### Product
- `id` - uint
- `name` - string
- `description` - string
- `locationId` - uint
- `quantity` - float64
- `formula` - string
- `unit` - string
- `expirationDate` - timestamp
- `laboratoryId` - uint

### Equipment
- `id` - uint
- `name` - string
- `description` - string
- `propertyNumber` - string
- `locationId` - uint
- `laboratoryId` - uint

### Location
- `id` - uint
- `name` - string
- `description` - string
- `subLocationId` - uint (nullable, for hierarchical locations)
- `laboratoryId` - uint

### Researcher
- `id` - uint
- `name` - string
- `clerkUserId` - string (unique)
- `email` - string (unique)
- `laboratoryId` - uint

### EquipmentUsage
- `id` - uint
- `equipmentId` - uint
- `researcherId` - uint
- `description` - string
- `startDate` - timestamp
- `endDate` - timestamp
- `laboratoryId` - uint

---

## Notes for Frontend Integration

1. **Authentication Flow:**
   - Use Clerk for user authentication
   - Include Clerk token in all API requests: `Authorization: Bearer {token}`
   - Token contains user's laboratory ID for multi-tenant isolation

2. **Pagination:**
   - Product list supports pagination (10 items per page)
   - Equipment usage supports pagination (10 items per page)
   - Use `page` query parameter (starts at 1)

3. **Date Handling:**
   - API accepts ISO 8601 format with timezone
   - API returns dates in `2006-01-02T15:04:05Z07:00` format
   - For simple dates, `2006-01-02` format is also accepted

4. **Validation:**
   - Required fields will return 400 Bad Request if missing
   - Business logic validation (e.g., insufficient stock) returns 400
   - Conflicts (e.g., equipment double-booking) return 409

5. **Missing Features (see TODO.md):**
   - Location management endpoints not yet available (need to create locations manually in DB)
   - Researcher management endpoints not yet available (need to create researchers manually in DB)
   - Product/Equipment models missing some fields shown in Postman collection

---

**Last Updated:** 2025-01-15
**API Version:** 1.0 (Development)
