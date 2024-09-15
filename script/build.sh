#!/usr/bin/env bash

set -euo pipefail

GO_LDFLAGS=' -w -extldflags "-static"'
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

cd "$(dirname "$0")/.."

echo "Building newsletter-assignment..."

go build -ldflags "$GO_LDFLAGS" -o "../build/newsletter-assignment" -buildvcs=false "/go/src/newsletter-assignment"
echo "Built: $(ls ../build/*)"
