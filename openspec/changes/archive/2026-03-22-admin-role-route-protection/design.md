## Context

The API already has two middleware layers: `AuthRequired` (validates JWT and stores `UserClaims` in context) and `RoleRequired(role string)` (reads claims and enforces a role). The six resource handler groups — customers, maintenances, orders, products, users, vehicles — sit inside a `protected` group that applies `AuthRequired` but not `RoleRequired`. Any user with a valid JWT, regardless of role, can currently call these endpoints.

The `RoleRequired` middleware already handles the administrator bypass (`claims.Role == "administrator"` always passes), so no new logic needs to be written — only the router wiring needs to change.

## Goals / Non-Goals

**Goals:**
- Gate all routes under `CustomerHandler`, `MaintenanceHandler`, `OrderHandler`, `ProductHandler`, `UserHandler`, and `VehicleHandler` with `RoleRequired("administrator")`.
- Non-administrator authenticated users receive `403 Forbidden`.
- Public order approval/rejection routes remain unchanged.

**Non-Goals:**
- Changing `CompanyHandler` routes (out of scope).
- Introducing new roles or role hierarchies.
- Changing the JWT structure or claim format.
- Adding per-endpoint role granularity (all six groups get the same `administrator` restriction).

## Decisions

### Apply `RoleRequired` at the group level, not per-route

Adding the middleware once per route group (e.g., `customerGroup.Use(middlewares.RoleRequired("administrator"))`) is cleaner and less error-prone than adding it to each individual route. It also makes the intent visible at the group definition rather than scattered across individual registrations.

**Alternatives considered:**
- Per-route middleware: More flexible for future role granularity but verbose and easy to miss a route.
- A combined `AdminRequired` middleware that does auth + role check: Would duplicate the `AuthRequired` logic; keeping them separate follows single-responsibility and allows other roles to reuse `RoleRequired` independently.

### Middleware order: `AuthRequired` → `RoleRequired`

`RoleRequired` reads claims set by `AuthRequired`, so it must run after. The existing `protected` group already applies `AuthRequired`; adding `Use(RoleRequired(...))` inside each sub-group ensures the correct order without restructuring.

## Risks / Trade-offs

- **BREAKING change for non-admin users** → Non-administrator tokens that currently succeed will now receive `403`. Any client or integration test using non-admin credentials against these endpoints will fail. Mitigation: ensure all E2E tests authenticate as an administrator, and document the change clearly.
- **No role granularity** → All six handler groups are locked to `administrator`. If a future role (e.g., `mechanic`) needs read-only access to orders, the router will need per-route or per-group refinement. This is acceptable given current requirements.
- **E2E test breakage** → Tests that call these endpoints with non-admin JWTs will return `403` instead of `200`. Mitigation: update `NewAuthenticatedReq` or test setup to use admin credentials.
