#!/usr/bin/env bash
set -euo pipefail

if [[ "$(id -u)" -ne 0 ]]; then
  echo "error: run this script as root" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
service_user="${REPO_SERVICE_USER:-reposvc}"

dnf install -y podman systemd
"$repo_root/scripts/setup-service-user.sh"
runuser -u "$service_user" -- bash -lc "cd '$repo_root' && bash scripts/install-quadlet.sh"

cat <<EOF
WSL repo-service bootstrap completed for $service_user.

Final step:
  sudo -iu $service_user systemctl --user enable --now fedora-package-repo.service
EOF