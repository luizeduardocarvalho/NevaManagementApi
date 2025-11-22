#!/bin/bash

echo "🚀 Starting LabFlux API with Docker Compose..."
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "Please create .env file with your Clerk credentials."
    echo "You can copy .env.example: cp .env.example .env"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running!"
    echo "Please start Docker Desktop and try again."
    exit 1
fi

echo "✅ Docker is running"
echo "✅ .env file found"
echo ""

# Start services
echo "📦 Building and starting containers..."
docker compose up -d --build

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 5

# Check service status
docker compose ps

echo ""
echo "🔍 Testing API health..."
sleep 2

# Test health endpoint
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health 2>/dev/null)

if [ $? -eq 0 ]; then
    echo "✅ API is healthy: $HEALTH_RESPONSE"
    echo ""
    echo "🎉 LabFlux API is running!"
    echo ""
    echo "📝 Available endpoints:"
    echo "   - Health: http://localhost:8080/health"
    echo "   - Ping: http://localhost:8080/api/ping"
    echo "   - Sign Up: http://localhost:8080/api/auth/signup"
    echo "   - Sign In: http://localhost:8080/api/auth/signin"
    echo ""
    echo "📚 Documentation:"
    echo "   - API Endpoints: ENDPOINTS.md"
    echo "   - Docker Guide: DOCKER_GUIDE.md"
    echo "   - TODO List: TODO.md"
    echo ""
    echo "🔧 Useful commands:"
    echo "   - View logs: docker compose logs -f"
    echo "   - Stop services: docker compose down"
    echo "   - Rebuild: docker compose up -d --build"
else
    echo "⚠️  API is starting... may take a few more seconds"
    echo ""
    echo "Run 'docker compose logs -f api' to see logs"
fi
