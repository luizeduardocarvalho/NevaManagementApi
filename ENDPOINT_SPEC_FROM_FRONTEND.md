# LabFlux API Endpoints Specification

This document specifies all the API endpoints needed for the LabFlux React frontend integration.

---

## 1. Products

### GET `/api/products`
**Description:** Get paginated list of products for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID
- `page` (number, optional, default: 1): Page number
- `pageSize` (number, optional, default: 9): Items per page

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Ethanol 99.5%",
      "description": "High-purity ethanol for laboratory use",
      "formula": "C2H5OH",
      "quantity": 5.5,
      "unit": "L",
      "location": {
        "id": 1,
        "name": "Storage Room A",
        "description": "Main storage area"
      },
      "expiration_date": "2025-12-31T00:00:00.000Z",
      "quantity_used_in_the_last_three_months": 2.5,
      "laboratory_id": 1
    }
  ],
  "nextPage": 2,
  "totalCount": 25
}
```

---

### GET `/api/products/:id`
**Description:** Get detailed product information

**Path Parameters:**
- `id` (number, required): Product ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "Ethanol 99.5%",
  "description": "High-purity ethanol for laboratory use",
  "formula": "C2H5OH",
  "quantity": 5.5,
  "unit": "L",
  "location": {
    "id": 1,
    "name": "Storage Room A",
    "description": "Main storage area"
  },
  "location_id": 1,
  "expiration_date": "2025-12-31T00:00:00.000Z",
  "quantity_used_in_the_last_three_months": 2.5,
  "laboratory_id": 1
}
```

---

### GET `/api/products/low-stock`
**Description:** Get products with low stock (quantity < 3-month usage average)

**Response:**
```json
[
  {
    "id": 1,
    "name": "Ethanol 99.5%",
    "description": "High-purity ethanol for laboratory use",
    "formula": "C2H5OH",
    "quantity": 1.0,
    "unit": "L",
    "location": {
      "id": 1,
      "name": "Storage Room A",
      "description": "Main storage area"
    },
    "expiration_date": "2025-12-31T00:00:00.000Z",
    "quantity_used_in_the_last_three_months": 2.5,
    "laboratory_id": 1
  }
]
```

---

### POST `/api/products`
**Description:** Create a new product

**Request Body:**
```json
{
  "name": "Ethanol 99.5%",
  "description": "High-purity ethanol for laboratory use",
  "formula": "C2H5OH",
  "quantity": 5.5,
  "unit": "L",
  "location_id": 1,
  "expiration_date": "2025-12-31T00:00:00.000Z",
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Product created successfully"
}
```

---

### PUT `/api/products/:id`
**Description:** Update an existing product

**Path Parameters:**
- `id` (number, required): Product ID

**Request Body:**
```json
{
  "id": 1,
  "name": "Ethanol 99.5%",
  "description": "High-purity ethanol for laboratory use",
  "formula": "C2H5OH",
  "quantity": 5.5,
  "unit": "L",
  "location_id": 1,
  "expiration_date": "2025-12-31T00:00:00.000Z",
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Product updated successfully"
}
```

---

### POST `/api/products/:id/add-quantity`
**Description:** Add quantity to an existing product

**Path Parameters:**
- `id` (number, required): Product ID

**Request Body:**
```json
{
  "quantity": 2.5
}
```

**Response:**
```json
{
  "message": "Quantity added successfully"
}
```

---

### POST `/api/products/:id/use`
**Description:** Record product usage (deducts from quantity)

**Path Parameters:**
- `id` (number, required): Product ID

**Request Body:**
```json
{
  "quantity": 1.5,
  "unit": "L"
}
```

**Response:**
```json
{
  "message": "Product usage recorded successfully"
}
```

---

### DELETE `/api/products/:id`
**Description:** Delete a product

