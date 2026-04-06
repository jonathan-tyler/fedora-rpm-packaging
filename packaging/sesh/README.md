# sesh RPM packaging

This directory holds the first RPM packaging scaffold for `sesh`.

Current flow:

1. Check out the exact upstream tag that matches `Version` in `sesh.spec`.
2. Run `go run ./src/cmd/fpb source vendor sesh /path/to/sesh-worktree`.
3. Run `bash packaging/sesh/prepare-sources.sh /path/to/sesh-worktree`.
4. Run `go run ./src/cmd/fpb build srpm sesh`.
5. Rebuild that SRPM with mock.

After step 3, `packaging/sesh/SOURCES/` contains:

- `sesh-2.24.2.tar.gz` created locally from the checked-out upstream source tree
- `vendor/` materialized from the locally generated vendored dependency archive
- `sesh-vendor.tar.gz` rebuilt locally from that `vendor/` directory

This flow is building from source. It does not download prebuilt release archives from GitHub.

The `build srpm` step runs in a one-shot Podman helper container so RPM-specific
tooling does not need to be installed directly on the host.