#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
	printf 'usage: %s PACKAGE\n' "$0" >&2
	exit 1
fi

package_name="$1"
repo_root="${FPB_REPO_ROOT:-/workspace}"
packaging_dir="${repo_root}/packaging/${package_name}"
spec_path="${packaging_dir}/${package_name}.spec"
sources_dir="${packaging_dir}/SOURCES"
result_dir="${repo_root}/state/results/${package_name}"

if [[ ! -f "${spec_path}" ]]; then
	printf 'spec file not found: %s\n' "${spec_path}" >&2
	exit 1
fi

if [[ ! -d "${sources_dir}" ]]; then
	printf 'sources directory not found: %s\n' "${sources_dir}" >&2
	exit 1
fi

mkdir -p "${result_dir}"

topdir="$(mktemp -d)"
trap 'rm -rf "${topdir}"' EXIT

mkdir -p "${topdir}/BUILD" "${topdir}/BUILDROOT" "${topdir}/RPMS" "${topdir}/SOURCES" "${topdir}/SPECS" "${topdir}/SRPMS"

cp "${spec_path}" "${topdir}/SPECS/${package_name}.spec"
cp -a "${sources_dir}/." "${topdir}/SOURCES/"

rpmbuild \
	--define "_topdir ${topdir}" \
	--define "_srcrpmdir ${result_dir}" \
	-bs "${topdir}/SPECS/${package_name}.spec"