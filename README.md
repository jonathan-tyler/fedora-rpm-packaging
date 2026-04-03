# Fedora Package Repo Container

Serve the staged Fedora repository from a dedicated WSL service user with Podman, Quadlet, and Nginx.

This repo owns the repository-hosting side of the split. The builder and publisher logic stays in the sibling repo at `repos/fedora-local-builder/`.

## Service Model

The current recommendation is unchanged: run the HTTP repo service as a different WSL user instead of your normal build user.

- use a dedicated service account such as `reposvc`
- let that account own `~/fedora-package-repo/`
- sync staged repository content into that directory from the builder repo
- mount the live repo read-only into the container

That gives you a clean boundary between build activity and the long-running HTTP service.

## Repository Layout

```text
container/nginx/   Nginx container image and config
quadlet/           rootless Quadlet units for the repo service
scripts/           WSL setup and install scripts
```

## WSL Install

### 1. Host dependencies

Run this as root on the Fedora WSL host:

```bash
sudo dnf install podman systemd
```

### 2. Create the service user

Run this as root:

```bash
cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-package-repo-container
sudo bash scripts/setup-service-user.sh
```

Defaults:

- `REPO_SERVICE_USER=reposvc`
- live repo root at `~reposvc/fedora-package-repo/repo`

### 3. Install the Quadlet units and image

Run this as the service user:

```bash
sudo -iu reposvc bash -lc 'cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-package-repo-container && bash scripts/install-quadlet.sh'
sudo -iu reposvc systemctl --user enable --now fedora-package-repo.service
```

If you want a one-command bootstrap wrapper, use:

```bash
cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-package-repo-container
sudo bash scripts/install-wsl-service.sh
```

That wrapper installs Podman if needed, creates the service user, and runs the Quadlet install step as that user. It still leaves `systemctl --user enable --now ...` as an explicit final step.

## Live Repo Layout

```text
~reposvc/fedora-package-repo/
  repo/
    fedora/
    downloads/
    keys/
```

The builder repo syncs staged content into that `repo/` directory with `fpb repo sync-service`.

## Notes

- Podman is the validated runtime here.
- The builder repo still owns repository generation, signing, and sync.
- This repo only owns the long-running containerized HTTP service.