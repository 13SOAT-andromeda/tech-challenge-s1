## Why

The application previously delegated all JWT validation to an AWS Lambda Authorizer, which injected trusted identity headers (`X-User-ID`, `X-User-Role`, etc.) into every request. That session layer was removed, leaving all endpoints unprotected. The app needs its own lightweight JWT guard so that protected routes can enforce authentication and role-based access control without depending on an external authorizer.

## What Changes

- Add a `pkg/jwt` package that validates a JWT's signature and extracts its claims (sub, email, role).
- Add an `AuthRequired` middleware in `internal/adapter/http/middlewares/auth.go` that reads the `Authorization: Bearer <token>` header, validates the JWT, and stores the parsed claims in the Gin context.
- Update `RoleRequired` (and the `ExtractClaims` helper in `role.go`) to read claims from the Gin context (set by the JWT middleware) instead of Lambda-forwarded headers.
- Add `JWT_SECRET` to `config.go` and `.env.example` so the signing key is injected at runtime.
- Wire the `AuthRequired` middleware into the protected route group in `router.go` and pass the config to `NewRouter`.

## Capabilities

### New Capabilities
- `jwt-auth`: JWT signature validation and claims extraction middleware for protected endpoints.

### Modified Capabilities
- `role-authorization`: Identity source changes from Lambda-injected headers to JWT claims stored in the Gin context by the new auth middleware.

## Impact

- **`pkg/jwt/`** — new package (validator + claims struct); replaces deleted `pkg/jwt/jwt.go`.
- **`internal/adapter/http/middlewares/auth.go`** — new middleware file; replaces deleted one.
- **`internal/adapter/http/middlewares/role.go`** — `ExtractClaims` and getter helpers updated to read from Gin context keys instead of HTTP headers.
- **`internal/adapter/config/config.go`** — new `JWTConfig` struct with `Secret` field; loaded from `JWT_SECRET` env var.
- **`internal/adapter/http/router.go`** — `AuthRequired` applied to the `protected` group; `JWTConfig` passed in.
- **`cmd/api/main.go`** — pass JWT config when constructing the router.
- **`.env.example`** — add `JWT_SECRET` example variable.
- No new external dependencies; uses `github.com/golang-jwt/jwt/v5` (already a transitive dep or standard choice for Go JWT).
