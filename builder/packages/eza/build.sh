#!/usr/bin/env bash
set -euo pipefail

cd "$FPB_CLONE_DIR"
export CARGO_NET_OFFLINE=true
cargo build --locked --offline --release --bin eza
cp target/release/eza "$FPB_OUT_DIR/eza"