## Why

The previous change (`remove-jwt-auth-middleware`) removed the session infrastructure from the application, but leftover references to the `Session` domain entity still exist in the codebase, Swagger docs, and the Postman collection. These orphaned references are misleading and should be removed to keep the codebase consistent with the actual system behavior.

## What Changes

- Delete `internal/domain/session.go` — the `Session` domain entity is no longer used by any service, repository, or handler.
- Delete `internal/domain/session_test.go` — tests for the now-deleted domain entity.
- Remove session-related endpoint documentation from `swagger/swagger.yaml` (`/sessions` paths and `Session*` schema definitions).
- Remove session-related requests from `misc/Tech Challenge S1.postman.json` (login, logout, refresh, validate entries).

## Capabilities

### New Capabilities

- None. This is a pure cleanup — no new capabilities are introduced.

### Modified Capabilities

- None. No spec-level behavior changes.

## Impact

- **Domain layer**: `internal/domain/session.go` and its test deleted; no other domain file references `Session`.
- **API docs**: Swagger spec loses `/sessions/*` paths and session-related request/response schemas.
- **Postman collection**: Session-related requests removed from `misc/`.
- **No runtime impact**: The `Session` type is already unused by all handlers, services, and repositories after the previous change.
