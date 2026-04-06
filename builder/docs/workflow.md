# Builder Workflow

This repo owns the local build, signing, and staging side of the workflow.

## What It Does

- mirrors upstream sources locally
- vendors dependencies before offline rebuilds
- builds SRPMs and mock rebuilds through one-shot helper containers
- signs RPMs and repository metadata
- stages a DNF repository under `state/repo/`

## Typical Flow

1. Mirror upstream sources.

   ```bash
   go run ./src/cmd/fpb upstream fetch
   ```

2. Vendor dependencies from a checked-out worktree when the package needs it.

   ```bash
   go run ./src/cmd/fpb source vendor PACKAGE /path/to/worktree
   ```

3. Freeze source inputs into `packaging/<name>/`.

   ```bash
   bash packaging/PACKAGE/prepare-sources.sh /path/to/worktree
   ```

4. Build an SRPM.

   ```bash
   go run ./src/cmd/fpb build srpm PACKAGE
   ```

5. Rebuild the SRPM with mock.

   ```bash
   go run ./src/cmd/fpb build mock-rebuild PACKAGE path/to/package.src.rpm
   ```

6. Initialize or refresh staged repository metadata.

   ```bash
   export FEDORA_PACKAGE_GPG_KEY="YOUR KEY ID OR FINGERPRINT"
   go run ./src/cmd/fpb repo init
   ```

7. Publish a package into the staged repository.

   ```bash
   export GPG_TTY="$(tty)"
   export FEDORA_PACKAGE_GPG_KEY="YOUR KEY ID OR FINGERPRINT"
   go run ./src/cmd/fpb repo publish PACKAGE
   ```

8. Import `state/repo/` into the live service tree with the companion service repo.

## Command Surface

```text
fpb package list
fpb upstream fetch [PACKAGE...]
fpb source vendor PACKAGE /path/to/worktree
fpb build container-binary PACKAGE [REF]
fpb build srpm PACKAGE
fpb build mock-rebuild PACKAGE path/to/package.src.rpm
fpb repo init
fpb repo publish PACKAGE
fpb repo sync-service
```

## Key Directories

- `packages/` holds per-package metadata and helper-container build inputs
- `packaging/` holds spec files and source-prep assets
- `mock/` holds repo-local mock overrides
- `tools/` holds shared helper-container assets
- `state/` holds generated mirrors, results, artifacts, and the staged repo tree

The sibling service repo is responsible for serving the staged repository over HTTP.