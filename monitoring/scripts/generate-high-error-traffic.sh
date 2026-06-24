#!/usr/bin/env bash
set -u

while true; do
  curl -sS -o /dev/null http://localhost:8080/notes

  curl -sS -o /dev/null \
    -X POST \
    -H "Content-Type: application/json" \
    --data '{"title":' \
    http://localhost:8080/notes

  sleep 1
done
