# Fedora Package Builder

Build signed Fedora RPMs locally and stage a DNF repository for a separate publishing host or service.

## Requirements

- Fedora-compatible build host with Go, Podman, Git, GPG, and rpmsign
- A configured RPM signing identity for repository metadata and RPM signatures
- A companion service repo or equivalent HTTP host for the staged repository

## Quick Start

```bash
cp config/rpmmacros.example ~/.rpmmacros
go run ./src/cmd/fpb help
```

Follow [docs/workflow.md](docs/workflow.md) for the build, signing, and publish flow.

## Key Paths

- `packages/` defines package manifests and helper-container inputs
- `packaging/` holds spec files and source-preparation assets
- `state/` is generated local state and is not intended for Git
- `clients/` contains an example DNF repo file for client machines

## Disclaimer

This project is built for personal use, experimentation, and learning. If you choose to use it in a production environment, you are responsible for validating and operating it safely.
