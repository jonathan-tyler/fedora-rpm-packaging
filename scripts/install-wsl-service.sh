#!/usr/bin/env bash
set -euo pipefail

if [[ "$(id -u)" -ne 0 ]]; then
  echo "error: run this script as root" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"$repo_root/scripts/bootstrap-fedora-service.sh"

cat <<EOF
WSL wrapper completed.

Final step:
  start the generated service through the helper bootstrap path
EOF