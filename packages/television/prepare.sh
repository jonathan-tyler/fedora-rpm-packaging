#!/usr/bin/env bash
set -euo pipefail

rm -rf "$FPB_CLONE_DIR"

if [[ -n "${FPB_REF:-}" ]]; then
  git clone --depth 1 --branch "$FPB_REF" "$FPB_REPO_URL" "$FPB_CLONE_DIR"
else
  git clone --depth 1 "$FPB_REPO_URL" "$FPB_CLONE_DIR"
fi

cd "$FPB_CLONE_DIR"
mkdir -p .cargo
cargo vendor vendor > .cargo/config.toml