## ADDED Requirements

### Requirement: Administrator-only access on resource management routes
The router SHALL apply `RoleRequired("administrator")` middleware to the route groups for `CustomerHandler`, `MaintenanceHandler`, `OrderHandler`, `ProductHandler`, `UserHandler`, and `VehicleHandler`. Authenticated users whose JWT `role` claim is not `administrator` MUST receive `403 Forbidden` when accessing any route in these groups.

#### Scenario: Administrator accesses a protected resource route
- **WHEN** a request with a valid JWT where `role=administrator` is sent to `GET /customers`
- **THEN** the request SHALL proceed to the handler and return a success response

#### Scenario: Non-administrator is rejected on customer routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `GET /customers`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Non-administrator is rejected on maintenance routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `POST /maintenances`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Non-administrator is rejected on order routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `GET /orders`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Non-administrator is rejected on product routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `GET /products`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Non-administrator is rejected on user routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `GET /users`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Non-administrator is rejected on vehicle routes
- **WHEN** a request with a valid JWT where `role=mechanic` is sent to `GET /vehicles`
- **THEN** the middleware SHALL return `403 Forbidden` and abort the request

#### Scenario: Public order approval route remains accessible without role restriction
- **WHEN** a request without any JWT is sent to `GET /orders/:id/approve`
- **THEN** the request SHALL proceed to the handler (no auth or role check applies)

### Requirement: Middleware order preserved within protected groups
The `RoleRequired("administrator")` middleware SHALL always run after `AuthRequired` within each route group, so that `UserClaims` are already present in the Gin context when role enforcement executes.

#### Scenario: Claims are available when RoleRequired runs
- **WHEN** a request with a valid JWT reaches a protected resource route
- **THEN** `AuthRequired` sets `UserClaims` in context before `RoleRequired` reads them
