# Backend Boilerplate

A Go REST API boilerplate with JWT authentication, PostgreSQL, Ent ORM, and OpenAPI documentation.

## Tech Stack

- **Go 1.25** – Runtime
- **Chi** – HTTP router
- **Huma v2** – OpenAPI-first REST framework with automatic docs
- **Ent** – ORM with schema migrations via Atlas
- **PostgreSQL 16** – Database
- **JWT** – Authentication (JWK-based signing)
- **Reflex** – Hot reload during development

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [Task](https://taskfile.dev/) (optional, for running tasks)

## Quick Start

1. **Start PostgreSQL**

   ```bash
   task local:docker:up
   # or: docker compose up -d
   ```

2. **Apply migrations**

   ```bash
   task migrate:apply
   ```

3. **Run the API (with hot reload)**

   ```bash
   task local:api:watch
   # or: go run ./apps/api
   ```

4. **Open the docs**
   - Swagger UI: http://localhost:5000/docs
   - OpenAPI JSON: http://localhost:5000/openapi.json

## Environment Variables

Create a `.env` file (copy from `.env.example`) to override defaults. If `.env` is present, its values take precedence over Taskfile defaults.

| Variable            | Default     | Description                 |
| ------------------- | ----------- | --------------------------- |
| `APP_NAME`          | TEARest     | Application name in OpenAPI |
| `PORT`              | 5000        | HTTP server port            |
| `POSTGRES_USER`     | postgres    | PostgreSQL user             |
| `POSTGRES_PASSWORD` | postgres    | PostgreSQL password         |
| `POSTGRES_DB`       | tearest     | Database name               |
| `JWT_MASTER_KEY`    | (32 chars)  | JWT signing key (32 bytes)  |
| `GO_ENV`            | development | Environment mode            |

## Project Structure

```
├── apps/api/              # API entry point & controllers
│   ├── controller/        # HTTP handlers (health, auth, protected)
│   └── rest/              # Router, middleware, OpenAPI setup
├── internal/
│   ├── actions/           # JWT, JWK, auth logic
│   ├── db/                # Database client
│   ├── ent/               # Ent schema & generated code
│   │   ├── schema/        # Entity definitions
│   │   └── migrate/       # Atlas migrations
│   ├── filters/           # Query filters
│   ├── usecases/          # Business logic
│   └── validation/        # Custom validators
├── types/                 # Shared DTOs & types
├── constants/             # App constants
└── dist/                  # Generated OpenAPI clients (gitignored)
```

## API Endpoints

| Method | Path               | Auth               | Description                   |
| ------ | ------------------ | ------------------ | ----------------------------- |
| GET    | `/health`          | —                  | Liveness probe                |
| POST   | `/auth/login`      | —                  | Login (identifier + password) |
| POST   | `/auth/register`   | —                  | Register new user             |
| GET    | `/protected/ping`  | —                  | Public ping                   |
| GET    | `/protected/me`    | Bearer             | Current user (authenticated)  |
| GET    | `/protected/admin` | Bearer + SuperUser | Admin-only endpoint           |

## Tasks (Taskfile)

| Task                       | Description                                 |
| -------------------------- | ------------------------------------------- |
| `local:docker:up`          | Start PostgreSQL 16 container               |
| `local:api:watch`          | Run API with hot reload (reflex)            |
| `migrate:diff`             | Generate migration from schema changes      |
| `migrate:diff:go`          | Same, using Go script (no Atlas CLI)        |
| `migrate:apply`            | Apply pending migrations                    |
| `migrate:status`           | Show migration status                       |
| `ent:generate`             | Regenerate Ent code from schema             |
| `openapi:download`         | Download OpenAPI spec (API must be running) |
| `openapi:generate-clients` | Generate TypeScript client from spec        |
| `openapi`                  | Download spec + generate clients            |

## Generating Migrations

After changing Ent schemas in `internal/ent/schema/`:

```bash
# Option 1: Atlas CLI (uses Docker for dev DB)
task migrate:diff -- my_migration_name

# Option 2: Go script (Postgres must be running)
task migrate:diff:go -- my_migration_name
```

Then apply:

```bash
task migrate:apply
```

## Generating API Clients

With the API running:

```bash
task openapi
```

This downloads the OpenAPI spec and generates a TypeScript Axios client in `dist/clients/`.

## License

MIT
