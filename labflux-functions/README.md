# Neva Functions - Go Serverless API

A Go-based serverless API for the Neva Management system using Google Cloud Functions.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Google Cloud CLI (for deployment)

### Local Development

1. **Clone and setup**
   ```bash
   cd neva-functions
   cp .env.example .env
   ```

2. **Start local environment**
   ```bash
   docker-compose up
   ```

3. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # API ping
   curl http://localhost:8080/api/ping
   ```

### Available Endpoints

#### Public Endpoints
- `GET /health` - Health check
- `GET /api/ping` - API ping

#### Protected Endpoints (require JWT token)
- `POST /api/laboratories` - Create laboratory
- `GET /api/laboratories/{id}` - Get laboratory
- `POST /api/laboratories/{id}/invitations` - Create invitation
- `GET /api/laboratories/{id}/invitations` - List invitations

### Authentication

Add JWT token to requests:
```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/laboratories/1
```

## 📁 Project Structure

```
neva-functions/
├── cmd/
│   └── main.go              # Local development server
├── functions/
│   └── laboratory.go        # Laboratory functions
├── pkg/
│   ├── database/            # Database connection
│   ├── middleware/          # Auth middleware
│   └── models/              # GORM models
├── internal/
│   └── router.go            # Chi router setup
├── docker-compose.yml       # Local development
└── Dockerfile              # Container build
```

## 🚀 Deployment to GCP

### Deploy Individual Function
```bash
gcloud functions deploy create-laboratory \
  --runtime go121 \
  --trigger-http \
  --source . \
  --entry-point CreateLaboratory \
  --set-env-vars DATABASE_URL=$DATABASE_URL,JWT_SECRET=$JWT_SECRET
```

### Deploy All Functions
```bash
# Deploy each function endpoint
gcloud functions deploy laboratory-api \
  --runtime go121 \
  --trigger-http \
  --source . \
  --entry-point HandleRequest
```

## 🔧 Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret
- `PORT` - Server port (default: 8080)

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 📊 Cost Estimation

For 50 users (~30k requests/month):
- **GCP Cloud Functions**: $1-3/month
- **Neon Database**: Free tier
- **Total**: $1-3/month