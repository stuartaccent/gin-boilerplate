# Gin Boilerplate

* Manage the database schema and apply migrations in `db/migrations`.
* Manage the queries in `db/sqlc`.
* Generated query directory `db/dbx`.
* Manage the templates in `ui`.

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
migrate create -ext sql -dir db/migrations \
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
templ generate -watch
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

## Load Testing

Just a basic test with wrk, no real load testing was done. Logging was disabled for the test.

logged in with real user and session cookie, postgres db on localhost:
```bash
wrk -c 400 -t 10 -d 30s -H "Cookie: session=..." http://localhost
```

    Running 30s test @ http://localhost
        10 threads and 400 connections
        Thread Stats   Avg      Stdev     Max   +/- Stdev
          Latency    20.20ms    2.00ms  71.96ms   88.29%
          Req/Sec     1.99k   144.93     2.24k    85.63%
        593341 requests in 30.02s, 796.16MB read
        Socket errors: connect 0, read 402, write 0, timeout 0
    Requests/sec:  19763.59
    Transfer/sec:     26.52MB

not logged in:
```bash
wrk -c 400 -t 10 -d 30s http://localhost/auth/login
```

    Running 30s test @ http://localhost/auth/login
        10 threads and 400 connections
        Thread Stats   Avg      Stdev     Max   +/- Stdev
          Latency     5.40ms    3.32ms  86.58ms   84.56%
          Req/Sec     7.68k   645.82    16.49k    89.63%
        2294661 requests in 30.04s, 4.52GB read
        Socket errors: connect 0, read 406, write 0, timeout 0
    Requests/sec:  76392.72
    Transfer/sec:    154.23MB
