# Repository Guidelines (LabFlux – Go / GCP / Pulumi)

## General rules

- Do not apologize  
- Do not thank me  
- Talk to me like a human  
- Verify information before making changes  
- Preserve existing code structures and patterns  
- Provide concise and relevant responses  
- Prefer minimal changes that fully solve the problem  
- Verify all information before making changes  
- Do not change behavior that is not explicitly in scope  

### You will be penalized if you:

- Skip steps in your thought process  
- Add placeholders or TODOs for other developers  
- Deliver code that is not production-ready  
- Break multi-tenancy guarantees or data isolation  
- Introduce regressions to existing endpoints or flows  

I'm tipping **$9000** for an optimal, elegant, minimal, world-class solution that meets all specifications.  
Your code changes should be specific and complete.  
Think through the problem step-by-step.

---

# More General Rules

## Premature Optimization is the Root of All Evil

Most code is not performance-critical.  
Optimize only when it truly matters.

## Principle of Least Astonishment

Code should behave in the most expected, least surprising way.  
Match names, comments, and behavior.  
Avoid side effects that surprise maintainers or users.

---

# Project Context (LabFlux)

LabFlux is a multi-tenant laboratory management API written in **Go**, deployed as a **Google Cloud Function (Gen 2)**, and managed with **Pulumi**.

### Key constraints:

- Strict multi-tenant isolation using `laboratory_id`  
- Serverless execution model (cold starts, limited resources)  
- Cost-optimized architecture  
- Clerk-based authentication  
- Domain-driven organization via `functions/`  

### Stack:

- **Go 1.23+**
- **Chi router** (`github.com/go-chi/chi/v5`)
- **GORM** + PostgreSQL
- **Clerk JWT authentication**
- **Pulumi (Go)** for IaC
- **Google Cloud Functions (Gen 2)**
- **Neon PostgreSQL**

---

# Project Structure & Module Organization

```text
labflux-api/
├── function.go
├── cmd/
│   └── migrate/main.go
├── functions/
│   ├── auth.go
│   ├── laboratory.go
│   ├── products.go
│   ├── equipment.go
│   └── equipment_usage.go
├── internal/
│   └── router.go
├── pkg/
│   ├── database/
│   ├── middleware/
│   ├── migrations/
│   └── models/
├── infra/
│   └── main.go
└── scripts/
```

### Rules:

- **Domain handlers** live in `functions/*`
- **Middleware, DB, shared utilities** go under `pkg/*`
- **Routing** lives in `internal/router.go`
- **Models** currently grouped under `pkg/models/laboratory.go`
- **Do not add top-level folders** unless architecture requires it

---

# Build, Test, and Development Commands

### Local Dev

```bash
go fmt ./...
go vet ./...
go test ./...

docker compose up -d postgres
export DATABASE_URL="postgres://dev:dev@localhost:5433/neva_local?sslmode=disable"

go run cmd/migrate/main.go
```

### Deployment (Pulumi)

```bash
cd infra
pulumi up
```

### Rules

- Use `go vet` and fix warnings for any changed code
- Do not add new build tools; follow Go’s minimal tooling philosophy
- Preserve GCF entrypoint signature (`HandleRequest`)

---

# Multi-Tenancy & Data Isolation (CRITICAL)

Breaking multi-tenancy is unacceptable.

### Rules:

1. **All queries must filter by `laboratory_id`:**
   ```go
   db.Where("laboratory_id = ?", labID).Find(&items)
   ```

2. **Never trust client-provided `laboratoryId`** if authenticated claims include a lab ID  
   - Reject mismatched lab IDs  
   - Default to claim-based scoping  

3. **All newly-created domain records must set `laboratory_id` from JWT claims**  
   not from request bodies.

4. No cross-lab reads, writes, reports, or analytics unless explicitly defined.

---

# Go / API Development Rules

You are a senior Go backend engineer. Write code like one.

## Code Style

- Always formatted via `gofmt`
- Keep functions small and composable
- Reject unnecessary abstractions
- Avoid global mutable state
- Use early returns for error handling

## Naming

- Packages: lowercase (no underscores)
- Exported: `PascalCase`
- Unexported & locals: `camelCase`
- Receivers: short (`h`, `s`, `db`)
- Errors: always named `err`

## Recommended Go Patterns

