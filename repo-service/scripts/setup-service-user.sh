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
install_root="${REPO_SERVICE_INSTALL_ROOT:-$service_home/fedora-package-repo-service}"

install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root"
install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root/repo"
install -d -o "$service_user" -g "$service_user" -m 0755 "$repo_root/.repo-sync"
install -d -o "$service_user" -g "$service_user" -m 0755 "$install_root"

cat <<EOF
Prepared service user $service_user.

Next step:
  sudo bash $(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/scripts/bootstrap-fedora-service.sh
EOF