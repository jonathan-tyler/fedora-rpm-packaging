#!/usr/bin/env bash
set -euo pipefail

cd "$FPB_CLONE_DIR"
cargo build --locked --offline --release --bin tv
cp target/release/tv "$FPB_OUT_DIR/tv"