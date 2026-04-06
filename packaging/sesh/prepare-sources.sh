#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'usage: %s /path/to/sesh-worktree\n' "$0" >&2
	exit 1
fi

worktree="$(realpath "$1")"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
spec_path="${script_dir}/sesh.spec"
sources_dir="${script_dir}/SOURCES"
vendor_archive_source="${repo_root}/state/vendor/sesh/sesh-vendor.tar.gz"
vendor_dir="${sources_dir}/vendor"
vendor_archive_target="${sources_dir}/sesh-vendor.tar.gz"

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

expected_tag="v${version}"
actual_tag="$(git -C "${worktree}" describe --tags --exact-match 2>/dev/null || true)"
if [[ "${actual_tag}" != "${expected_tag}" ]]; then
	printf 'worktree must be checked out at tag %s; found %s\n' "${expected_tag}" "${actual_tag:-<none>}" >&2
	exit 1
fi

mkdir -p "${sources_dir}"
rm -rf "${vendor_dir}"

git -C "${worktree}" archive \
	--format=tar.gz \
	--prefix="sesh-${version}/" \
	-o "${sources_dir}/sesh-${version}.tar.gz" \
	HEAD

mkdir -p "${vendor_dir}"
tar -xzf "${vendor_archive_source}" -C "${sources_dir}"
tar -czf "${vendor_archive_target}" -C "${sources_dir}" vendor

printf 'prepared %s\n' "${sources_dir}/sesh-${version}.tar.gz"
printf 'prepared %s\n' "${vendor_dir}"
printf 'prepared %s\n' "${vendor_archive_target}"