## REMOVED Requirements

### Requirement: In-app JWT validation
**Reason**: JWT validation is now fully handled by the upstream AWS Lambda authorizer. Re-validating in the app is redundant and couples the app to JWT infrastructure.
**Migration**: Remove `AuthRequired()` middleware from all route groups. The Lambda authorizer ensures only authenticated requests reach the app. No in-app replacement is needed for authentication.

### Requirement: Session lifecycle management
**Reason**: Sessions were only maintained to support token revocation checks inside the app. With the Lambda owning auth, the app no longer manages token lifecycle.
**Migration**: Remove `POST /sessions`, `GET /sessions/validate`, `POST /sessions/refresh`, `DELETE /sessions/logout` endpoints. Clients MUST use the upstream auth service for login, token refresh, and logout.

## ADDED Requirements

### Requirement: Role-based authorization via Lambda-forwarded headers
Protected routes that require a specific role SHALL validate the user's role by reading the `X-User-Role` request header injected by the Lambda authorizer. A user with the `admin` role MUST be granted access regardless of the required role.

#### Scenario: Request with matching role is allowed
- **WHEN** a request arrives with `X-User-Role: mechanic` and the route requires role `mechanic`
- **THEN** the request SHALL proceed to the handler

#### Scenario: Admin bypasses role restriction
- **WHEN** a request arrives with `X-User-Role: admin` and the route requires any non-admin role
- **THEN** the request SHALL proceed to the handler

#### Scenario: Request with insufficient role is rejected
- **WHEN** a request arrives with `X-User-Role: mechanic` and the route requires role `admin`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Request with missing role header is rejected
- **WHEN** a request arrives without the `X-User-Role` header
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

### Requirement: User identity available in handler context
Handlers that need user identity SHALL be able to read `X-User-ID`, `X-User-Email`, and `X-User-Role` from the request headers forwarded by the Lambda authorizer. The app MUST NOT re-parse or re-validate these values from a JWT.

#### Scenario: Handler reads user ID from header
- **WHEN** the Lambda authorizer injects `X-User-ID: 42` into the request
- **THEN** the handler SHALL be able to retrieve the value `42` as the current user's ID

#### Scenario: Handler reads user role from header
- **WHEN** the Lambda authorizer injects `X-User-Role: admin` into the request
- **THEN** the handler SHALL be able to retrieve `admin` as the current user's role
