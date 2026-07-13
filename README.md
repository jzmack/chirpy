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

## How to deploy

These instructions are for setting up on Linux systems. Instructions validated on WSL2.

Prerequisites:
- Go 1.26.x
- PostgreSQL
- Goose - for database migrations
- sqlc - to regenerate DB code

### 1. Clone Repo
```bash
git clone https://github.com/jzmack/chirpy.git
cd chirpy
```

### 2. Install prerequisites

#### Go installation 
Go documentation: https://go.dev/doc/install

That said, I recommend installing Go with the following:

```bash
curl -sS https://webi.sh/golang | sh
```

Validate Go installation with:
```bash
go version
```

#### Goose & sqlc installation

Once Go is installed, Goose and sqlc can be installed with:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
Validate with:
```bash
goose --version
sqlc version
```

#### PostgreSQL installation

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
Validate PostgreSQL installation with:
```bash
psql --version # validate install
sudo service postgres status # should be active
```

### 3. Setting up environment

Ensure you have a user that can launch the database. The user `postgres` can be used and is the user created when PostgreSQL was installed. Password can be set with:

```bash
sudo passwd postgres
```

To create the PostgreSQL database, first launch the psql client with:
```bash
sudo -u postgres psql
# should bring you to a  prompt that looks like --> postgres=#
```

While in the postgres prompt, create the database with:

```sql
CREATE DATABASE chirpy;
```

Set the database user (`postgres` is also the name of the user in this example) and password with the following and then `exit` the prompt:

```sql
ALTER USER postgres WITH PASSWORD 'yourpassword';
```

Test connection to the database using the following from the Linux shell (be sure to update `user` and `pass` in the connection string):

```bash
psql postgres://user:pass@localhost:5432/chirpy
# if connection succeeds, prompt should look like --> chirpy=#
# type 'exit' to leave the prompt
```

Once connection string is confirmed, change into the `sql/schema` directory and run a goose up migration with:

```bash
cd sql/schema
goose postgres postgres://user:pass@localhost:5432/chirpy up
```
Example output looks like:

```plain
jzm@sandor-lnx:~/workspace/chirpy/sql/schema$ goose postgres postgres://postgres:postgres@localhost:5432/chirpy up
2026/07/13 15:13:53 OK   001_users.sql (7.2ms)
2026/07/13 15:13:53 OK   002_chirps.sql (4.63ms)
2026/07/13 15:13:53 OK   003_passwords.sql (2.45ms)
2026/07/13 15:13:53 OK   004_refresh_tokens.sql (4.38ms)
2026/07/13 15:13:53 OK   005_chirpy_red.sql (2.04ms)
2026/07/13 15:13:53 goose: successfully migrated database to version: 5
```

A `.env` file can be used(the app uses godotenv) with the following the variables:

```plain
DB_URL="postgres://user:pass@localhost:5432/chirpy"
API_SECRET="replace_with_API_secret"
POLKA_KEY="replace_with_polka_key"
PLATFORM="local"
```

#### 4. (Optional) Regenerate sqlc code

The code included under the `sql` will work as is but if you edit SQL files, regenerate the Go DB layer with sqlc:

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
# http://localhost:8080/app
```

## Test API endpoints

to be continued...