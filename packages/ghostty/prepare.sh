#!/usr/bin/env bash
set -euo pipefail

resolve_ref() {
  if [[ -n "${FPB_REF:-}" ]]; then
    printf '%s\n' "$FPB_REF"
    return
  fi

  git ls-remote --refs --tags "$FPB_REPO_URL" 'v*' \
    | awk '{print $2}' \
    | sed 's#refs/tags/##' \
    | sort -V \
    | tail -n 1
}

ref="$(resolve_ref)"
if [[ -z "$ref" ]]; then
  printf 'failed to resolve a Ghostty release tag from %s\n' "$FPB_REPO_URL" >&2
  exit 1
fi

rm -rf "$FPB_CLONE_DIR" /tmp/ghostty-zig-cache
git clone --depth 1 --branch "$ref" "$FPB_REPO_URL" "$FPB_CLONE_DIR"

cd "$FPB_CLONE_DIR"
mkdir -p /tmp/ghostty-zig-cache
export ZIG_GLOBAL_CACHE_DIR=/tmp/ghostty-zig-cache

while IFS= read -r url; do
  [[ -n "$url" ]] || continue

  attempt=1
  while true; do
    if zig fetch "$url" >/dev/null 2>&1; then
      break
    fi

    if [[ "$attempt" -ge 5 ]]; then
      printf 'Ghostty Zig cache prefetch failed for %s after %s attempts\n' "$url" "$attempt" >&2
      exit 1
    fi

    attempt="$((attempt + 1))"
  done
done < build.zig.zon.txt