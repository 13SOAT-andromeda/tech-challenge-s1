## 1. Configuration

- [x] 1.1 Add `JWTConfig` struct to `internal/adapter/config/config.go` with a `Secret string` field
- [x] 1.2 Load `JWT_SECRET` env var into `JWTConfig` inside `config.Init()`
- [x] 1.3 Add `JWT_SECRET=your-secret-here` to `.env.example`
- [x] 1.4 Attach `JWT *JWTConfig` to the top-level `Config` struct

## 2. JWT Package

- [x] 2.1 Create `pkg/jwt/jwt.go` with a `Claims` struct embedding `jwt.RegisteredClaims` plus `Email string` and `Role string` fields
- [x] 2.2 Implement `ValidateToken(tokenString, secret string) (*Claims, error)` that parses with `jwt.ParseWithClaims`, enforces HS256 only, and returns the validated claims

## 3. Auth Middleware

- [x] 3.1 Create `internal/adapter/http/middlewares/auth.go` with `AuthRequired(secret string) gin.HandlerFunc`
- [x] 3.2 Extract the Bearer token from the `Authorization` header; return `401` if absent or malformed
- [x] 3.3 Call `pkg/jwt.ValidateToken`; return `401` on any validation error
- [x] 3.4 Store claims as `*middlewares.UserClaims` under context key `"userClaims"` via `c.Set`

## 4. Update Role Middleware

- [x] 4.1 Update `ExtractClaims` in `internal/adapter/http/middlewares/role.go` to read `"userClaims"` from Gin context (`c.Get`) instead of Lambda-forwarded headers
- [x] 4.2 Update `GetUserID`, `GetUserEmail`, `GetUserRole` helpers to read from the context `UserClaims` struct
- [x] 4.3 Update `RoleRequired` to return `401` (not `403`) when claims are missing from context

## 5. Router Wiring

- [x] 5.1 Add `jwtSecret string` (or `JWTConfig`) parameter to `NewRouter` in `internal/adapter/http/router.go`
- [x] 5.2 Apply `middlewares.AuthRequired(jwtSecret)` to the `protected` route group
- [x] 5.3 Update `cmd/api/main.go` to pass `cfg.JWT.Secret` when calling `http.NewRouter`

## 6. Tests

- [x] 6.1 Write unit tests for `pkg/jwt.ValidateToken` covering: valid token, expired token, wrong algorithm, bad secret, malformed string
- [x] 6.2 Write unit tests for `AuthRequired` middleware: valid token proceeds, missing header returns 401, invalid token returns 401
- [x] 6.3 Update or add unit tests for `ExtractClaims` / `RoleRequired` covering context-based claims instead of header-based
