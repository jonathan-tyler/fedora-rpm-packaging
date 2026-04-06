#!/usr/bin/env bash
set -euo pipefail

cd "$FPB_CLONE_DIR"

version="$(git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)"
CGO_ENABLED=0 go build -mod=vendor -buildvcs=false -ldflags "-X main.version=${version}" -o "$FPB_OUT_DIR/sesh" ./