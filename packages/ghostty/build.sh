#!/usr/bin/env bash
set -euo pipefail

cd "$FPB_CLONE_DIR"
zig build \
	--system /tmp/ghostty-zig-cache/p \
	-fno-sys=fontconfig \
	-fno-sys=freetype \
	-fno-sys=harfbuzz \
	-fno-sys=libpng \
	-fno-sys=libxml2 \
	-fno-sys=oniguruma \
	-fno-sys=zlib \
	-Doptimize=ReleaseFast \
	-Dcpu=baseline
cp zig-out/bin/ghostty "$FPB_OUT_DIR/ghostty"