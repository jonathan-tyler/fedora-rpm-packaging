# Fedora Package Builder

Build reproducible Fedora RPMs with mock, sign them with your own GPG key, and stage them for a separate repository-hosting service.

This repo implementation lives under `repos/fedora-local-builder/` inside the parent project.

The Nginx, Quadlet, and WSL repo-service setup assets now live in the sibling repo `repos/fedora-package-repo-container/`.

## Why This Project Exists

This project replaces the older build-yazi-only note with a broader package builder that can cover yazi and other tools from the same workflow.

The current target clients are Odin and Fenrir.

## Design Summary

- mock runs on the Fedora host instead of inside a rootless container
- package builds run separately as one-off jobs under the normal WSL user
- the build phase is offline by default
- source fetch and dependency vendoring are explicit pre-build steps
- packages and repository metadata are signed with your own GPG key
- the live repo is updated through a handoff step instead of being owned directly by the daily user

The long-running HTTP service is intentionally split out from this builder repo.

The implementation is now a Go CLI with a hexagonal layout. There is no script wrapper layer anymore.

External tool invocation is also centralized now. Command names and default extra arguments live in `config/tooling.json` instead of being hardcoded across the app layer.

The mock decision is intentional. Nested mock inside a rootless container is possible to chase, but it is the wrong default because mock is most reliable when it owns the Fedora build root directly on the host.

## Recommended Deployment Model

For your current setup, this is the right default:

- one dedicated WSL user such as reposvc owns and runs the HTTP repo service
- your normal WSL user does source fetch, vendoring, SRPM creation, and mock rebuilds on demand
- the build output is staged locally first
- a separate sync step copies the staged repository into the reposvc-owned live repo directory

That gives you the separation you want without pretending you need a full bare-metal package server.

## Security Posture

For a low-to-moderate home-lab threat model, I would recommend this order of priorities:

1. Separate the service identity from your daily user.
2. Keep the live repo content in a service-owned directory.
3. Mount the served repo read-only into the container.
4. Restrict Windows firewall ingress to the specific LAN clients you care about.
5. Keep the build path offline by default.

That is a better return on complexity than trying to fully containerize the build path.

If you later want a stricter no-egress guarantee for the HTTP container itself, add a host-owned socket plus internal-network proxy pattern. For now, that is probably more complexity than your threat model justifies.

## Target Packages

- sesh: Go, vendor with go mod vendor
- television: Rust Cargo build, stages `tv`
- yazi: Rust Cargo workspace build, stages `yazi` and `ya`

## Containerized Builds

The implemented containerized build targets are `sesh`, `television`, and `yazi`.

The flow is:

1. build a local Fedora builder image
2. clone the target repository inside the container
3. vendor dependencies while the container still has network access
4. commit that prepared container state to a temporary local image
5. run the build from that snapshot with the configured offline container runtime arguments
6. write the built binaries into the staged server tree under `state/repo/downloads/`

This keeps the source checkout out of your host filesystem while still giving you a host-visible artifact at the end.

By default, `config/tooling.json` sets Git clone extra arguments to `--depth 1` and the offline container runtime extra arguments to `--network none`.

The container runtime defaults to Podman. There is also a Docker shim that reuses the Podman-compatible argument model, but it is only a placeholder for future work and is not validated or supported yet.

Run it with:

```bash
go run ./src/cmd/fpb build container-binary sesh
go run ./src/cmd/fpb build container-binary television
go run ./src/cmd/fpb build container-binary yazi
```

Or pin a tag or branch:

```bash
go run ./src/cmd/fpb build container-binary sesh v2.15.0
```

The artifact lands here:

```text
state/repo/downloads/PACKAGE/TIMESTAMP/
   sesh|tv|yazi|ya
   build-info.txt
```

For `yazi`, the staged directory includes both `yazi` and `ya`.

Then publish it to the live service with the existing sync step:

```bash
sudo -iu reposvc bash -lc 'cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-local-builder && go run ./src/cmd/fpb repo sync-service'
```

At that point the binary is available from the repo server as a normal static download once the sibling service repo is installed.

## Important Distinction

This container flow produces raw staged binary downloads, not RPMs.

That is intentional. It is the lowest-risk way to prove the online-fetch then offline-build boundary before we invest in a full RPM recipe for the same package.

## Directory Layout

```text
config/                 project-local configuration and package registry
container/build/        ephemeral package build images
clients/                example DNF repo files for client machines
mock/                   project-local mock overrides
packaging/              spec files, source tarballs, and packaging notes
state/                  generated mirrors, build outputs, and staged repo data
src/cmd/fpb/            main CLI entrypoint
src/commands/           CLI adapters and command routing
src/app/                application use cases
src/core/               domain models and ports
src/infra/              filesystem, config, archive, command, and external-tool adapters
```

The package registry lives in `config/packages.json` instead of shell functions.

External tool names and default extra arguments live in `config/tooling.json`.

Repository hosting assets live in the sibling repo `repos/fedora-package-repo-container/`.

