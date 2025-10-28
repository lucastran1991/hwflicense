#!/bin/bash
# Wrapper for Frontend with log tagging
cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
  echo "[FE] $(date '+%Y-%m-%d %H:%M:%S') Installing dependencies..."
  npm install --production
fi

# Start frontend
npm start 2>&1 | while IFS= read -r line; do
  echo "[FE] $(date '+%Y-%m-%d %H:%M:%S') $line"
done

