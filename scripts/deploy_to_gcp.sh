#!/bin/bash

# Deployment script for LabFlux Cloud Functions
# This script deploys the Go functions to Google Cloud Platform

set -e  # Exit on error

# Configuration
PROJECT_ID="labflux-475521"
REGION="us-central1"
FUNCTION_NAME="labflux-api"
RUNTIME="go123"
ENTRY_POINT="HandleRequest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting LabFlux Cloud Functions deployment...${NC}"

# Change to project root
cd "$(dirname "$0")/.."

# Get secrets from Pulumi config
echo -e "${YELLOW}Retrieving secrets from Pulumi...${NC}"
cd infra
DATABASE_URL=$(pulumi config get --path 'labflux:database-url')
CLERK_PUBLISHABLE_KEY=$(pulumi config get --path 'labflux:clerk-publishable-key')
CLERK_SECRET_KEY=$(pulumi config get --path 'labflux:clerk-secret-key')
CLERK_WEBHOOK_SECRET=$(pulumi config get --path 'labflux:clerk-webhook-secret')
cd ..

# Deploy function
echo -e "${YELLOW}Deploying function to GCP...${NC}"
gcloud functions deploy "$FUNCTION_NAME" \
  --gen2 \
  --runtime="$RUNTIME" \
  --region="$REGION" \
  --source=. \
  --entry-point="$ENTRY_POINT" \
  --trigger-http \
  --allow-unauthenticated \
  --ignore-file=.gcloudignore \
  --set-env-vars="DATABASE_URL=$DATABASE_URL,CLERK_PUBLISHABLE_KEY=$CLERK_PUBLISHABLE_KEY,CLERK_SECRET_KEY=$CLERK_SECRET_KEY,CLERK_WEBHOOK_SECRET=$CLERK_WEBHOOK_SECRET" \
  --project="$PROJECT_ID"

# Get function URL
echo -e "${GREEN}Deployment complete!${NC}"
FUNCTION_URL=$(gcloud functions describe "$FUNCTION_NAME" --region="$REGION" --gen2 --format="value(serviceConfig.uri)" --project="$PROJECT_ID")
echo -e "${GREEN}Function URL: ${FUNCTION_URL}${NC}"

# Test health endpoint
echo -e "${YELLOW}Testing health endpoint...${NC}"
curl -s "${FUNCTION_URL}/health" || echo -e "${RED}Health check failed${NC}"

echo -e "${GREEN}Done!${NC}"
