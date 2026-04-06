#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'usage: %s /path/to/lazygit-worktree\n' "$0" >&2
	exit 1
fi

worktree="$(realpath "$1")"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
spec_path="${script_dir}/lazygit.spec"
sources_dir="${script_dir}/SOURCES"

if [[ ! -d "${worktree}" ]]; then
	printf 'worktree not found: %s\n' "${worktree}" >&2
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

git -C "${worktree}" archive \
	--format=tar.gz \
	--prefix="lazygit-${version}/" \
	-o "${sources_dir}/lazygit-${version}.tar.gz" \
	HEAD

printf 'prepared %s\n' "${sources_dir}/lazygit-${version}.tar.gz"