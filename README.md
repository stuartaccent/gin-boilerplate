# Gin Boilerplate

* Manage the database schema and apply migrations in `db/migrations`.
* Manage the queries in `db/sqlc`.
* Generated query directory `internal/db`.
* Manage the templates in `components`.

## Tooling
for migrations:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
for sqlc:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
for template rendering:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

## Config
For configuration see the `config.toml` passed in as the `-app-config` flag to server and cli commands.

> [!TIP]
> Generate new session auth keys ask [ChatGPT](https://chat.openai.com)
> 
> "Could you generate random hex keys of 32 bytes and 16 bytes for me?"

## SQL Helpers
make sure to use the correct db dsn in `sqlc.yml` and that the db is fully migrated.
```bash
sqlc generate
```

## Run Server
```bash
go run cmd/server/main.go
```

## Build Server
```bash
go build -o server cmd/server/main.go
```

## Migrations

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

## Templates

Generate template code with [templ.guide](https://templ.guide)
```bash
templ generate -watch
```

## CLI
list commands:
```bash
go run cmd/cli/main.go
```

usage:
```bash
go run cmd/cli/main.go createuser -help
```
```bash
go run cmd/cli/main.go createuser \
-email admin@example.com \
-password password \
-firstname Admin \
-lastname User
```

## Tailwind
install:
```bash
npm i
```
watch:
```bash
npm run css
```

## Docker
build:
```bash
docker build -t gin-boilerplate:latest .
```
run:
```bash
docker run \
--rm \
--name gin-boilerplate \
--publish "80:80" \
--env "DATABASE_HOST=host.docker.internal" \
--env "SERVER_MODE=release" \
gin-boilerplate:latest \
-app-config config.toml
```
cli:
```bash
docker exec -it <container-id> ./cli
```

## Postgres Container
create the volume:
```bash
docker volume create gin-postgres-data
```
create the container:
```bash
docker run --detach \
--name "gin-postgres" \
--volume "gin-postgres-data:/var/lib/postgresql/data" \
--publish "5432:5432" \
--env POSTGRES_USER=postgres \
--env POSTGRES_PASSWORD=password \
--env POSTGRES_DB=gin-boilerplate \
postgres:latest
```
cleanup:
```bash
docker stop gin-postgres
docker rm gin-postgres
docker volume rm gin-postgres-data
```