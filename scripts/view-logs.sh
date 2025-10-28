#!/bin/bash
# View unified logs with optional filtering
# Usage: ./scripts/view-logs.sh [BE|FE|LS]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG_FILE="$PROJECT_ROOT/logs/system.log"

# Colors for service tags
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

if [ -z "$1" ]; then
  # Show all logs with color coding
  tail -f "$LOG_FILE" 2>/dev/null | \
    sed "s/\[LS\]/${BLUE}[LS]${NC}/g" | \
    sed "s/\[BE\]/${GREEN}[BE]${NC}/g" | \
    sed "s/\[FE\]/${YELLOW}[FE]${NC}/g"
else
  # Filter by service
  if [ "$1" == "BE" ] || [ "$1" == "FE" ] || [ "$1" == "LS" ]; then
    tail -f "$LOG_FILE" 2>/dev/null | grep "\[$1\]"
  else
    echo "Usage: $0 [BE|FE|LS]"
    echo ""
    echo "Examples:"
    echo "  $0       # Show all logs with colors"
    echo "  $0 BE    # Show only Backend logs"
    echo "  $0 FE    # Show only Frontend logs"
    echo "  $0 LS    # Show only License Server logs"
    exit 1
  fi
fi

