## Context

The `remove-jwt-auth-middleware` change deleted all runtime session infrastructure (handler, use case, service, repository, DB model, JWT package). However, the following orphaned session references were not removed in that change:

- `internal/domain/session.go` — `Session` struct and helper methods (`NewSession`, `IsExpired`, `IsValid`)
- `internal/domain/session_test.go` — unit tests for the above
- `swagger/swagger.yaml` — `/sessions/*` path documentation and `Session*` schema definitions
- `misc/Tech Challenge S1.postman.json` — session-related API requests

None of these files are referenced by any live code path. Removing them completes the cleanup.

## Goals / Non-Goals

**Goals:**
- Delete the `Session` domain entity and its tests.
- Strip `/sessions` paths and session-related schema definitions from the Swagger spec.
- Remove session requests from the Postman collection.

**Non-Goals:**
- Changing any other domain entity, service, or handler.
- Modifying the database migration file `drop_session_model.sql` (it remains as the mechanism to drop the table from existing environments).
- Changing any Kubernetes, Terraform, or CI/CD configuration.

## Decisions

### Decision 1: Delete domain entity rather than leave it unused

**Choice**: Delete `internal/domain/session.go` entirely.

**Rationale**: Unused domain types add noise to the codebase and are misleading — they imply a concept the system still manages. The hexagonal architecture principle of keeping the domain layer free of dead code applies here.

**Alternative**: Leave it in place as documentation — rejected because it is actively misleading after removing all callers.

---

### Decision 2: Edit Swagger and Postman manually rather than regenerating

**Choice**: Remove session-related entries directly from `swagger/swagger.yaml` and `misc/*.postman.json`.

**Rationale**: These files appear to be committed manually (not auto-generated from annotations). Deleting or regenerating from scratch risks losing non-session content. Surgical edits are safer.

## Risks / Trade-offs

- **Swagger breaking change** → Any client code generated from the Swagger spec that references session schemas or endpoints will break. Mitigation: this is expected and intentional — the `/sessions` endpoints no longer exist.
- **Postman collection** → Team members with cached local copies of the collection may still see session requests. Mitigation: commit the updated collection so the team pulls the clean version.
