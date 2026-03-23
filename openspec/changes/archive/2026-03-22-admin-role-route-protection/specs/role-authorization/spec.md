## MODIFIED Requirements

### Requirement: Role-based authorization via JWT claims
Protected routes that require a specific role SHALL validate the user's role by reading the `UserClaims` struct stored in the Gin context by the `AuthRequired` middleware. A user with the `administrator` role MUST be granted access regardless of the required role. The route groups for `CustomerHandler`, `MaintenanceHandler`, `OrderHandler`, `ProductHandler`, `UserHandler`, and `VehicleHandler` MUST be configured with `RoleRequired("administrator")` in the router.

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

#### Scenario: Non-administrator is denied access to resource management endpoints
- **WHEN** a request carries a valid JWT with `role=mechanic` and targets any route under `/customers`, `/maintenances`, `/orders`, `/products`, `/users`, or `/vehicles`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request
