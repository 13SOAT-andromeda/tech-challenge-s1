## REMOVED Requirements

### Requirement: Session domain entity exists in codebase
**Reason**: The `Session` domain entity (`internal/domain/session.go`) is no longer used by any service, repository, or handler after the `remove-jwt-auth-middleware` change. Keeping it creates misleading dead code in the domain layer.
**Migration**: No migration needed. The type has no callers. Delete `internal/domain/session.go` and `internal/domain/session_test.go`.

#### Scenario: Session entity is unreachable
- **WHEN** any application code path is executed
- **THEN** no code SHALL reference or instantiate the `Session` domain struct

### Requirement: Session endpoints documented in Swagger spec
**Reason**: The `/sessions` endpoints (login, logout, refresh, validate) were removed in the `remove-jwt-auth-middleware` change. Their presence in `swagger/swagger.yaml` is inaccurate and misleading to API consumers.
**Migration**: Remove `/sessions` paths and all `Session*` schema definitions from `swagger/swagger.yaml`.

#### Scenario: Swagger spec contains no session endpoints
- **WHEN** a client reads the OpenAPI spec at `/swagger/swagger.yaml`
- **THEN** the spec SHALL NOT contain any `/sessions` path entries or session-related schema definitions

### Requirement: Session requests present in Postman collection
**Reason**: The Postman collection at `misc/Tech Challenge S1.postman.json` contains session-related requests that correspond to now-deleted endpoints. These requests will always fail and mislead developers.
**Migration**: Remove session-related request entries (login, logout, refresh, validate) from the Postman collection file.

#### Scenario: Postman collection contains no session requests
- **WHEN** a developer opens the Postman collection
- **THEN** the collection SHALL NOT contain any request targeting `/sessions/*` endpoints
