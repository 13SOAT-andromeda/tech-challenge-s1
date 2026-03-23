## Why

All resource management routes (customers, maintenances, orders, products, users, vehicles) are currently protected only by JWT authentication, meaning any authenticated user regardless of role can call them. The system has a `RoleRequired` middleware and an `administrator` role already defined but not enforced on these routes.

## What Changes

- Apply `RoleRequired("administrator")` middleware to the route groups for `CustomerHandler`, `MaintenanceHandler`, `OrderHandler`, `ProductHandler`, `UserHandler`, and `VehicleHandler` in the router.
- Non-administrator authenticated users will receive `403 Forbidden` when attempting to access any of those routes.
- Public order approval/rejection routes (`GET /orders/:id/approve` and `GET /orders/:id/reject`) remain unprotected as-is.
- The `CompanyHandler` routes are out of scope (not listed in the request).

## Capabilities

### New Capabilities

- `admin-route-protection`: Enforce administrator-only access on all listed resource management route groups via the existing `RoleRequired` middleware.

### Modified Capabilities

- `role-authorization`: Add requirement that the six handler groups (customers, maintenances, orders, products, users, vehicles) MUST be gated by `RoleRequired("administrator")` in addition to `AuthRequired`.

## Impact

- `internal/adapter/http/router.go`: Add `middlewares.RoleRequired("administrator")` to the six route groups.
- No changes to domain, application, or database layers.
- Any existing non-administrator JWT users will lose access to these endpoints — **BREAKING** for non-admin API consumers.
- E2E tests that create requests for these endpoints must use an administrator token.
