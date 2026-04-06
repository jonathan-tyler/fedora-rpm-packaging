#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
  printf 'usage: %s /path/to/eza-worktree\n' "$0" >&2
  exit 1
fi

worktree="$(realpath "$1")"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
spec_path="${script_dir}/eza.spec"
sources_dir="${script_dir}/SOURCES"
vendor_archive_source="${repo_root}/state/vendor/eza/eza-vendor.tar.gz"
vendor_dir="${sources_dir}/vendor"
vendor_config="${sources_dir}/cargo-config.toml"
vendor_archive_target="${sources_dir}/eza-vendor.tar.gz"

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
actual_tags="$(git -C "${worktree}" tag --points-at HEAD || true)"
if ! printf '%s\n' "${actual_tags}" | grep -Fxq "${expected_tag}"; then
  found_tag="$(printf '%s\n' "${actual_tags}" | head -n 1)"
  printf 'worktree must be checked out at tag %s; found %s\n' "${expected_tag}" "${found_tag:-<none>}" >&2
  exit 1
fi

mkdir -p "${sources_dir}"
rm -rf "${vendor_dir}" "${vendor_config}"

git -C "${worktree}" archive \
  --format=tar.gz \
  --prefix="eza-${version}/" \
  -o "${sources_dir}/eza-${version}.tar.gz" \
  HEAD

tar -xzf "${vendor_archive_source}" -C "${sources_dir}"
tar -czf "${vendor_archive_target}" -C "${sources_dir}" vendor cargo-config.toml

printf 'prepared %s\n' "${sources_dir}/eza-${version}.tar.gz"
printf 'prepared %s\n' "${vendor_dir}"
printf 'prepared %s\n' "${vendor_config}"
printf 'prepared %s\n' "${vendor_archive_target}"