#!/usr/bin/env bash
set -euo pipefail

if [[ "$(id -u)" -ne 0 ]]; then
  echo "error: run this script as root" >&2
  exit 1
fi

service_user="${REPO_SERVICE_USER:-reposvc}"

if ! id "$service_user" >/dev/null 2>&1; then
  useradd -m -s /usr/sbin/nologin "$service_user"
fi

loginctl enable-linger "$service_user"

service_home="$(getent passwd "$service_user" | cut -d: -f6)"
repo_root="${REPO_SERVICE_ROOT:-$service_home/fedora-package-repo}"

install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root"
install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root/repo"
install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root/.repo-sync"

cat <<EOF
Prepared service user $service_user.

Next step:
  sudo -iu $service_user bash -lc 'cd $(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd) && bash scripts/install-quadlet.sh'
EOF