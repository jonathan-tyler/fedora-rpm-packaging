# television RPM packaging

This directory holds the RPM packaging scaffold for `television`.

Steps:

1. Check out the exact upstream tag that matches `Version` in `television.spec`.
2. Run `go run ./src/cmd/fpb source vendor television /path/to/television-worktree`.
3. Run `bash packaging/television/prepare-sources.sh /path/to/television-worktree`.
4. Run `go run ./src/cmd/fpb build srpm television`.

After step 3, `packaging/television/SOURCES/` contains:

- `television-0.15.4.tar.gz` created locally from the checked-out upstream source tree
- `vendor/` materialized from locally vendored Cargo dependencies
- `cargo-config.toml` emitted by `cargo vendor`
- `television-vendor.tar.gz` rebuilt locally from `vendor/` and `cargo-config.toml`