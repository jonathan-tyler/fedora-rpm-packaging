#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'usage: %s /path/to/ghostty-worktree\n' "$0" >&2
	exit 1
fi

worktree="$(realpath "$1")"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
spec_path="${script_dir}/ghostty.spec"
sources_dir="${script_dir}/SOURCES"
container_dir="${repo_root}/packages/ghostty"
container_tag="localhost/fedora-package-builder-ghostty-sourceprep:latest"
cache_archive="${sources_dir}/ghostty-zig-cache.tar.gz"

if [[ ! -d "${worktree}" ]]; then
	printf 'worktree not found: %s\n' "${worktree}" >&2
	exit 1
fi

version="$(sed -n 's/^Version:[[:space:]]*//p' "${spec_path}" | head -n 1)"
zig_version="$(sed -n 's/^%global[[:space:]]\+zig_version[[:space:]]\+//p' "${spec_path}" | head -n 1)"
if [[ -z "${version}" || -z "${zig_version}" ]]; then
	printf 'could not determine version metadata from: %s\n' "${spec_path}" >&2
	exit 1
fi
source_tarball_url="https://release.files.ghostty.org/${version}/ghostty-${version}.tar.gz"

expected_tag="v${version}"
actual_tags="$(git -C "${worktree}" tag --points-at HEAD || true)"
if ! printf '%s\n' "${actual_tags}" | grep -Fxq "${expected_tag}"; then
	found_tag="$(printf '%s\n' "${actual_tags}" | head -n 1)"
	printf 'worktree must be checked out at tag %s; found %s\n' "${expected_tag}" "${found_tag:-<none>}" >&2
	exit 1
fi

if [[ ! -d "${container_dir}" ]]; then
	printf 'container build context not found: %s\n' "${container_dir}" >&2
	exit 1
fi

mkdir -p "${sources_dir}"

curl -fsSL "${source_tarball_url}" \
	-o "${sources_dir}/ghostty-${version}.tar.gz"

curl -fsSL "https://ziglang.org/download/${zig_version}/zig-x86_64-linux-${zig_version}.tar.xz" \
	-o "${sources_dir}/zig-x86_64-linux-${zig_version}.tar.xz"

tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

podman build -t "${container_tag}" "${container_dir}" >/dev/null
podman run --rm --network host \
	-v "${worktree}:/src/ghostty:Z" \
	-v "${tmpdir}:/work:Z" \
	"${container_tag}" \
	bash -lc 'set -euo pipefail; cd /src/ghostty; export ZIG_GLOBAL_CACHE_DIR=/work/cache; ./nix/build-support/fetch-zig-cache.sh'

if [[ ! -d "${tmpdir}/cache/p" ]]; then
	printf 'ghostty zig cache was not created under %s\n' "${tmpdir}/cache/p" >&2
	exit 1
fi

tar -czf "${cache_archive}" -C "${tmpdir}/cache" p

printf 'prepared %s\n' "${sources_dir}/ghostty-${version}.tar.gz"
printf 'prepared %s\n' "${cache_archive}"
printf 'prepared %s\n' "${sources_dir}/zig-x86_64-linux-${zig_version}.tar.xz"