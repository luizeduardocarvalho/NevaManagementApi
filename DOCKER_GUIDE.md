# 🐳 Docker Setup Guide - LabFlux API

Complete guide for running LabFlux API with Docker and Docker Compose.

---

## 📋 Prerequisites

- **Docker** installed ([Get Docker](https://docs.docker.com/get-docker/))
- **Docker Compose** installed (included with Docker Desktop)
- **Clerk Account** for authentication ([Sign up](https://clerk.com))

---

## 🚀 Quick Start

### 1. Set Up Environment Variables

Your `.env` file already exists with Clerk credentials. Verify it contains:

```bash
cat .env
```

Should show:
```env
DATABASE_URL=postgresql://...
CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
CLERK_WEBHOOK_SECRET=whsec_...
PORT=8080
```

### 2. Start All Services

```bash
# Start both API and PostgreSQL
docker compose up -d

# View logs
docker compose logs -f

# Or view logs for specific service
docker compose logs -f api
docker compose logs -f postgres
```

### 3. Verify Services Are Running

```bash
# Check container status
docker compose ps

# Should show:
# labflux-postgres   running   5433->5432/tcp
# labflux-api        running   8080->8080/tcp

# Test API health
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

### 4. Test Authentication Endpoints

```bash
# Sign up a new user
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@labflux.com",
    "password": "SecurePassword123!",
    "first_name": "Test",
    "last_name": "User"
  }'

# Sign in
curl -X POST http://localhost:8080/api/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@labflux.com",
    "password": "SecurePassword123!"
  }'

# Save the token from the response and use it
TOKEN="your_token_here"

# Get current user
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🛠️ Docker Commands Reference

### Service Management

```bash
# Start services
docker compose up -d

# Start and rebuild (after code changes)
docker compose up -d --build

# Stop services
docker compose down

# Stop and remove volumes (⚠️ deletes database data)
docker compose down -v

# Restart services
docker compose restart

# Restart specific service
docker compose restart api
```

### Viewing Logs

```bash
# All services
docker compose logs -f

# API only
docker compose logs -f api

# PostgreSQL only
docker compose logs -f postgres

# Last 100 lines
docker compose logs --tail=100 api
```

### Container Access

```bash
# Access API container shell
docker compose exec api sh

# Access PostgreSQL container
docker compose exec postgres psql -U dev -d neva_local

# Run commands in API container
docker compose exec api go version
docker compose exec api ls -la
```

### Database Management

```bash
# Access PostgreSQL CLI
docker compose exec postgres psql -U dev -d neva_local

# Run SQL query
docker compose exec postgres psql -U dev -d neva_local -c "SELECT * FROM users;"

# Backup database
docker compose exec postgres pg_dump -U dev neva_local > backup.sql

# Restore database
cat backup.sql | docker compose exec -T postgres psql -U dev -d neva_local
```

### Cleanup

```bash
# Stop and remove containers
docker compose down

# Stop and remove volumes (⚠️ deletes all data)
docker compose down -v

# Remove all unused Docker resources
docker system prune -a
```

---

## 🔧 Development Workflow

### Making Code Changes

The API container has hot-reload via volume mounting:

```bash
# 1. Make changes to your code
# 2. Rebuild and restart the API container
docker compose up -d --build api

# 3. View logs to verify changes
docker compose logs -f api
```

### Running Tests

```bash
# Run tests inside container
docker compose exec api go test ./...

# Run specific test
docker compose exec api go test -v ./functions -run TestSignUp
```

### Rebuilding After Dependency Changes

```bash
# If you changed go.mod or go.sum
docker compose build --no-cache api
docker compose up -d api
```

---

## 📊 Container Architecture

```
┌─────────────────────────────────────────┐
│          labflux-network (bridge)       │
│                                         │
│  ┌──────────────┐    ┌───────────────┐ │
│  │ labflux-api  │───▶│ labflux-      │ │
│  │ Port: 8080   │    │ postgres      │ │
│  │              │    │ Port: 5432    │ │
│  └──────────────┘    │ (5433 host)   │ │
│                      └───────────────┘ │
│                                         │
└─────────────────────────────────────────┘
         │
         ▼
  Host: localhost:8080
```

---

## 🔍 Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change the port in .env
echo "PORT=8081" >> .env
docker compose up -d --build api
```

### Database Connection Failed

```bash
# Check PostgreSQL is healthy
docker compose ps postgres

# View PostgreSQL logs
docker compose logs postgres

# Restart PostgreSQL
docker compose restart postgres
```

### API Container Won't Start

```bash
# View detailed logs
docker compose logs api

# Check if binary built correctly
docker compose exec api ls -la

# Rebuild from scratch
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Clerk Authentication Errors

```bash
# Verify environment variables
docker compose exec api env | grep CLERK

# Check .env file has correct Clerk keys
cat .env

# Restart API after updating .env
docker compose restart api
```

### Database Persistence Issues

```bash
# Check volumes
docker volume ls

# Inspect postgres volume
docker volume inspect labflux-api_postgres_data

# Remove and recreate volume (⚠️ deletes data)
docker compose down -v
docker compose up -d
```

---

## 🔐 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://dev:dev@postgres:5432/neva_local?sslmode=disable` |
| `PORT` | API server port | `8080` |
| `CLERK_SECRET_KEY` | Clerk secret key | `sk_test_...` |
| `CLERK_PUBLISHABLE_KEY` | Clerk publishable key | `pk_test_...` |
| `CLERK_WEBHOOK_SECRET` | Clerk webhook secret | `whsec_...` |

---

## 📝 Useful Docker Compose Commands

```bash
# Show running containers
docker compose ps

# Show images
docker compose images

# Show resource usage
docker stats

# Execute command in running container
docker compose exec <service> <command>

# Run one-off command in new container
docker compose run --rm api go version

# View container configuration
docker compose config

# Pull latest images
docker compose pull

# Build without cache
docker compose build --no-cache

# Scale service (not recommended for API)
docker compose up -d --scale api=2
```

---

## 🎯 Production Considerations

For production deployment:

1. **Use production database URL** in `.env`
2. **Set secure Clerk credentials** (production keys)
3. **Enable SSL/TLS** for database connections
4. **Set up proper logging** and monitoring
5. **Use container orchestration** (Kubernetes, Cloud Run, etc.)
6. **Implement health checks** (already configured)
7. **Use secrets management** (don't commit .env to git)

---

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Clerk Documentation](https://clerk.com/docs)
- [PostgreSQL Docker Image](https://hub.docker.com/_/postgres)
- [LabFlux API Endpoints](./ENDPOINTS.md)
- [LabFlux TODO List](./TODO.md)

---

## 🆘 Getting Help

If you encounter issues:

1. Check logs: `docker compose logs -f`
2. Verify environment variables: `docker compose exec api env`
3. Check container status: `docker compose ps`
4. Review [ENDPOINTS.md](./ENDPOINTS.md) for API documentation
5. Review [TODO.md](./TODO.md) for known issues

---

**Last Updated:** 2025-01-15
**Docker Compose Version:** 3.8
**Go Version:** 1.23
**PostgreSQL Version:** 15
