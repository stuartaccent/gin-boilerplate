# Gin Boilerplate

* Manage the database schema and apply migrations in `db/migrations`.
* Manage the queries in `db/sqlc`.
* Generated query directory `internal/db`.

### Install Go Tools
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Config
For configuration see the `config.toml` passed in as the `-app-config` flag to server and cli commands.

### Generate SQL Helpers
make sure to use the correct db dsn in `sqlc.yml` and that the db is fully migrated.
```bash
sqlc generate
```

### Run Server
```bash
go run cmd/server/server.go
```

### Build Server
```bash
go build cmd/server/server.go
```

### Migrations

new:
```bash
migrate create -ext sql -dir db/migrations \
-seq <do_something>
```

up:
```bash
migrate -source file://db/migrations \
-database "postgres://postgres:password@localhost:5432/gin-boilerplate?sslmode=disable" up
```

down:
```bash
migrate -source file://db/migrations \
-database "postgres://postgres:password@localhost:5432/gin-boilerplate?sslmode=disable" down
```

### CLI
list commands:
```bash
go run cmd/cli/cli.go
```

usage:
```bash
go run cmd/cli/cli.go createuser -help
go run cmd/cli/cli.go hexauthkey -help
```
```bash
go run cmd/cli/cli.go createuser \
-email admin@example.com \
-password password \
-firstname Admin \
-lastname User
```

### Tailwind
install:
```bash
npm i
```
watch:
```bash
npm run css
```

### Docker
build:
```bash
docker build -t wip-go:latest .
```
run:
```bash
docker run \
--rm \
--name wip-go \
--publish "80:80" \
wip-go:latest \
-app-config config.toml
```
cli:
```bash
docker exec -it <container-id> ./cli
```