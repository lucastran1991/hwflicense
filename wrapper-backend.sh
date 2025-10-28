#!/bin/bash
# Wrapper for Backend with log tagging
./server 2>&1 | while IFS= read -r line; do
  echo "[BE] $(date '+%Y-%m-%d %H:%M:%S') $line"
done

