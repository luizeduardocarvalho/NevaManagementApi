#!/bin/bash
# Create a deployment archive excluding infra directory

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
ARCHIVE_PATH="$PROJECT_ROOT/deployment.zip"

cd "$PROJECT_ROOT"

# Remove old archive if exists
rm -f "$ARCHIVE_PATH"

# Create zip excluding infra, .git, and other unnecessary files
zip -r "$ARCHIVE_PATH" . \
  -x "infra/*" \
  -x ".git/*" \
  -x ".env" \
  -x "*.log" \
  -x "main" \
  -x "deployment.zip" \
  -x "scripts/*" \
  -x "docs/*" \
  -x "*.md"

echo "Deployment archive created at: $ARCHIVE_PATH"