**Path Parameters:**
- `id` (number, required): Product ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "message": "Product deleted successfully"
}
```

---

## 2. Equipment

### GET `/api/equipment`
**Description:** Get paginated list of equipment for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID
- `page` (number, optional, default: 1): Page number
- `pageSize` (number, optional, default: 9): Items per page

**Response:**
```json
{
  "equipment": [
    {
      "id": 1,
      "name": "PCR Thermocycler"
    }
  ],
  "nextPage": 2,
  "totalCount": 15
}
```

---

### GET `/api/equipment/:id`
**Description:** Get detailed equipment information

**Path Parameters:**
- `id` (number, required): Equipment ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "PCR Thermocycler",
  "description": "Bio-Rad T100 Thermal Cycler for PCR amplification",
  "propertyNumber": "EQ-2023-001",
  "location": {
    "id": 1,
    "name": "Main Lab",
    "description": "Primary laboratory space"
  },
  "location_id": 1,
  "laboratory_id": 1
}
```

---

### POST `/api/equipment`
**Description:** Create new equipment

**Request Body:**
```json
{
  "name": "PCR Thermocycler",
  "description": "Bio-Rad T100 Thermal Cycler for PCR amplification",
  "propertyNumber": "EQ-2023-001",
  "location_id": 1,
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Equipment created successfully"
}
```

---

### PUT `/api/equipment/:id`
**Description:** Update existing equipment

**Path Parameters:**
- `id` (number, required): Equipment ID

**Request Body:**
```json
{
  "id": 1,
  "name": "PCR Thermocycler",
  "description": "Bio-Rad T100 Thermal Cycler for PCR amplification",
  "propertyNumber": "EQ-2023-001",
  "location_id": 1,
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Equipment updated successfully"
}
```

---

### DELETE `/api/equipment/:id`
**Description:** Delete equipment and all associated usage records

**Path Parameters:**
- `id` (number, required): Equipment ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "message": "Equipment deleted successfully"
}
```

---

### POST `/api/equipment/:id/schedule`
**Description:** Schedule equipment usage

**Path Parameters:**
- `id` (number, required): Equipment ID

**Request Body:**
```json
{
  "equipmentId": 1,
  "researcherId": 2,
  "description": "PCR amplification for gene cloning",
  "startDate": "2025-01-20T09:00:00.000Z",
  "endDate": "2025-01-20T11:30:00.000Z"
}
```

**Response:**
```json
{
  "message": "Equipment usage scheduled successfully"
}
```

---

### GET `/api/equipment/:id/usage-history`
**Description:** Get paginated usage history for equipment

**Path Parameters:**
- `id` (number, required): Equipment ID

**Query Parameters:**
- `page` (number, optional, default: 1): Page number
- `pageSize` (number, optional, default: 10): Items per page

**Response:**
```json
{
  "usages": [
    {
      "id": 1,
      "equipmentId": 1,
      "researcher": {
        "id": 2,
        "name": "Dr. Maria Silva",
        "email": "maria@lab.com"
      },
      "description": "PCR amplification for gene cloning",
      "startDate": "2025-01-20T09:00:00.000Z",
      "endDate": "2025-01-20T11:30:00.000Z"
    }
  ],
  "nextPage": 2,
  "totalCount": 45
}
```

---

### GET `/api/equipment/:id/calendar`
**Description:** Get calendar view of equipment usage for a specific month

**Path Parameters:**
- `id` (number, required): Equipment ID

**Query Parameters:**
- `year` (number, required): Year (e.g., 2025)
- `month` (number, required): Month (1-12)

**Response:**
```json
{
  "month": 1,
  "year": 2025,
  "days": [
    {
      "dayNumber": 20,
      "appointments": [
        {
          "id": 1,
          "startDate": "2025-01-20T09:00:00.000Z",
          "endDate": "2025-01-20T11:30:00.000Z",
          "researcher": {
            "id": 2,
            "name": "Dr. Maria Silva",
            "email": "maria@lab.com"
          },
          "description": "PCR amplification"
        }
      ]
    }
  ]
}
```

---

### POST `/api/equipment/:id/check-overlap`
**Description:** Check for scheduling conflicts for equipment

**Path Parameters:**
- `id` (number, required): Equipment ID

**Request Body:**
```json
{
  "startDate": "2025-01-20T09:00:00.000Z",
  "endDate": "2025-01-20T11:30:00.000Z",
  "excludeUsageId": null
}
```

**Response:**
```json
{
  "hasOverlap": true,
  "conflictingUsages": [
    {
      "id": 5,
      "equipmentId": 1,
      "researcher": {
        "id": 3,
        "name": "Dr. João Santos",
        "email": "joao@lab.com"
      },
      "description": "DNA gel electrophoresis",
      "startDate": "2025-01-20T10:00:00.000Z",
      "endDate": "2025-01-20T12:00:00.000Z"
    }
  ]
}
```

---

## 3. Samples

### GET `/api/samples`
**Description:** Get all samples for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "E. coli K12",
    "description": "Wild-type strain"
  }
]
```

