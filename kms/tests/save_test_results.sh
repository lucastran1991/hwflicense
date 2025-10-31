#!/bin/bash

# Save test results and artifacts to tests folder

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KMS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

cd "$KMS_DIR"

echo "=== Saving test artifacts to tests/ folder ==="

# Create tests directory if it doesn't exist
mkdir -p "$SCRIPT_DIR"

# Copy database
if [ -f "kms_test.db" ]; then
    cp "kms_test.db" "$SCRIPT_DIR/test_database_${TIMESTAMP}.db"
    echo "✓ Database saved: test_database_${TIMESTAMP}.db"
fi

# Copy log files
if [ -f "kms.log" ]; then
    cp "kms.log" "$SCRIPT_DIR/kms_service_${TIMESTAMP}.log"
    echo "✓ Service log saved: kms_service_${TIMESTAMP}.log"
fi

# Save current configuration
echo "=== Test Configuration ===" > "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
echo "Timestamp: $(date)" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
echo "Database: $(ls -lh kms_test.db 2>/dev/null | awk '{print $5}' || echo 'Not found')" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
echo "Service PID: $(cat kms.pid 2>/dev/null || echo 'Not running')" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
echo "" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
echo "Settings:" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"
cat config/setting.json 2>/dev/null >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt" || echo "No setting.json" >> "$SCRIPT_DIR/test_config_${TIMESTAMP}.txt"

echo "✓ Configuration saved: test_config_${TIMESTAMP}.txt"

echo ""
echo "All test artifacts saved to: $SCRIPT_DIR"