- Use `context.Context` everywhere
- Prefer explicit types for clarity, `var` when the type is obvious
- Wrap errors with context using `fmt.Errorf("message: %w", err)`

---

# HTTP Routing & Middleware

Chi router lives in `internal/router.go`.

### Domain Routing Pattern

```go
func RegisterProductRoutes(r chi.Router) {
    r.Route("/products", func(r chi.Router) {
        r.Get("/GetAll", GetAllProducts)
        r.Post("/Create", CreateProduct)
    })
}
```

- All protected routes go inside the Auth group
- Public routes include:
  - `/health`
  - `/api/ping`
  - Clerk webhooks  
  - Invitation acceptance

### Middleware Rules

- Auth middleware:
  - Parses JWT (`Authorization: Bearer <token>`)
  - Supports local test tokens
  - Fallback to Clerk JWT validation
  - Injects claims into context for handler use

- Handlers must **never** parse JWT manually.

---

# GORM / Database Rules

### Use the global DB singleton:
```go
db := database.GetDB()
```

### Query Rules

- Always use parameterized queries:
  ```go
  db.Where("field = ?", value)
  ```

- Always include `laboratory_id` in filters
- Handle `err` after every DB call
- Use `Preload()` only when needed
- Use transactions for multi-record operations:
  ```go
  db.Transaction(func(tx *gorm.DB) error { ... })
  ```

### Soft deletes

- All models use GORM’s soft delete
- Do not manually delete records with raw SQL

---

# Error Handling & Response Rules

Use Chi render or middleware helpers:

```go
middleware.ErrorResponse(w, r, http.StatusBadRequest, "Invalid request")
```

or

```go
render.Status(r, http.StatusBadRequest)
render.JSON(w, r, map[string]string{"error": "validation failed"})
```

### Rules:

- Never expose internal errors or SQL details
- Use proper HTTP status codes (400, 401, 403, 404, 409, 422, 500)
- Be explicit, consistent, and safe

---

# Testing Guidelines

Use:

- `testing`
- `httptest`
- Table-driven tests

### Rules:

- Test filenames mirror source (`products.go` → `products_test.go`)
- Test names describe behavior (`TestGetAllProducts_ReturnsPaginatedList`)
- For handlers:
  - Use `httptest.NewRecorder()`
  - Validate status and key JSON fields

- Avoid external testing frameworks unless already present

---

# Commit & PR Guidelines

- Use lowercase imperative commit subjects:

  ```
  feat: add equipment reservation conflict detection
  fix: validate lab id on product edit
  ```

- PRs must include:
  - Summary of changes
  - Affected modules
  - Steps to repro / verify
  - Any new environment variables
  - Screenshots or example responses for API changes

Avoid large mixed-purpose PRs.

---

# Security & Configuration

### Secrets

Never commit real secrets:

- `CLERK_SECRET_KEY`
- `CLERK_WEBHOOK_SECRET`
- `DATABASE_URL`

### Environment Variables (Prod)

- DATABASE_URL  
- CLERK_PUBLISHABLE_KEY  
- CLERK_SECRET_KEY  
- CLERK_WEBHOOK_SECRET  

### Security Rules

- Webhook signature verification must be implemented
- Validate all user input
- Clerk validation must not be bypassed in production
- Error messages must not leak system details

---

# Pulumi / Infrastructure Rules

Infrastructure definitions live in `infra/main.go`.

### Cloud Function Requirements

- Gen 2 HTTP Function  
- Runtime: Go 1.22+  
- Entry point: `HandleRequest`  
- Memory: 256MB  
- Timeout: 60s  
- Pulumi manages:
  - Deployment  
  - IAM  
  - Environment variables  

### Rules:

- Never deploy functions manually through the console  
- All infrastructure changes must go through Pulumi  
- No hardcoded secrets  
- Keep IaC declarative and minimal  

---

# Performance & Cost Awareness

Because LabFlux is fully serverless:

- Minimize cold start penalties
- Reduce allocations in hot paths
- Avoid unnecessary DB calls
- Paginate all large queries
- Avoid preload explosions

---

# Final Expectations

When modifying this repository:

- Respect the existing architecture & naming
- Maintain multi-tenancy
- Produce complete, production-ready code
- No TODOs  
- No speculative changes  
- Provide clear, step-by-step explanations when asked  
- Prefer small, surgical diffs  
- Always validate assumptions before coding  