---

### GET `/api/samples/:id`
**Description:** Get detailed sample information

**Path Parameters:**
- `id` (number, required): Sample ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "E. coli K12",
  "description": "Wild-type strain for cloning experiments",
  "origin": "ATCC 10798",
  "isolationDate": "2024-01-15T00:00:00.000Z",
  "latitude": -23.5505,
  "longitude": -46.6333,
  "subcultureMedium": "LB Agar",
  "subcultureIntervalDays": 7,
  "laboratory_id": 1,
  "createdAt": "2024-01-15T10:30:00.000Z"
}
```

---

### POST `/api/samples`
**Description:** Create a new sample (automatically creates first replica)

**Request Body:**
```json
{
  "name": "E. coli K12",
  "description": "Wild-type strain for cloning experiments",
  "origin": "ATCC 10798",
  "isolationDate": "2024-01-15T00:00:00.000Z",
  "latitude": -23.5505,
  "longitude": -46.6333,
  "subcultureMedium": "LB Agar",
  "subcultureIntervalDays": 7,
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Sample created successfully with first replica"
}
```

---

### PUT `/api/samples/:id`
**Description:** Update an existing sample

**Path Parameters:**
- `id` (number, required): Sample ID

**Request Body:**
```json
{
  "name": "E. coli K12",
  "description": "Wild-type strain for cloning experiments",
  "origin": "ATCC 10798",
  "isolationDate": "2024-01-15T00:00:00.000Z",
  "latitude": -23.5505,
  "longitude": -46.6333,
  "subcultureMedium": "LB Agar",
  "subcultureIntervalDays": 7,
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Sample updated successfully"
}
```

---

### DELETE `/api/samples/:id`
**Description:** Delete a sample and all associated replicas

**Path Parameters:**
- `id` (number, required): Sample ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "message": "Sample deleted successfully"
}
```

---

## 4. Replicas

### GET `/api/replicas`
**Description:** Get all replicas for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "E. coli K12 - Replica 1",
    "sample": {
      "id": 1,
      "name": "E. coli K12"
    },
    "location": {
      "id": 2,
      "name": "Freezer A"
    },
    "status": "active"
  }
]
```

---

### GET `/api/replicas/sample/:sampleId`
**Description:** Get all replicas for a specific sample

**Path Parameters:**
- `sampleId` (number, required): Sample ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "E. coli K12 - Replica 1",
    "sample": {
      "id": 1,
      "name": "E. coli K12"
    },
    "sample_id": 1,
    "location": {
      "id": 2,
      "name": "Freezer A",
      "description": "Main freezer"
    },
    "location_id": 2,
    "status": "active",
    "lastSubcultureDate": "2025-01-10T00:00:00.000Z",
    "nextSubcultureDate": "2025-01-17T00:00:00.000Z",
    "laboratory_id": 1
  }
]
```

---

### GET `/api/replicas/:id`
**Description:** Get detailed replica information

