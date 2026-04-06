# ghostty RPM packaging notes

This package uses the same source-first RPM flow as the other packaged tools,
but Ghostty needs two extra offline inputs:

- the upstream Ghostty release source tarball
- a prefetched Zig dependency cache
- the exact Zig toolchain tarball required by the pinned release

Prepare the sources from a pinned worktree:

```bash
go run ./src/cmd/fpb upstream fetch ghostty
go run ./src/cmd/fpb upstream worktree add ghostty v1.3.1
bash packaging/ghostty/prepare-sources.sh state/worktrees/ghostty-v1.3.1
```

Then build the SRPM and validate it in mock:

```bash
go run ./src/cmd/fpb build srpm ghostty
go run ./src/cmd/fpb build mock-rebuild ghostty state/results/ghostty/ghostty-1.3.1-1.fc42.src.rpm
```

The source prep step reuses the existing `packages/ghostty/Containerfile` image
with `podman run --network host` to populate the offline Zig cache. That keeps
the networked fetch stage containerized while still leaving the mock rebuild
self-contained.

Unlike the other packaged tools, Ghostty should use the upstream release source
tarball rather than a plain Git archive. Ghostty's own packaging guidance notes
that the release tarball contains preprocessed files that downstream packagers
should rely on.