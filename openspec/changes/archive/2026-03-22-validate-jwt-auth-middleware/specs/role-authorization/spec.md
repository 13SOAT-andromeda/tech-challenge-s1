## MODIFIED Requirements

### Requirement: Role-based authorization via JWT claims
Protected routes that require a specific role SHALL validate the user's role by reading the `UserClaims` struct stored in the Gin context by the `AuthRequired` middleware. A user with the `administrator` role MUST be granted access regardless of the required role.

#### Scenario: Request with matching role is allowed
- **WHEN** a request carries a valid JWT with `role=mechanic` and the route requires role `mechanic`
- **THEN** the request SHALL proceed to the handler

#### Scenario: Administrator bypasses role restriction
- **WHEN** a request carries a valid JWT with `role=administrator` and the route requires any non-administrator role
- **THEN** the request SHALL proceed to the handler

#### Scenario: Request with insufficient role is rejected
- **WHEN** a request carries a valid JWT with `role=mechanic` and the route requires role `administrator`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Request with missing claims is rejected
- **WHEN** no `userClaims` key is present in the Gin context (i.e., `AuthRequired` was not applied or failed silently)
- **THEN** `RoleRequired` SHALL return `401 Unauthorized` and abort the request

### Requirement: User identity available in handler context
Handlers that need user identity SHALL retrieve it by calling `ExtractClaims(c)`, which reads the `*UserClaims` struct from the Gin context key `"userClaims"`. The app MUST NOT read identity from raw HTTP headers or re-parse the JWT in handlers.

#### Scenario: Handler reads user ID from context
- **WHEN** the `AuthRequired` middleware stores `UserClaims{ID: "42"}` in the context
- **THEN** `GetUserID(c)` SHALL return `"42"`

#### Scenario: Handler reads user role from context
- **WHEN** the `AuthRequired` middleware stores `UserClaims{Role: "administrator"}` in the context
- **THEN** `GetUserRole(c)` SHALL return `"administrator"`
