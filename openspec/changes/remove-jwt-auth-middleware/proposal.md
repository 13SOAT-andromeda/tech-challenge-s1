## Why

The application is deployed behind an AWS Lambda authorizer that validates JWT tokens before routing requests to private endpoints, making in-app JWT validation redundant and adding unnecessary complexity. Removing it simplifies the codebase, reduces dependencies, and eliminates a now-unnecessary layer of auth logic that duplicates what the infrastructure already enforces.

## What Changes

- **BREAKING** Remove the AuthRequired() middleware from all protected routes. Replace the existing RoleRequired() usage to ensure that every route explicitly validates authorization based on user roles, allowing access when the user has either the role configured for the endpoint or the admin role.
- **BREAKING** Remove the session endpoints (`POST /sessions`, `GET /sessions/validate`, `POST /sessions/refresh`, `DELETE /sessions/logout`) — login and token lifecycle are no longer managed by this service.
- Remove `internal/adapter/http/middlewares/auth.go` and its test.
- Remove `internal/application/usecases/session/` use case and its test.
- Remove session-related ports (`SessionService`, `SessionRepository`) from `internal/application/ports/`.
- Remove session DB model, repository implementation, and any session migration/table.
- Remove the `pkg/jwt/` package (no longer used internally).
- Remove the `session_handler` and `session_service` wiring from `cmd/api/main.go` and `router.go`.
- Remove JWT-related environment variables (`JWT_SECRET`, `JWT_ACCESS_TOKEN_EXPIRY`, `JWT_REFRESH_TOKEN_EXPIRY`) from config.
- Remove context helpers (`GetUserID`, `GetUserEmail`, `GetClaims`) from the middleware package, or replace with simpler header-extraction helpers if the Lambda forwards identity headers.

## Capabilities

### New Capabilities

- None. This is a pure removal — no new capabilities are introduced.

### Modified Capabilities

- None. No existing spec files are present; no spec-level behavior changes.

## Impact

- **Routes**: All protected routes lose their in-app auth guard; security is now fully delegated to the Lambda authorizer upstream.
- **Session endpoints**: `POST/GET/DELETE /sessions/*` are removed entirely — clients must use the Lambda/auth service for token issuance.
- **Code removed**: `middlewares/auth.go`, `usecases/session/`, `ports/session_*.go`, `adapter/database/repository/session*.go`, `adapter/database/model/session*.go`, `pkg/jwt/`.
- **Config**: JWT config block and related env vars removed from `config.go` and `.env.example`.
- **Tests**: Auth middleware tests and session use case tests removed.
- **`main.go` / `router.go`**: Session service, session handler, and auth middleware initialization removed; `NewRouter` signature simplified.
