#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'usage: %s /path/to/television-worktree\n' "$0" >&2
	exit 1
fi

worktree="$(realpath "$1")"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
spec_path="${script_dir}/television.spec"
sources_dir="${script_dir}/SOURCES"
vendor_archive_source="${repo_root}/state/vendor/television/television-vendor.tar.gz"
vendor_dir="${sources_dir}/vendor"
vendor_config="${sources_dir}/cargo-config.toml"
vendor_archive_target="${sources_dir}/television-vendor.tar.gz"

if [[ ! -d "${worktree}" ]]; then
	printf 'worktree not found: %s\n' "${worktree}" >&2
	exit 1
fi

if [[ ! -f "${vendor_archive_source}" ]]; then
	printf 'vendor archive not found: %s\n' "${vendor_archive_source}" >&2
	exit 1
fi

version="$(sed -n 's/^Version:[[:space:]]*//p' "${spec_path}" | head -n 1)"
if [[ -z "${version}" ]]; then
	printf 'could not determine version from: %s\n' "${spec_path}" >&2
	exit 1
fi

actual_tag="$(git -C "${worktree}" describe --tags --exact-match 2>/dev/null || true)"
if [[ "${actual_tag}" != "${version}" ]]; then
	printf 'worktree must be checked out at tag %s; found %s\n' "${version}" "${actual_tag:-<none>}" >&2
	exit 1
fi

mkdir -p "${sources_dir}"
rm -rf "${vendor_dir}" "${vendor_config}"

git -C "${worktree}" archive \
	--format=tar.gz \
	--prefix="television-${version}/" \
	-o "${sources_dir}/television-${version}.tar.gz" \
	HEAD

tar -xzf "${vendor_archive_source}" -C "${sources_dir}"
tar -czf "${vendor_archive_target}" -C "${sources_dir}" vendor cargo-config.toml

printf 'prepared %s\n' "${sources_dir}/television-${version}.tar.gz"
printf 'prepared %s\n' "${vendor_dir}"
printf 'prepared %s\n' "${vendor_config}"
printf 'prepared %s\n' "${vendor_archive_target}"