# eza RPM packaging

This directory holds the RPM packaging scaffold for `eza`.

Steps:

1. Check out the exact upstream tag that matches `Version` in `eza.spec`.
2. Run `go run ./src/cmd/fpb source vendor eza /path/to/eza-worktree`.
3. Run `bash packaging/eza/prepare-sources.sh /path/to/eza-worktree`.
4. Run `go run ./src/cmd/fpb build srpm eza`.

After step 3, `packaging/eza/SOURCES/` contains:

- `eza-0.23.4.tar.gz` created locally from the checked-out upstream source tree
- `vendor/` materialized from locally vendored Cargo dependencies
- `cargo-config.toml` emitted by `cargo vendor`
- `eza-vendor.tar.gz` rebuilt locally from `vendor/` and `cargo-config.toml`