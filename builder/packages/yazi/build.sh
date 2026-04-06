#!/usr/bin/env bash
set -euo pipefail

cd "$FPB_CLONE_DIR"
CARGO_BUILD_JOBS=1 cargo build --locked --offline --release --bin yazi --bin ya
cp target/release/yazi "$FPB_OUT_DIR/yazi"
cp target/release/ya "$FPB_OUT_DIR/ya"