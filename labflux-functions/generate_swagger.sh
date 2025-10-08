#!/bin/bash

echo "📚 Generating Swagger Documentation for LabFlux Functions"

# Create docs directory if it doesn't exist
mkdir -p docs

# Try to use swag CLI if available, otherwise use Docker
if command -v swag &> /dev/null; then
    echo "Using local swag CLI..."
    swag init -g cmd/main.go --output docs --parseDependency --parseInternal
else
    echo "Generating swagger documentation using Docker..."
    docker run --rm -v $(pwd):/src -w /src swaggerapi/swagger-codegen-cli:2.4.28 generate -i /src/cmd/main.go -l swagger -o /src/docs
fi

# Alternative: Generate a basic swagger.json manually for now
cat > ./docs/swagger.json << 'EOF'
{
  "swagger": "2.0",
  "info": {
    "title": "LabFlux Functions API",
    "version": "1.0.0",
    "description": "A serverless API for laboratory management with multi-tenancy support"
  },
  "host": "localhost:8080",
  "basePath": "/api",
  "schemes": ["http"],
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Type 'Bearer' followed by a space and JWT token"
    }
  },
  "paths": {
    "/laboratories": {
      "post": {
        "tags": ["laboratories"],
        "summary": "Create Laboratory",
        "description": "Creates a new laboratory and assigns the current user as admin",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "in": "body",
            "name": "laboratory",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name": {"type": "string"},
                "description": {"type": "string"},
                "address": {"type": "string"}
              }
            }
          }
        ],
        "responses": {
          "201": {"description": "Laboratory created successfully"},
          "400": {"description": "Invalid request"},
          "401": {"description": "Unauthorized"}
        }
      }
    },
    "/laboratories/{laboratoryId}": {
      "get": {
        "tags": ["laboratories"],
        "summary": "Get Laboratory",
        "description": "Get laboratory details by ID",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "in": "path",
            "name": "laboratoryId",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "200": {"description": "Laboratory details"},
          "404": {"description": "Laboratory not found"}
        }
      }
    }
  }
}
EOF

# Also create docs.go file that's needed
cat > ./docs/docs.go << 'EOF'
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "LabFlux Functions API",
    "version": "1.0.0",
    "description": "A serverless API for laboratory management with multi-tenancy support"
  },
  "host": "localhost:8080",
  "basePath": "/api"
}`

func init() {
    swag.Register(swag.Name, &s{})
}

type s struct{}

func (s *s) ReadDoc() string {
    return docTemplate
}
EOF

echo "✅ Swagger documentation generated in ./docs/"
echo "📋 Files created:"
ls -la docs/

echo ""
echo "🚀 Next steps:"
echo "1. Run: docker compose up"
echo "2. Visit: http://localhost:8080/swagger/"
echo "3. Test the API endpoints interactively!"