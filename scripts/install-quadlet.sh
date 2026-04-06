#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
target_home="${QUADLET_TARGET_HOME:-$HOME}"
xdg_config_home="${XDG_CONFIG_HOME:-$target_home/.config}"
unit_dir="$xdg_config_home/containers/systemd"
repo_server_root="${REPO_SERVER_ROOT:-$target_home/fedora-package-repo/repo}"
container_runtime="${CONTAINER_RUNTIME_COMMAND:-podman}"
image_name="${REPO_IMAGE_NAME:-localhost/fedora-package-repo-nginx:latest}"

mkdir -p "$unit_dir"
mkdir -p "$repo_server_root"

"$container_runtime" build -t "$image_name" "$repo_root/container/nginx"
ln -sfn "$repo_root/quadlet/fedora-package-repo.network" "$unit_dir/fedora-package-repo.network"
ln -sfn "$repo_root/quadlet/fedora-package-repo.container" "$unit_dir/fedora-package-repo.container"
systemctl --user daemon-reload

cat <<EOF
Installed Quadlet files.

Next step:
  systemctl --user start fedora-package-repo.service
EOF