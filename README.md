# Chirpy

Chirpy is a minimal Go microblogging API (chirps) backed by PostgreSQL. It demonstrates:
- A small REST API with endpoints for chirps and users
- Authentication with JWT + refresh tokens
- Password hashing (argon2id)
- SQL-first DB access via sqlc

## Features
- CRUD-ish endpoints for chirps (list, create, get by id, delete)
- User registration and update
- Login / refresh / revoke token flows
- Admin endpoints: reset DB (/admin/reset) and metrics (/admin/metrics)
- Static files served under /app/

## Tech stack
- Go 1.26.x
- PostgreSQL
- sqlc for type-safe DB access (internal/database)
- JWT (github.com/golang-jwt/jwt/v5)
- Argon2id for password hashing (github.com/alexedwards/argon2id)

## Quick start (local)

Prerequisites:
- Go 1.26.x
- PostgreSQL
- (optional) sqlc if you want to regenerate DB code

#### 1. Clone Repo
```bash
git clone https://github.com/jzmack/chirpy.git
cd chirpy
```

#### 2. Environment
Create a `.env` file (the app uses godotenv) or export the variables:

```bash
export DB_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
export API_SECRET="replace-with-a-secure-secret"
export POLKA_KEY="optional-polka-key"
export PLATFORM="local"
```

Or create a `.env` file with the same key/value pairs.

#### 3. Create DB schema

Apply the SQL located in the `sql/` directory to your PostgreSQL instance. Example (replace the filename if different):

```bash
psql "$DB_URL" -f sql/schema.sql
```

If you're unsure which file is the schema, inspect the `sql/` directory for files that create tables (e.g., `schema.sql`, `migrations.sql`, etc.).

#### 4. (Optional) Regenerate sqlc code

If you edit SQL files, regenerate the Go DB layer with sqlc:

```bash
sqlc generate
```

This uses `sqlc.yaml` in the repository to find SQL and output locations.

#### 5. Run the server

From the repository root you can run the app directly with `go run` or build a binary.

```bash
# run without building a binary
go run .

# or build and run
go build -o chirpy .
./chirpy
```

By default the server listens on port 8080 and serves static files under `/app/`.

#### 6. Quick smoke test

```bash
curl http://localhost:8080/api/healthz
curl http://localhost:8080/api/chirps
# Visit static UI
# http://localhost:8080/app/index.html
```

## API endpoints

to be continued...