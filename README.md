# Fedora Package Infra

Monorepo for the Fedora package builder and the Fedora package repo service.

## Layout

- `builder/` contains the Go CLI, packaging inputs, helper containers, and staged repo workflow
- `repo-service/` contains the containerized HTTP service, Quadlet units, and service bootstrap scripts

## Quick Start

```bash
cd builder
go test ./...
```

Use [builder/README.md](builder/README.md) for builder workflow details and [repo-service/README.md](repo-service/README.md) for service installation.

## Disclaimer

This project is built for personal use, experimentation, and learning. If you choose to use it in a production environment, you are responsible for validating and operating it safely.