**Path Parameters:**
- `id` (number, required): Replica ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "E. coli K12 - Replica 1",
  "sample": {
    "id": 1,
    "name": "E. coli K12",
    "description": "Wild-type strain",
    "subcultureMedium": "LB Agar",
    "subcultureIntervalDays": 7
  },
  "sample_id": 1,
  "location": {
    "id": 2,
    "name": "Freezer A",
    "description": "Main freezer"
  },
  "location_id": 2,
  "status": "active",
  "lastSubcultureDate": "2025-01-10T00:00:00.000Z",
  "nextSubcultureDate": "2025-01-17T00:00:00.000Z",
  "laboratory_id": 1
}
```

---

### POST `/api/replicas`
**Description:** Create a new replica

**Request Body:**
```json
{
  "name": "E. coli K12 - Replica 2",
  "sample_id": 1,
  "location_id": 2,
  "status": "active",
  "lastSubcultureDate": "2025-01-15T00:00:00.000Z",
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Replica created successfully"
}
```

---

### PUT `/api/replicas/:id`
**Description:** Update an existing replica

**Path Parameters:**
- `id` (number, required): Replica ID

**Request Body:**
```json
{
  "name": "E. coli K12 - Replica 1",
  "sample_id": 1,
  "location_id": 2,
  "status": "active",
  "lastSubcultureDate": "2025-01-15T00:00:00.000Z",
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Replica updated successfully"
}
```

---

### DELETE `/api/replicas/:id`
**Description:** Delete a replica

**Path Parameters:**
- `id` (number, required): Replica ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "message": "Replica deleted successfully"
}
```

---

### POST `/api/replicas/:id/subculture`
**Description:** Record a subculture event (updates lastSubcultureDate and nextSubcultureDate)

**Path Parameters:**
- `id` (number, required): Replica ID

**Request Body:**
```json
{
  "subcultureDate": "2025-01-20T00:00:00.000Z"
}
```

**Response:**
```json
{
  "message": "Subculture recorded successfully"
}
```

---

## 5. Locations

### GET `/api/locations`
**Description:** Get all locations for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "Main Lab",
    "description": "Primary laboratory space",
    "laboratory_id": 1
  },
  {
    "id": 2,
    "name": "Freezer A",
    "description": "Main freezer at -80°C",
    "laboratory_id": 1
  }
]
```

---

## 6. Researchers

### GET `/api/researchers`
**Description:** Get all researchers for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "Dr. Ana Coordinator",
    "email": "ana@lab.com",
    "role": "coordinator",
    "laboratory_id": 1
  },
  {
    "id": 2,
    "name": "Dr. Maria Silva",
    "email": "maria@lab.com",
    "role": "researcher",
    "laboratory_id": 1
  }
]
```

---

## 7. Routines

### GET `/api/routines`
**Description:** Get all routines for a laboratory

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "70% Alcohol Preparation",
    "description": "Weekly preparation of 70% ethanol solution",
    "scheduleType": "recurring",
    "recurrence": {
      "frequency": "weekly",
      "interval": 1,
      "daysOfWeek": [1, 3, 5],
      "startDate": "2025-01-01T00:00:00.000Z"
    },
    "assignedTo": [2, 3],
    "assignedToNames": ["Dr. Maria Silva", "Dr. João Santos"],
    "createdBy": 1,
    "createdByName": "Dr. Ana Coordinator",
    "laboratory_id": 1
  }
]
```

---

### GET `/api/routines/:id`
**Description:** Get detailed routine information including materials, equipment, and steps

**Path Parameters:**
- `id` (number, required): Routine ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "id": 1,
  "name": "70% Alcohol Preparation",
  "description": "Weekly preparation of 70% ethanol solution",
  "scheduleType": "recurring",
  "recurrence": {
    "frequency": "weekly",
    "interval": 1,
    "daysOfWeek": [1, 3, 5],
    "startDate": "2025-01-01T00:00:00.000Z"
  },
  "deadline": null,
  "assignedTo": [2, 3],
  "assignedToNames": ["Dr. Maria Silva", "Dr. João Santos"],
  "createdBy": 1,
  "createdByName": "Dr. Ana Coordinator",
  "materials": [
    {
      "productId": 1,
      "productName": "Ethanol 99.5%",
      "quantity": 0.7,
      "unit": "L"
    },
    {
      "productId": 3,
      "productName": "Distilled Water",
      "quantity": 0.3,
      "unit": "L"
    }
  ],
  "equipment": [
    {
      "equipmentId": 7,
      "equipmentName": "Graduated Cylinder 1L",
      "estimatedDuration": 15,
      "required": true
    }
  ],
  "steps": [
    {
      "order": 1,
      "description": "Measure 700mL of ethanol 99.5% using graduated cylinder",
      "notes": "Ensure cylinder is clean and dry"
    },
    {
      "order": 2,
      "description": "Add 300mL of distilled water",
      "notes": "Add water slowly while stirring"
    },
    {
      "order": 3,
      "description": "Mix thoroughly and transfer to labeled bottle",
      "notes": "Label with date and concentration"
    }
  ],
  "laboratory_id": 1
}
```

