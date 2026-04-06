#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
image_tag="localhost/fedora-package-builder-createrepo:latest"
context_dir="${repo_root}/tools/createrepo"

if [[ $# -ne 2 || "${1}" != "--update" ]]; then
	printf 'usage: %s --update REPO_DIR\n' "$0" >&2
	exit 1
fi

repo_dir="$(realpath "${2}")"

if [[ ! -d "${repo_dir}" ]]; then
	printf 'repo directory not found: %s\n' "${repo_dir}" >&2
	exit 1
fi

if ! podman image exists "${image_tag}"; then
	podman build -t "${image_tag}" "${context_dir}"
fi

exec podman run --rm \
	--userns=keep-id \
	--user "$(id -u):$(id -g)" \
	-v "${repo_dir}:/repo:Z" \
	"${image_tag}" \
	createrepo_c --update /repo