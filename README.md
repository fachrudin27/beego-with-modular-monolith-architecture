# beego-with-modular-monolith-architecture

Starter backend API using **Beego**, **Domain-Driven Design style layering**, and a **modular monolith** structure.

The project is designed as a practical foundation for building Go APIs where each business module owns its own delivery, application, and domain layer while still running as one deployable application.

## Highlights

- Beego v2 HTTP server and routing
- Modular monolith package structure
- DDD-inspired layering per module
- JWT authentication middleware
- Request ID middleware
- Standard API response and error format
- Native Zap structured logging
- Loki-friendly JSON logs
- Request and response logging with sensitive field redaction
- PostgreSQL database integration for the ordering module
- Repository layer with application-level repository contracts
- Module-owned database migrations
- IP-based rate limiting
- Endpoint tests using Go test

## Architecture

The codebase is organized by business capability, not by technical type only.

```text
internal/
  auth/
    delivery/api/     HTTP controller layer
    app/              application use cases
    domain/           request/response/domain models
    client.go         module contract

  ordering/
    delivery/api/     HTTP controller layer
    app/              application use cases
      repository.go   repository contract used by the app layer
    domain/           request/response/domain models
    infra/postgres/   PostgreSQL connection, repository implementation, migrations
    client.go         module contract

  shared/
    error.go          standard API errors and responses
    jwt.go            JWT generation, parsing, middleware
    request_id.go     request correlation ID
    zap.go            native Zap logger
    log_context.go    log metadata passed through context
    rate_limiting.go  IP rate limiting
    cors.go           CORS middleware

routers/
  router.go           API route registration and middleware wiring
```

## Layering Rules

The intended boundaries are:

```text
delivery/api  -> handles Beego HTTP details
app           -> contains use case/business flow
domain        -> contains module data contracts
infra         -> contains concrete storage/external integrations
shared        -> cross-cutting infrastructure helpers
```

Practical rules:

- Controllers may know Beego.
- Services should use `context.Context`, not Beego context.
- Services depend on repository interfaces, not concrete PostgreSQL implementations.
- PostgreSQL implementations live under module `infra/postgres`.
- Migrations are owned by the module storage implementation.
- Domain types should not depend on HTTP, logger, or Beego.
- Technical error details should be logged, not returned directly to users.
- User responses should stay stable and safe.

## Request Flow

Example protected endpoint flow:

```text
HTTP request
  -> CORS middleware
  -> request ID middleware
  -> rate limit middleware
  -> JWT middleware
  -> Beego controller
  -> application service
  -> repository interface
  -> PostgreSQL repository implementation
  -> shared response writer
```

## API Response Format

Success:

```json
{
  "message": "success get data",
  "data": {
    "product_id": 1
  },
  "request_id": "request-id",
  "create_in": 1780000000
}
```

Error:

```json
{
  "message": "request failed",
  "error": {
    "kind": "validation",
    "code": "invalid_object_id",
    "message": "object id must be a number"
  },
  "request_id": "request-id",
  "create_in": 1780000000
}
```

## Logging

This project uses native Zap directly, not Beego logs.

Logs are structured and intended to be friendly for Grafana Loki:

```json
{
  "level": "info",
  "ts": "2026-06-30T10:00:00+07:00",
  "msg": "Ordering Get API Log",
  "service": "ordering",
  "position": "/api",
  "request_id": "request-id",
  "url": "/v1/ordering/1",
  "request": "",
  "response": "{\"product_id\":1}"
}
```

Sensitive fields are redacted before being written to logs:

- `password`
- `token`
- `access_token`
- `refresh_token`
- `authorization`

Log bodies are also truncated to avoid oversized log entries.

## Authentication

Login endpoint:

```http
POST /v1/auth/login
```

Example:

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'
```

Protected endpoint:

```bash
curl http://localhost:8080/v1/ordering/1 \
  -H "Authorization: Bearer <access_token>"
```

Current demo credentials:

```text
email: admin@example.com
password: password
```

Replace the dummy authentication logic with a real user repository before using this project beyond local development.

## Configuration

Copy the example config:

```bash
cp conf/app.conf.example conf/app.conf
```

Important values:

```ini
httpport = 8080
runmode = dev
copyrequestbody = true
jwtsecret = change-this-secret
jwtissuer = firstbeegoapi
jwtexpiresin = 3600

