# Service Install

This repo owns the long-running HTTP service that exposes a staged Fedora package repository.

## Bootstrap

Run the bootstrap script as root on the Fedora host:

```bash
sudo bash scripts/bootstrap-fedora-service.sh
```

That script:

- installs Podman and systemd if needed
- creates the service user when it does not already exist
- copies this repo into a service-owned install path
- installs the Quadlet units
- starts the generated user service

## Default Paths And Overrides

- `REPO_SERVICE_USER` defaults to `reposvc`
- `REPO_SERVICE_ROOT` defaults to `~SERVICE_USER/fedora-package-repo`
- `REPO_SERVICE_INSTALL_ROOT` defaults to `~SERVICE_USER/fedora-package-repo-service`
- `LIVE_REPO_ROOT` defaults to `~SERVICE_USER/fedora-package-repo/repo`
- `QUADLET_TARGET_HOME` defaults to the target user's home directory

## Importing The Staged Repo

Import a staged repo tree from the builder side with:

```bash
sudo bash scripts/import-staged-repo.sh /absolute/path/to/staged/repo
```

By default the import excludes `artifacts/` so the public HTTP surface stays RPM-focused. Set `INCLUDE_PUBLIC_ARTIFACTS=1` only if you intentionally want to expose raw build outputs.

## Quadlet Install

For an already prepared service account, install or refresh the user units with:

```bash
bash scripts/install-quadlet.sh
systemctl --user start fedora-package-repo.service
```

The service publishes host port `8090` by default.

## WSL Notes

The generic Fedora bootstrap is the primary path. For WSL-specific convenience, use:

```bash
sudo bash scripts/install-wsl-service.sh
```

If mirrored networking needs a Windows-side Hyper-V firewall rule, the helper script `scripts/allow-mirrored-client.ps1` can create one for a chosen remote address.