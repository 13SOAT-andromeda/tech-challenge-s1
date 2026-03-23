## Context

The application uses Gin + hexagonal architecture. Until recently, JWT authentication was handled externally by an AWS Lambda Authorizer that injected identity headers (`X-User-ID`, `X-User-Role`, `X-User-Email`) before requests reached the app. That authorizer layer was removed along with the internal session management code (`pkg/jwt/`, `middlewares/auth.go`). The current `middlewares/role.go` reads those trusted headers, but since no one sets them anymore, every endpoint is effectively unprotected.

The existing `role.go` `RoleRequired` middleware and `ExtractClaims` helper are still valid — only the *source* of the identity data needs to change (JWT → Gin context → role check).

## Goals / Non-Goals

**Goals:**
- Validate the JWT signature on every request that hits a protected route.
- Extract claims (`sub`/user ID, `email`, `role`) from the token and store them in the Gin context.
- Update `RoleRequired` and identity helpers to read from the Gin context instead of Lambda headers.
- Keep the change minimal and non-breaking for handlers that already use `GetUserID` / `GetUserRole`.

**Non-Goals:**
- Token issuance / login endpoint (no session creation).
- Refresh tokens or token revocation.
- Claims beyond sub, email, and role.
- Changing any business logic or handler code.

## Decisions

### 1. Where to validate the JWT

**Decision:** A Gin middleware `AuthRequired` in `internal/adapter/http/middlewares/auth.go`, applied to the existing `protected` group in `router.go`.

**Alternatives considered:**
- Validate inside each handler — rejected: too much duplication.
- Use a separate application-layer port — rejected: overkill for a stateless signature check with no business rules.

### 2. Claims storage in context

**Decision:** After validation, store claims as a `*UserClaims` struct under the key `"userClaims"` in the Gin context (`c.Set`). Update `ExtractClaims` in `role.go` to call `c.Get("userClaims")` instead of reading headers.

**Alternatives considered:**
- Keep reading from headers (set them ourselves after validating) — rejected: indirection with no benefit.
- Store individual string values per key — rejected: a struct is safer and easier to extend.

### 3. JWT library

**Decision:** Use `github.com/golang-jwt/jwt/v5`, already present in `go.mod`.

### 4. Config injection

**Decision:** Add a `JWTConfig` struct to `config.go` with a `Secret` field (from `JWT_SECRET` env var). Pass it into `NewRouter` so the middleware is initialized with the signing key.

**Alternatives considered:**
- Read env var directly in middleware — rejected: untestable and bypasses the existing config pattern.

### 5. Algorithm

**Decision:** Accept only `HS256` (HMAC-SHA256). Reject tokens signed with any other algorithm to prevent the "alg:none" attack.

## Risks / Trade-offs

- **Secret rotation:** Changing `JWT_SECRET` instantly invalidates all issued tokens. → Mitigation: document this; rotation requires a coordinated redeploy.
- **Stateless validation only:** Revoked/logged-out tokens remain valid until expiry. → Accepted trade-off per non-goals; a token denylist is out of scope.
- **Lambda header trust removed:** If the app is ever placed behind API Gateway + Lambda Authorizer again, the headers will be ignored and the Bearer token will be validated instead. Both flows are compatible since the Gin context keys remain identical.

## Migration Plan

1. Add `JWT_SECRET` to `.env` / deployment secrets.
2. Deploy the new code (no DB migrations needed).
3. All clients must include a valid `Authorization: Bearer <token>` header on protected requests.
4. Rollback: revert the middleware wiring in `router.go`; no data changes to undo.

## Open Questions

- None — scope is narrowly defined and all dependencies are already in-tree.
