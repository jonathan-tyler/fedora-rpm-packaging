# Fedora Package Repo Service

Serve a staged Fedora package repository with rootless Podman, Quadlet, and Nginx.

## Requirements

- Fedora host with Podman and systemd user services
- Root access for bootstrap and staged-repo imports
- A staged repository tree from the builder repo or an equivalent source

## Quick Start

```bash
sudo bash scripts/bootstrap-fedora-service.sh
sudo bash scripts/import-staged-repo.sh /absolute/path/to/staged/repo
```

Follow [docs/service-install.md](docs/service-install.md) for service-user overrides, Quadlet installation details, and WSL notes.

## Key Paths

- The live repo defaults to `~reposvc/fedora-package-repo/repo`
- The service-owned repo copy defaults to `~reposvc/fedora-package-repo-service`
- Quadlet units are installed under `~/.config/containers/systemd`