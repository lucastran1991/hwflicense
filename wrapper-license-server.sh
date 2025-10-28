#!/bin/bash
# Wrapper for License Server with log tagging
./license-server 2>&1 | while IFS= read -r line; do
  echo "[LS] $(date '+%Y-%m-%d %H:%M:%S') $line"
done

