#!/bin/bash

echo "🐳 Testing Docker Setup with Swagger Documentation"

# Clean up any existing containers
echo "1. Cleaning up existing containers..."
docker compose down

# Generate Swagger docs
echo "2. Generating Swagger documentation..."
./generate_swagger.sh

echo ""
echo "3. Building and starting containers..."
docker compose up --build -d

# Wait for services to start
echo "4. Waiting for services to start..."
sleep 10

# Test health endpoints
echo "5. Testing endpoints..."

echo "   📊 Health Check:"
curl -s http://localhost:8080/health | jq '.' || echo "Failed to reach health endpoint"

echo "   📡 API Ping:"
curl -s http://localhost:8080/api/ping | jq '.' || echo "Failed to reach ping endpoint"

echo "   📚 Swagger JSON:"
curl -s http://localhost:8080/swagger/doc.json | jq '.info' || echo "Failed to reach swagger JSON"

# Show running containers
echo ""
echo "6. Container status:"
docker compose ps

# Show logs if there are issues
echo ""
echo "7. Recent logs:"
docker compose logs --tail=10 functions

echo ""
echo "✅ Docker setup test complete!"
echo ""
echo "🌐 Available endpoints:"
echo "   • Health Check:  http://localhost:8080/health"
echo "   • API Ping:      http://localhost:8080/api/ping"
echo "   • Swagger UI:    http://localhost:8080/swagger/"
echo "   • Swagger JSON:  http://localhost:8080/swagger/doc.json"
echo ""
echo "🔧 Manual testing (if compose fails):"
echo "   docker run --rm -p 8080:8080 \\"
echo "     -e DATABASE_URL=\"postgres://dev:dev@host.docker.internal:5433/neva_local?sslmode=disable\" \\"
echo "     -e JWT_SECRET=\"your-local-jwt-secret-key\" \\"
echo "     neva-functions-functions"
echo ""
echo "🔧 Troubleshooting:"
echo "   • View logs:     docker compose logs functions"
echo "   • Stop services: docker compose down"
echo "   • Rebuild:       docker compose up --build"