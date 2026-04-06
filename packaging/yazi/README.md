# yazi RPM packaging

This directory holds the RPM packaging scaffold for `yazi`.

Steps:

1. Check out the exact upstream tag that matches `Version` in `yazi.spec`.
2. Run `go run ./src/cmd/fpb source vendor yazi /path/to/yazi-worktree`.
3. Run `bash packaging/yazi/prepare-sources.sh /path/to/yazi-worktree`.
4. Run `go run ./src/cmd/fpb build srpm yazi`.

After step 3, `packaging/yazi/SOURCES/` contains:

- `yazi-26.1.22.tar.gz` created locally from the checked-out upstream source tree
- `vendor/` materialized from locally vendored Cargo dependencies
- `cargo-config.toml` emitted by `cargo vendor`
- `yazi-vendor.tar.gz` rebuilt locally from `vendor/` and `cargo-config.toml`