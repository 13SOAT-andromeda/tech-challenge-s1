## 1. Router Update

- [x] 1.1 Add `middlewares.RoleRequired("administrator")` to the `customerGroup` in `internal/adapter/http/router.go`
- [x] 1.2 Add `middlewares.RoleRequired("administrator")` to the `maintenances` group in `internal/adapter/http/router.go`
- [x] 1.3 Add `middlewares.RoleRequired("administrator")` to the `orderGroup` in `internal/adapter/http/router.go`
- [x] 1.4 Add `middlewares.RoleRequired("administrator")` to the `productGroup` in `internal/adapter/http/router.go`
- [x] 1.5 Add `middlewares.RoleRequired("administrator")` to the `userGroup` in `internal/adapter/http/router.go`
- [x] 1.6 Add `middlewares.RoleRequired("administrator")` to the `vehicleGroup` in `internal/adapter/http/router.go`

## 2. Test Updates

- [x] 2.1 Verify E2E test helper `NewAuthenticatedReq` (or equivalent setup) uses an administrator token for requests to the six protected handler groups
- [x] 2.2 Update any E2E tests in `test/e2e/` that call customer, maintenance, order, product, user, or vehicle endpoints with non-admin credentials — switch to admin credentials or expect `403`
- [x] 2.3 Add or update middleware unit test in `internal/adapter/http/middlewares/` to assert that a non-administrator JWT receives `403 Forbidden` on a route protected by `RoleRequired("administrator")`

## 3. Validation

- [x] 3.1 Run `go test ./internal/...` and confirm all unit tests pass
- [ ] 3.2 Run `go test ./test/e2e/...` against a running instance and confirm all E2E tests pass
- [ ] 3.3 Manually verify with a non-admin JWT that `GET /customers` returns `403 Forbidden`
- [ ] 3.4 Manually verify with an administrator JWT that `GET /customers` returns `200 OK`
