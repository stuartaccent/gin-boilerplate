# Gin Boilerplate

* Manage the database schema and apply migrations in `pkg/storage/db/migrations`.
* Manage the queries in `pkg/storage/db/sqlc`.
* Generated query directory `pkg/storage/db/dbx`.
* Manage the templates in `pkg/ui`.

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
for air reloading:
```bash
go install github.com/air-verse/air@latest
```

## Config
For configuration see the `config.toml` passed in as the `--config` flag to app.

> [!TIP]
> Generate new session auth keys ask [ChatGPT](https://chat.openai.com)
> 
> "Could you generate random hex keys of 32 bytes and 16 bytes for me?"

## SQL Helpers
make sure to use the correct db dsn in `sqlc.yml` and that the db is fully migrated.
```bash
sqlc generate
```

## Migrations

new, requires `golang-migrate` cli to be installed:
```bash
migrate create -ext sql -dir pkg/storage/db/migrations \
-seq <do_something>
```

up:
```bash
go run . migrate --config config.dev.toml up
```

down:
```bash
go run . migrate --config config.dev.toml down
```

## Templates

Generate template code with [templ.guide](https://templ.guide)
```bash
templ generate
```

## Usage

```bash
go run . help
```

run the server:
```bash
go run . server --config config.dev.toml
```

create a user:
```bash
go run . createuser --config config.dev.toml \
--email admin@example.com \
--password password \
--firstname Admin \
--lastname User
```

## Tailwind

download the cli to the tmp folder (created once air is run):
```bash 
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 tmp/tailwindcss
```
build and watch:
```bash
./tmp/tailwindcss \
-i ./pkg/ui/css/global.css \
-o ./pkg/static/public/css/global.css \
--watch \
--minify
```

## Docker
This can be used to build the app as a binary and run it in a container.

build:
```bash
docker build --build-arg EXPOSE_PORT=80 -t gin-boilerplate:latest .
```

run:
```bash
docker run \
--rm \
--name gin-boilerplate \
--publish "80:80" \
--env "DATABASE_HOST=host.docker.internal" \
--env "SERVER_MODE=release" \
--env "SERVER_PORT=80" \
gin-boilerplate:latest \
server --config config.toml
```

## Postgres Container

create the container:
```bash
docker run \
--detach \
--name "gin-postgres" \
--mount type=tmpfs,destination=/var/lib/postgresql/data \
--publish "5432:5432" \
--env POSTGRES_USER=postgres \
--env POSTGRES_PASSWORD=password \
--env POSTGRES_DB=gin-boilerplate \
postgres:latest
```

## Links

* [Boilerplate](https://github.com/stuartaccent/gin-boilerplate)
* [Golang](https://go.dev)
* [Templ](https://templ.guide)
* [Air](https://github.com/air-verse/air)
* [Icônes](https://icones.js.org/collection/lucide)
* [Tailwind](https://tailwindcss.com)
* [Owl](https://github.com/AccentDesign/owl)
* [Shadcn](https://ui.shadcn.com/docs)
* [Htmx](https://htmx.org/)