---

### GET `/api/routines/upcoming`
**Description:** Get upcoming routine instances for next N days (includes one-time and recurring)

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID
- `days` (number, optional, default: 7): Number of days to look ahead

**Response:**
```json
[
  {
    "id": 1,
    "routineId": 1,
    "routineName": "70% Alcohol Preparation",
    "routineDescription": "Weekly preparation",
    "scheduledDate": "2025-01-20T00:00:00.000Z",
    "dueDate": "2025-01-20T23:59:59.000Z",
    "daysUntilDue": 2,
    "assignedTo": [2, 3],
    "assignedToNames": ["Dr. Maria Silva", "Dr. João Santos"],
    "status": "pending",
    "laboratory_id": 1
  }
]
```

---

### POST `/api/routines`
**Description:** Create a new routine

**Request Body:**
```json
{
  "name": "70% Alcohol Preparation",
  "description": "Weekly preparation of 70% ethanol solution",
  "scheduleType": "recurring",
  "recurrence": {
    "frequency": "weekly",
    "interval": 1,
    "daysOfWeek": [1, 3, 5],
    "startDate": "2025-01-01T00:00:00.000Z"
  },
  "deadline": null,
  "assignedTo": [2, 3],
  "materials": [
    {
      "productId": 1,
      "quantity": 0.7
    }
  ],
  "equipment": [
    {
      "equipmentId": 7,
      "estimatedDuration": 15,
      "required": true
    }
  ],
  "steps": [
    {
      "order": 1,
      "description": "Measure 700mL of ethanol",
      "notes": "Use clean cylinder"
    }
  ],
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "message": "Routine created successfully"
}
```

---

### PUT `/api/routines/:id`
**Description:** Update an existing routine

**Path Parameters:**
- `id` (number, required): Routine ID

**Request Body:** Same as POST `/api/routines` with `id` field included

**Response:**
```json
{
  "message": "Routine updated successfully"
}
```

---

### DELETE `/api/routines/:id`
**Description:** Delete a routine

**Path Parameters:**
- `id` (number, required): Routine ID

**Query Parameters:**
- `laboratory_id` (number, required): Laboratory ID

**Response:**
```json
{
  "message": "Routine deleted successfully"
}
```

---

### POST `/api/routines/:id/execute`
**Description:** Start executing a routine (creates execution record)

**Path Parameters:**
- `id` (number, required): Routine ID

**Request Body:**
```json
{
  "executedBy": 2,
  "laboratory_id": 1
}
```

**Response:**
```json
{
  "executionId": 15,
  "message": "Routine execution started"
}
```

---

### GET `/api/routines/executions/:executionId`
**Description:** Get execution details including step completion

**Path Parameters:**
- `executionId` (number, required): Execution ID

