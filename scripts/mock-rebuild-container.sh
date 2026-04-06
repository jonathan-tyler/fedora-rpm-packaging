#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
image_tag="localhost/fedora-package-builder-mock-rebuild:latest"
context_dir="${repo_root}/tools/mock-rebuild"

usage() {
	printf 'usage: %s --configdir DIR --root ROOT --resultdir DIR --rebuild PATH\n' "$0" >&2
	exit 1
}

abs_existing_path() {
	local input_path="$1"
	local parent_dir
	parent_dir="$(cd "$(dirname "$input_path")" && pwd)"
	printf '%s/%s\n' "$parent_dir" "$(basename "$input_path")"
}

abs_dir_path() {
	local input_path="$1"
	mkdir -p "$input_path"
	(cd "$input_path" && pwd)
}

map_into_workspace() {
	local host_path="$1"
	case "$host_path" in
		"$repo_root")
			printf '/workspace\n'
			;;
		"$repo_root"/*)
			printf '/workspace%s\n' "${host_path#$repo_root}"
			;;
		*)
			printf 'path must live under the repo root: %s\n' "$host_path" >&2
			exit 1
			;;
	esac
}

config_dir=''
root_name=''
result_dir=''
srpm_path=''

while [[ $# -gt 0 ]]; do
	case "$1" in
		--configdir)
			[[ $# -ge 2 ]] || usage
			config_dir="$2"
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
		--rebuild)
			[[ $# -ge 2 ]] || usage
			srpm_path="$2"
			shift 2
			;;
		*)
			usage
			;;
	esac
done

[[ -n "$config_dir" && -n "$root_name" && -n "$result_dir" && -n "$srpm_path" ]] || usage
[[ -d "$config_dir" ]] || { printf 'config directory not found: %s\n' "$config_dir" >&2; exit 1; }
[[ -f "$srpm_path" ]] || { printf 'SRPM not found: %s\n' "$srpm_path" >&2; exit 1; }

config_dir="$(abs_existing_path "$config_dir")"
result_dir="$(abs_dir_path "$result_dir")"
srpm_path="$(abs_existing_path "$srpm_path")"

container_config_dir="$(map_into_workspace "$config_dir")"
container_result_dir="$(map_into_workspace "$result_dir")"
container_srpm_path="$(map_into_workspace "$srpm_path")"

if ! podman image exists "$image_tag"; then
	podman build -t "$image_tag" "$context_dir"
fi

exec podman run --rm \
	--privileged \
	-v "${repo_root}:/workspace:Z" \
	-w /workspace \
	"$image_tag" \
	/usr/local/bin/mock-rebuild-helper \
	--workspace-root /workspace \
	--config-source "$container_config_dir" \
	--root "$root_name" \
	--resultdir "$container_result_dir" \
	--srpm "$container_srpm_path"