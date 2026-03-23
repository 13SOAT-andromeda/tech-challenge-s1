## 1. Remove Session Domain Entity

- [x] 1.1 Delete `internal/domain/session.go`
- [x] 1.2 Delete `internal/domain/session_test.go`
- [x] 1.3 Run `go build ./...` to confirm no remaining callers

## 2. Clean Up Swagger Spec

- [x] 2.1 Remove all `/sessions` path entries from `swagger/swagger.yaml`
- [x] 2.2 Remove `Session`-related schema definitions from `swagger/swagger.yaml`
- [x] 2.3 Verify the spec is valid YAML after edits

## 3. Clean Up Postman Collection

- [x] 3.1 Remove session-related requests (login, logout, refresh, validate) from `misc/Tech Challenge S1.postman.json`

## 4. Verification

- [x] 4.1 Run `go build ./...` — zero errors
- [x] 4.2 Run `go test ./internal/...` — all tests pass
- [x] 4.3 Confirm no session references remain: `grep -rn "domain/session\|domain.Session\|NewSession" . --include="*.go"`
