#!/usr/bin/env bash

set -euo pipefail

docker compose exec newsletter-assignment go test -race -v -count=1 -timeout 50s -coverpkg=./... -coverprofile=./tmp/coverage ./...

docker compose exec newsletter-assignment go tool cover -html=./tmp/coverage -o ./tmp/coverage.html
