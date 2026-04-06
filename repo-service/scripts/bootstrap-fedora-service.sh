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

service_home="$(getent passwd "$service_user" | cut -d: -f6)"
install_root="${REPO_SERVICE_INSTALL_ROOT:-$service_home/fedora-package-repo-service}"
runtime_uid="$(id -u "$service_user")"
runtime_dir="/run/user/$runtime_uid"
bus_address="unix:path=$runtime_dir/bus"

rm -rf "$install_root"
install -d -o "$service_user" -g "$service_user" -m 0755 "$install_root"
tar --exclude=.git -C "$repo_root" -cf - . | tar -xf - -C "$install_root"
chown -R "$service_user:$service_user" "$install_root"

systemctl start "user@$runtime_uid.service"

runuser -u "$service_user" -- env \
  XDG_RUNTIME_DIR="$runtime_dir" \
  DBUS_SESSION_BUS_ADDRESS="$bus_address" \
  bash -lc "cd '$install_root' && bash scripts/install-quadlet.sh"

runuser -u "$service_user" -- env \
  XDG_RUNTIME_DIR="$runtime_dir" \
  DBUS_SESSION_BUS_ADDRESS="$bus_address" \
  systemctl --user start fedora-package-repo.service

cat <<EOF
Fedora host repo-service bootstrap completed for $service_user.

Installed service repo copy:
  $install_root
EOF