## Service User Layout

The recommended live service layout is:

```text
~reposvc/fedora-package-repo/
  repo/
    fedora/
    keys/
```

The repo container mounts that directory read-only. Your normal user should not be the owner of that live path.

## Build And Handoff Model

Use this repo tree for connected preparation and local staging:

- fetch upstreams into state/upstream/
- create vendored source bundles under state/vendor/
- run mock rebuilds into state/results/
- publish the signed staged repo into state/repo/
- stage raw binary downloads into state/repo/downloads/

Then import that staged repo into the service-owned live repo with `fpb repo sync-service`.

That handoff is the right place to keep the trust boundary simple.

## Build Model

1. Fetch upstream Git repositories while connected.
2. Vendor language dependencies while connected.
3. Freeze source inputs into packaging/.
4. Build SRPMs.
5. Rebuild inside mock with networking disabled.
6. Sign RPMs and repository metadata into the staged repo.
7. Sync the staged repo into the service-owned live repo.
8. Serve the live repo with Nginx from a rootless Podman container.

That split keeps the actual build step deterministic and mostly offline while still acknowledging that upstream source acquisition has to happen somewhere.

## Quick Start

1. Install host dependencies.

   ```bash
   sudo dnf install golang mock rpm-sign gnupg2 createrepo_c git podman
   ```

2. Create a dedicated service user.

   ```bash
   sudo useradd -m -s /usr/sbin/nologin reposvc
   sudo loginctl enable-linger reposvc
   ```

3. Configure RPM signing for your normal build user.

   ```bash
   cp config/rpmmacros.example ~/.rpmmacros
   ```

   Replace the placeholder key ID before signing anything.

4. Mirror the upstream repositories.

   ```bash
   go run ./src/cmd/fpb upstream fetch
   ```

5. Prepare vendored sources from a checked-out worktree.

   ```bash
   go run ./src/cmd/fpb source vendor yazi /path/to/yazi-worktree
   ```

6. Build an SRPM from your spec tree, then rebuild it in mock.

   ```bash
   go run ./src/cmd/fpb build mock-rebuild yazi /path/to/yazi.src.rpm
   ```

7. Publish the resulting RPMs into the staged repository.

   ```bash
   FEDORA_PACKAGE_GPG_KEY="YOUR KEY ID" go run ./src/cmd/fpb repo publish yazi
   ```

8. Sync the staged repository into the live repo as the dedicated service user.

   ```bash
   sudo -iu reposvc bash -lc 'cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-local-builder && go run ./src/cmd/fpb repo sync-service'
   ```

9. Install the sibling repo-service container as the service user.

   ```bash
   sudo -iu reposvc bash -lc 'cd /home/him/my/forge/home/2026-03-mouse-free-migration/2026-04-fedora-package-builder/repos/fedora-package-repo-container && bash scripts/install-quadlet.sh'
   sudo -iu reposvc systemctl --user enable --now fedora-package-repo.service
   ```

10. Copy clients/fedora-local-builder.repo.example onto Odin and Fenrir, then replace REPO_HOST with the builder host name or LAN IP.

## Command Surface

The Go CLI is organized as small command adapters over focused use-case packages:

```text
fpb package list
fpb upstream fetch [PACKAGE...]
fpb source vendor PACKAGE /path/to/worktree
fpb build container-binary PACKAGE [REF]
fpb build mock-rebuild PACKAGE path/to/package.src.rpm
fpb repo publish PACKAGE
fpb repo sync-service
```

Use `go run ./src/cmd/fpb help` to print the available commands.

## Network Policy

The intended steady state is:

- no internet access during mock rebuilds
- no internet requirement for the live repo container
- Windows firewall restricts ingress to the machines you choose
- only the explicit fetch and vendoring step may touch the internet

The sibling repo-service scaffold is the simple mode: rootless container plus published LAN port. It is good enough for a low-threat environment.

If you later decide the container must have a hard no-egress boundary instead of a practical one, switch the service to a host-owned socket plus an internal Podman network.

If a package build breaks without the network, treat that as a packaging gap and add the missing vendored inputs. Do not relax the build policy unless the package genuinely cannot be packaged that way.

## Serving The Repo

The repo server publishes HTTP on port 8080 by default. A client repo file can point DNF to:

```text
http://REPO_HOST:8080/fedora/42/x86_64/
```

If you later want TLS, put a reverse proxy or LAN certificate in front of this service rather than complicating the first iteration.

## Recommendations

- Yes, use a dedicated WSL service account for the repo server.
- Yes, keep builds separate and run them on demand as your normal user.
- No, I would not add a separate long-running build container right now.
- Yes, use a service-owned live repo directory plus a sync or handoff step.
- Yes, use the Windows firewall to allow only the specific LAN clients you trust.

## Open Questions

- Should television stay in scope if its offline packaging story is worse than sesh and yazi?
- Should source acquisition live on the same host, or move to a separate connected machine later?
- Do you want one shared repo for Odin and Fenrir, or a per-host repo branch for risky packages?
