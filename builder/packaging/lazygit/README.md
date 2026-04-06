# lazygit RPM packaging

This directory holds the RPM packaging scaffold for `lazygit`.

Current flow:

1. Check out the exact upstream tag that matches `Version` in `lazygit.spec`.
2. Run `bash packaging/lazygit/prepare-sources.sh /path/to/lazygit-worktree`.
3. Run `go run ./src/cmd/fpb build srpm lazygit`.
4. Rebuild that SRPM with mock.

After step 2, `packaging/lazygit/SOURCES/` contains:

- `lazygit-0.60.0.tar.gz` created locally from the checked-out upstream source tree

This flow is building from source. It does not download prebuilt release archives from GitHub.

`lazygit` already commits its `vendor/` tree upstream, so this recipe packages the tagged source snapshot directly and does not need a separate `fpb source vendor` step.

The `build srpm` step runs in a one-shot Podman helper container so RPM-specific
tooling does not need to be installed directly on the host.