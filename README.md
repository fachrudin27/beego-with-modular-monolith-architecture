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
    domain/           request/response/domain models
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
shared        -> cross-cutting infrastructure helpers
```

Practical rules:

- Controllers may know Beego.
- Services should use `context.Context`, not Beego context.
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
```

For non-dev environments, set a real JWT secret. The application will reject the default secret outside `dev` mode.

You can also set:

```bash
export JWT_SECRET="your-secure-secret"
export JWT_ISSUER="firstbeegoapi"
```

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

Planned next updates:

- Database integration and repository layer
- Database migrations
- Dockerfile and Docker Compose setup
- Prometheus metrics endpoint
- Grafana dashboard setup
- Loki log shipping example
- Request access log middleware
- More unit tests for services and middleware
