# Refactor: Remove In-App Authentication and Introduce Claims-Based Authorization

## 1. Remove Authentication Responsibility

- [x] 1.1 Remove all JWT validation logic (no in-app token parsing or signature validation)
- [x] 1.2 Delete session-related components:
  - `internal/adapter/http/handlers/session.go`
  - `internal/application/usecases/session/usecase.go`
  - `internal/application/usecases/session/usecase_test.go`
  - `internal/application/ports/session.go`
  - `internal/adapter/database/repository/session_repository.go`
  - `internal/adapter/database/model/session_model.go`
  - `internal/adapter/database/model/session_model_test.go`
- [x] 1.3 Delete JWT package:
  - `pkg/jwt/jwt.go`
  - `pkg/jwt/jwt_test.go`
- [x] 1.4 Ensure the application no longer reads or parses the `Authorization` header

---

## 2. Introduce Trusted Identity via Headers

- [x] 2.1 Define standard headers injected by the Lambda Authorizer:
  - `X-User-ID`
  - `X-User-Email`
  - `X-User-Role`
  - (optional) `X-User-Scopes`
- [x] 2.2 Document that these headers are trusted **only when requests come through API Gateway**
- [x] 2.3 Ensure the application fails fast if required headers are missing

---

## 3. Replace Auth Middleware with Claims-Based Authorization

- [x] 3.1 Delete auth middleware:
  - `internal/adapter/http/middlewares/auth.go`
  - `internal/adapter/http/middlewares/auth_test.go`
- [x] 3.2 Create `internal/adapter/http/middlewares/role.go` with:

```go
func RoleRequired(role string) gin.HandlerFunc
```

Behavior:
- Reads `X-User-Role` from request headers
- Grants access if: role matches OR `role == administrator`
- Returns `403 Forbidden` otherwise

- [x] 3.3 Add header helper functions:
  - `func GetUserID(c *gin.Context) string`
  - `func GetUserEmail(c *gin.Context) string`
  - `func GetUserRole(c *gin.Context) string`
- [x] 3.4 Ensure helpers read directly from headers (not from gin context storage)

---

## 4. Introduce Claims Extraction

- [x] 4.1 Create a `UserClaims` struct:

```go
type UserClaims struct {
    ID    string
    Email string
    Role  string
}
```

- [x] 4.2 Implement claims extractor:

```go
func ExtractClaims(c *gin.Context) (*UserClaims, error)
```

- [x] 4.3 Validation rules:
  - Missing headers → return `401 Unauthorized`
  - Empty values → return `401 Unauthorized`

---

## 5. Update Router

- [x] 5.1 Remove auth middleware initialization from `internal/adapter/http/router.go`
- [x] 5.2 Remove `/sessions` route group
- [x] 5.3 Remove `protected.Use(authMiddleware.AuthRequired())`
- [x] 5.4 Keep route grouping but without automatic authentication
- [x] 5.5 Apply `RoleRequired(...)` middleware only where needed
- [x] 5.6 Ensure public routes do not use any authorization middleware

---

## 6. Update Dependency Wiring

- [x] 6.1 Remove session dependencies from `cmd/api/main.go`
- [x] 6.2 Remove `sessionService` from `NewRouter(...)`
- [x] 6.3 Remove JWT config from `internal/adapter/config/config.go`
- [x] 6.4 Remove JWT environment variables from `.env.example`:
  - `JWT_SECRET`
  - `JWT_ACCESS_TOKEN_EXPIRY`
  - `JWT_REFRESH_TOKEN_EXPIRY`

---

## 7. Security Guardrails

- [x] 7.1 Ensure the application is only accessible via API Gateway (not publicly exposed)
- [x] 7.2 Reject requests missing `X-User-ID` on protected routes
- [ ] 7.3 (Optional) Validate internal header: `X-Internal-Auth: true`
- [x] 7.4 Do not implement fallback authentication logic inside the application

---

## 8. Database Migration

- [x] 8.1 Create a migration to drop the `sessions` table

---

## 9. Verification

- [x] 9.1 Run build: `go build ./...`
- [x] 9.2 Run tests: `go test ./internal/...`
- [x] 9.3 Ensure no JWT usage remains: `grep -r "pkg/jwt" .`
- [x] 9.4 Ensure no session usage remains: `grep -r "usecases/session\|ports/session" .`
- [ ] 9.5 Run linter: `golangci-lint run`
- [ ] 9.6 Manual validation:
  - Request without headers → `401 Unauthorized`
  - Request with wrong role → `403 Forbidden`
  - Request with correct role → `200 OK`
  - Request with `administrator` role → bypass allowed