ordering_postgres_dsn =
ordering_postgres_host = localhost
ordering_postgres_port = 5432
ordering_postgres_user = postgres
ordering_postgres_password = postgres
ordering_postgres_database = firstbeegoapi
ordering_postgres_sslmode = disable
ordering_postgres_max_open_conns = 10
ordering_postgres_max_idle_conns = 5
ordering_postgres_conn_max_lifetime_seconds = 300
ordering_postgres_ping_timeout_seconds = 5
```

For non-dev environments, set a real JWT secret. The application will reject the default secret outside `dev` mode.

You can also set:

```bash
export JWT_SECRET="your-secure-secret"
export JWT_ISSUER="firstbeegoapi"
```

## Docker Setup

Docker support is now started, with the current layout split by responsibility:

```text
Dockerfile                  container build for the Beego app
conf/docker-compose.yaml    local orchestration for app, postgres, and migrate
conf/app.conf               Beego application config
conf/.env                   Docker/PostgreSQL environment values
```

Current progress:

- The application image build is available through `Dockerfile`
- Compose is wired for local container orchestration
- PostgreSQL and migration containers are separated from the app container
- Beego still reads its runtime settings from `app.conf`
- Docker-specific secrets and database values are kept in `.env`
- Vector is used to ship Zap logs from Docker stdout into Loki
- Grafana is provisioned with a dashboard for log table and status chart views
- Log entries expose `level`, `request_id`, `service`, `position`, `request`, and `response`

## Observability

Logging and metrics are separated into two paths:

- Prometheus scrapes application metrics from `/metrics`
- Vector reads Zap JSON logs from Docker and forwards them to Loki
- Grafana reads Loki logs and renders a table plus a status chart

The log dashboard is focused on:

- `create_date`
- `level`
- `request_id`
- `service`
- `position`
- `request`
- `response`

The status chart groups entries by:

- `error`
- `warn`
- `warning`
- `info`
- `debug`

## Database

The ordering module has its own PostgreSQL integration under:

```text
internal/ordering/infra/postgres/
```

The application layer defines the repository contract:

```text
internal/ordering/app/repository.go
```

The PostgreSQL implementation lives in:

```text
internal/ordering/infra/postgres/repository/
```

The dependency direction is:

```text
delivery/api -> app -> app.OrderingRepository interface
infra/postgres/repository -> implements app.OrderingRepository
main.go -> wires concrete PostgreSQL repository into the ordering service
```

At runtime, `main.go` initializes the ordering PostgreSQL connection from `ordering_postgres_*` config keys and injects the repository into the ordering service.

## Database Migrations

Migrations are stored inside the module infrastructure folder, because each module may own a different storage technology.

Ordering PostgreSQL migrations:

```text
internal/ordering/infra/postgres/migrations/
  000002_create_ordering_table_order.up.sql
  000002_create_ordering_table_order.down.sql
```

This keeps storage ownership close to the module:

```text
ordering module -> postgres infra -> postgres migrations
```

If another module later uses MySQL, Elasticsearch, or another database, it can keep its own migration/setup files under that module's `infra` folder.

## Running

```bash
go run .
```

The API runs on:

```text
http://localhost:8080
```

## Testing

```bash
go test ./...
```

If your environment blocks the default Go build cache path:

```bash
GOCACHE=/tmp/firstbeegoapi-gocache go test ./...
```

## Example Endpoints

Login:

```text
POST /v1/auth/login
```

Ordering:

```text
GET /v1/ordering/:objectId
```

## Roadmap

Implemented:

- Database integration for the ordering module
- Repository layer for the ordering module
- Module-owned PostgreSQL migrations
- Dockerfile for containerizing the Beego app
- Docker Compose scaffold for app, PostgreSQL, and migrations
- Loki log shipping through Vector
- Grafana dashboard provisioning for application logs
- Prometheus metrics endpoint and scrape configuration

Planned next updates:

- Request access log middleware
- More unit tests for services and middleware
