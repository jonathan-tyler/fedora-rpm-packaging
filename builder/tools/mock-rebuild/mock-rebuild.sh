#!/usr/bin/env bash

set -euo pipefail

usage() {
	printf 'usage: %s --workspace-root DIR --config-source DIR --root ROOT --resultdir DIR --srpm PATH\n' "$0" >&2
	exit 1
}

workspace_root=''
config_source=''
root_name=''
result_dir=''
srpm_path=''

while [[ $# -gt 0 ]]; do
	case "$1" in
		--workspace-root)
			[[ $# -ge 2 ]] || usage
			workspace_root="$2"
			shift 2
			;;
		--config-source)
			[[ $# -ge 2 ]] || usage
			config_source="$2"
			shift 2
			;;
		--root)
			[[ $# -ge 2 ]] || usage
			root_name="$2"
			shift 2
			;;
		--resultdir)
			[[ $# -ge 2 ]] || usage
			result_dir="$2"
			shift 2
			;;
		--srpm)
			[[ $# -ge 2 ]] || usage
			srpm_path="$2"
			shift 2
			;;
		*)
			usage
			;;
	esac
done

[[ -n "$workspace_root" && -n "$config_source" && -n "$root_name" && -n "$result_dir" && -n "$srpm_path" ]] || usage
[[ -d "$config_source" ]] || { printf 'config source not found: %s\n' "$config_source" >&2; exit 1; }
[[ -f "$srpm_path" ]] || { printf 'SRPM not found: %s\n' "$srpm_path" >&2; exit 1; }

temp_config_dir="$(mktemp -d)"
cleanup() {
	rm -rf "$temp_config_dir"
}
trap cleanup EXIT

cp -a /etc/mock/. "$temp_config_dir/"
cp -a "$config_source/." "$temp_config_dir/"

root_config_path="${temp_config_dir}/${root_name}.cfg"
[[ -f "$root_config_path" ]] || { printf 'mock root config not found: %s\n' "$root_config_path" >&2; exit 1; }

project_root_pattern="os.environ['FEDORA_PACKAGE_BUILDER_ROOT']"
if grep -Fq "$project_root_pattern" "$root_config_path"; then
	sed -i "s|os\.environ\['FEDORA_PACKAGE_BUILDER_ROOT'\]|'$workspace_root'|g" "$root_config_path"
fi

if grep -Fq "$project_root_pattern" "$root_config_path"; then
	printf 'failed to rewrite project root in mock config: %s\n' "$root_config_path" >&2
	exit 1
fi

mkdir -p "$result_dir"

exec mock \
	--configdir "$temp_config_dir" \
	--root "$root_name" \
	--resultdir "$result_dir" \
	--rebuild "$srpm_path"