## ADDED Requirements

### Requirement: JWT signature validation on protected routes
The `AuthRequired` middleware SHALL verify the JWT signature on every request to a protected route using the configured `JWT_SECRET` key and the HS256 algorithm. Tokens signed with any other algorithm MUST be rejected.

#### Scenario: Valid token is accepted
- **WHEN** a request arrives with `Authorization: Bearer <valid-hs256-token>`
- **THEN** the middleware SHALL parse the claims and call `c.Next()` to proceed to the handler

#### Scenario: Token with wrong algorithm is rejected
- **WHEN** a request arrives with a JWT signed using RS256 or `alg:none`
- **THEN** the middleware SHALL return `401 Unauthorized` and abort the request

#### Scenario: Expired token is rejected
- **WHEN** a request arrives with a JWT whose `exp` claim is in the past
- **THEN** the middleware SHALL return `401 Unauthorized` and abort the request

#### Scenario: Malformed or missing token is rejected
- **WHEN** a request arrives with no `Authorization` header, or a non-Bearer value, or a malformed JWT string
- **THEN** the middleware SHALL return `401 Unauthorized` and abort the request

### Requirement: JWT claims stored in Gin context
After successful validation the middleware SHALL store the parsed claims as a `*UserClaims` struct under the context key `"userClaims"` so that downstream middleware and handlers can retrieve identity without re-parsing the token.

#### Scenario: Claims available after middleware runs
- **WHEN** a valid JWT with claims `sub=42`, `email=user@example.com`, `role=mechanic` is validated
- **THEN** `c.Get("userClaims")` SHALL return a `*UserClaims` with `ID="42"`, `Email="user@example.com"`, `Role="mechanic"`

### Requirement: JWT secret injected via configuration
The signing secret SHALL be read from the `JWT_SECRET` environment variable through the existing `Config` struct and passed into the middleware at startup. The middleware MUST NOT read environment variables directly.

#### Scenario: Missing JWT_SECRET at startup
- **WHEN** the application starts without a `JWT_SECRET` environment variable
- **THEN** the config loader SHALL use an empty string default, and the middleware SHALL reject all tokens (effectively disabling protected routes)
