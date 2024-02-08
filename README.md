# Gin Boilerplate

* Manage the database schema and apply migrations in `db/migrations`.
* Manage the queries in `db/sqlc`.
* Generated query directory `internal/db`.

### Install Tools
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Generate SQL Helpers
make sure to use the correct db dsn in `sqlc.yml` and that the db is fully migrated.
```bash
sqlc generate
```

### Run Server
```bash
go run cmd/server/server.go \
-db-dns "postgres://postgres:password@localhost:5432/gin-boilerplate?sslmode=disable" \
-port 80
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
-db "postgres://postgres:password@localhost:5432/gin-boilerplate?sslmode=disable" \
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
docker build -t gin-boilerplate-go:latest .
```
run:
```bash
docker run \
--rm \
--name gin-boilerplate-go \
--publish "80:80" \
gin-boilerplate-go:latest \
-db-dns "postgres://postgres:password@host.docker.internal:5432/gin-boilerplate?sslmode=disable" \
-port 80
```
cli:
```bash
docker exec -it <container-id> ./cli
```