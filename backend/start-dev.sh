#!/bin/bash
# Load .env and start backend server

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Change to backend directory
cd "$SCRIPT_DIR"

set -a  # automatically export all variables
source .env
set +a

echo "🔧 Environment loaded from .env"
echo "✓ FIREBASE_PROJECT_ID: $FIREBASE_PROJECT_ID"
echo "✓ FIREBASE_CREDENTIALS_PATH: $FIREBASE_CREDENTIALS_PATH"
echo "✓ SERVER_PORT: $SERVER_PORT"
echo "✓ Working directory: $(pwd)"
echo ""

# Verify credentials file exists
if [ ! -f "$FIREBASE_CREDENTIALS_PATH" ]; then
  echo "❌ ERROR: Firebase credentials not found at $FIREBASE_CREDENTIALS_PATH"
  echo "   Looking for file at: $PROJECT_ROOT/$FIREBASE_CREDENTIALS_PATH"
  exit 1
fi

echo "✅ Firebase credentials found"
echo ""

go run ./cmd/server
