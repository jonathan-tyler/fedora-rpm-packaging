#!/usr/bin/env bash
set -euo pipefail

if [[ "$(id -u)" -ne 0 ]]; then
  echo "error: run this script as root" >&2
  exit 1
fi

if [[ "$#" -gt 1 ]]; then
  echo "usage: sudo bash scripts/import-staged-repo.sh [STAGED_REPO_ROOT]" >&2
  exit 1
fi

staged_repo_root="${1:-${STAGED_REPO_ROOT:-}}"
if [[ -z "$staged_repo_root" ]]; then
  echo "error: provide STAGED_REPO_ROOT as an argument or environment variable" >&2
  exit 1
fi
if [[ ! -d "$staged_repo_root" ]]; then
  echo "error: staged repo root not found: $staged_repo_root" >&2
  exit 1
fi

service_user="${REPO_SERVICE_USER:-reposvc}"
service_home="$(getent passwd "$service_user" | cut -d: -f6)"
live_repo_root="${LIVE_REPO_ROOT:-$service_home/fedora-package-repo/repo}"
temp_root="${LIVE_REPO_TMP_ROOT:-$service_home/fedora-package-repo/.repo-sync}"
exclude_artifacts="${INCLUDE_PUBLIC_ARTIFACTS:-0}"

tar_args=(-C "$staged_repo_root" -cf - .)
if [[ "$exclude_artifacts" != "1" ]]; then
  tar_args=(--exclude='./artifacts' "${tar_args[@]}")
fi

rm -rf "$temp_root"
install -d -o "$service_user" -g "$service_user" -m 0755 "$temp_root"
tar "${tar_args[@]}" | tar -xf - -C "$temp_root"

find "$temp_root" -type d -exec chmod 755 {} +
find "$temp_root" -type f -exec chmod 644 {} +

install -d -o "$service_user" -g "$service_user" -m 0755 "$live_repo_root"
# Preserve the repo root inode so a running bind-mounted container keeps seeing updates.
find "$live_repo_root" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
tar -C "$temp_root" -cf - . | tar -xf - -C "$live_repo_root"

chown -R "$service_user:$service_user" "$live_repo_root"
rm -rf "$temp_root"

cat <<EOF
Imported staged repo into:
  $live_repo_root

Source:
  $staged_repo_root
EOF