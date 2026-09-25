# sokosplit-core-api

**Status: early development**

The business API for [SokoSplit](https://github.com/sokosplit) — written in Go with Gin. Owns users, split-list definitions, and transaction history; calls into `sokosplit-wallet-service` for anything on-chain.

## How it fits into SokoSplit

- [sokosplit-contracts](https://github.com/sokosplit/sokosplit-contracts) — the on-chain escrow/split logic
- [sokosplit-wallet-service](https://github.com/sokosplit/sokosplit-wallet-service) — Stellar-facing service this API calls
- **sokosplit-core-api** (this repo) — the business layer
- [sokosplit-sdk](https://github.com/sokosplit/sokosplit-sdk) — client library + CLI that wraps this API
- [sokosplit-web](https://github.com/sokosplit/sokosplit-web) — dashboard, calls this API only

## Endpoints

See the full OpenAPI spec at [`/docs`](./docs/openapi.yaml) once running, or the file directly.

| Method | Path | Purpose |
|---|---|---|
| POST | `/api/v1/split-lists` | Create a split list |
| GET | `/api/v1/split-lists` | List the authenticated user's split lists |
| GET | `/api/v1/split-lists/:id` | Get a split list |
| POST | `/api/v1/split-lists/:id/release` | Request release of an escrowed split |
| POST | `/webhooks/wallet-service` | Wallet service callback — updates status |

## Setup

Requires Go 1.22+ and PostgreSQL.

```bash
cp .env.example .env   # fill in DATABASE_URL, JWT_SECRET, WALLET_SERVICE_URL
go mod download
go run main.go
```

## Testing

```bash
go test ./... -v
```

Handler tests run against an in-memory SQLite DB (no Postgres needed) and a fixed `user_id` in place of real JWT auth, so they stay fast and hermetic.

## Status

Skeleton in place: Gin app, GORM models (`User`, `SplitList`, `Recipient`), a CRUD flow for split lists, JWT auth middleware, an event-bus stub that calls `sokosplit-wallet-service`'s webhook, and a minimal OpenAPI spec. GitHub OAuth login is not yet wired up — `AuthMiddleware` expects a JWT already minted elsewhere.

## License

MIT
