# Contributing to sokosplit-core-api

## Getting started

1. Fork and clone the repo.
2. Install Go 1.22+ and PostgreSQL (or rely on the in-memory SQLite tests, which need neither).
3. `cp .env.example .env` and fill in values.
4. `go mod download && go test ./... -v` to confirm your environment is set up correctly.

## Making a change

- Keep on-chain calls isolated behind `sokosplit-wallet-service` — this repo should never talk to Stellar directly.
- Add or update tests in the relevant `internal/*` package for any behavior change.
- Run `go vet ./...` and `go test ./...` before opening a PR.
- Open a PR against `main`. At least one review is required before merge.
- If you change request/response shapes, update `docs/openapi.yaml` in the same PR.

## Reporting issues

Open a GitHub issue with steps to reproduce, or a proposed change and rationale for a feature request.