**Response:**
```json
{
  "id": 15,
  "routineId": 1,
  "routineName": "70% Alcohol Preparation",
  "executedBy": {
    "id": 2,
    "name": "Dr. Maria Silva"
  },
  "startedAt": "2025-01-20T14:30:00.000Z",
  "completedAt": null,
  "status": "in_progress",
  "stepsCompleted": [1, 2],
  "totalSteps": 3,
  "materials": [
    {
      "productId": 1,
      "productName": "Ethanol 99.5%",
      "quantity": 0.7,
      "unit": "L"
    }
  ],
  "laboratory_id": 1
}
```

---

### PUT `/api/routines/executions/:executionId/step`
**Description:** Mark a step as completed or uncompleted

**Path Parameters:**
- `executionId` (number, required): Execution ID

**Request Body:**
```json
{
  "stepOrder": 2,
  "completed": true
}
```

**Response:**
```json
{
  "message": "Step updated successfully"
}
```

---

### POST `/api/routines/executions/:executionId/complete`
**Description:** Complete routine execution (deducts materials from inventory)

**Path Parameters:**
- `executionId` (number, required): Execution ID

**Response:**
```json
{
  "message": "Routine execution completed successfully"
}
```

---

### POST `/api/routines/executions/:executionId/cancel`
**Description:** Cancel routine execution

**Path Parameters:**
- `executionId` (number, required): Execution ID

**Response:**
```json
{
  "message": "Routine execution cancelled"
}
```

---

### GET `/api/routines/:id/executions`
**Description:** Get execution history for a routine

**Path Parameters:**
- `id` (number, required): Routine ID

**Query Parameters:**
- `page` (number, optional, default: 1): Page number
- `pageSize` (number, optional, default: 10): Items per page

**Response:**
```json
{
  "executions": [
    {
      "id": 15,
      "executedBy": {
        "id": 2,
        "name": "Dr. Maria Silva"
      },
      "startedAt": "2025-01-20T14:30:00.000Z",
      "completedAt": "2025-01-20T15:00:00.000Z",
      "status": "completed",
      "duration": 30
    }
  ],
  "nextPage": null,
  "totalCount": 5
}
```

---

## Authentication & Authorization Notes

All endpoints require:
1. **Authentication**: Valid JWT token in `Authorization: Bearer <token>` header
2. **Laboratory Context**: User must have access to the specified `laboratory_id`
3. **Role-Based Access**:
   - **Coordinator**: Full access to all operations
   - **Researcher**: Can view, create, and edit; cannot delete
   - **Technician**: Can view and execute routines; limited create/edit access

The backend should enforce multi-tenancy by:
- Always filtering queries by `laboratory_id`
- Validating that the authenticated user has access to the specified laboratory
- Returning 403 Forbidden if user attempts to access resources from other laboratories

---

## Error Responses

All endpoints should return consistent error responses:

```json
{
  "error": "Error type (e.g., ValidationError, NotFoundError, UnauthorizedError)",
  "message": "Human-readable error message",
  "details": {} // Optional: Additional error details
}
```

**HTTP Status Codes:**
- `200 OK`: Successful GET/PUT/DELETE
- `201 Created`: Successful POST
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: User doesn't have permission
- `404 Not Found`: Resource doesn't exist
- `500 Internal Server Error`: Server error

---

## Notes for Backend Implementation

1. **Date/Time Format**: All dates should be in ISO 8601 format (e.g., `2025-01-20T14:30:00.000Z`)
2. **Pagination**: When implementing pagination, always return `nextPage: null` when there are no more pages
3. **Cascade Deletes**:
   - Deleting equipment should delete all usage records
   - Deleting samples should delete all replicas
   - Deleting routines should delete all execution records
4. **Overlap Detection**: The overlap check algorithm should use: `startDate < existingEndDate AND endDate > existingStartDate`
5. **Subculture Calculation**: `nextSubcultureDate = lastSubcultureDate + subcultureIntervalDays`
6. **Material Deduction**: When completing a routine execution, deduct materials from product quantities and update `quantity_used_in_the_last_three_months`

