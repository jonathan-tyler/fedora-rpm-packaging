#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
image_tag="localhost/fedora-package-builder-rpm-packaging:latest"
context_dir="${repo_root}/tools/rpm-packaging"

if [[ $# -ne 1 ]]; then
	printf 'usage: %s PACKAGE\n' "$0" >&2
	exit 1
fi

package_name="$1"
packaging_dir="${repo_root}/packaging/${package_name}"
spec_path="${packaging_dir}/${package_name}.spec"
sources_dir="${packaging_dir}/SOURCES"

if [[ ! -d "${packaging_dir}" ]]; then
	printf 'packaging directory not found: %s\n' "${packaging_dir}" >&2
	exit 1
fi

if [[ ! -f "${spec_path}" ]]; then
	printf 'spec file not found: %s\n' "${spec_path}" >&2
	exit 1
fi

if [[ ! -d "${sources_dir}" ]]; then
	printf 'sources directory not found: %s\n' "${sources_dir}" >&2
	exit 1
fi

if ! find "${sources_dir}" -mindepth 1 -type f -print -quit | grep -q .; then
	printf 'no source files found under: %s\n' "${sources_dir}" >&2
	exit 1
fi

if ! podman image exists "${image_tag}"; then
	podman build -t "${image_tag}" "${context_dir}"
fi

exec podman run --rm \
	--userns=keep-id \
	--user "$(id -u):$(id -g)" \
	-v "${repo_root}:/workspace:Z" \
	-w /workspace \
	"${image_tag}" \
	/usr/local/bin/build-srpm "${package_name}"