#!/usr/bin/env bash

set -euo pipefail

docker compose exec newsletter-assignment go test -race -v -count=1 -timeout 50s -coverpkg=./... -coverprofile=./tmp/coverage ./... | sed '
    s/^=== RUN.*/\x1b[1;36m&\x1b[0m/
    s/^--- PASS:.*/\x1b[1;32m&\x1b[0m/
    s/^--- FAIL:.*/\x1b[1;31m&\x1b[0m/
    s/^FAIL$/\x1b[1;31m&\x1b[0m/
    s/^PASS$/\x1b[1;32m&\x1b[0m/
    s/^ok.*/\x1b[1;32m&\x1b[0m/
    s/coverage:.*/\x1b[1;33m&\x1b[0m/
'

# Generate coverage HTML
docker compose exec newsletter-assignment go tool cover -html=./tmp/coverage -o ./tmp/coverage.html