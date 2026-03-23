## Context

The application currently handles its own JWT validation via an in-app `authMiddleware` that parses the `Authorization` header, validates token signatures, and checks session state in the database. This is being removed because an AWS Lambda authorizer sits upstream and validates JWTs before any request reaches this service on private endpoints.

However, some routes still need role-based authorization (e.g. admin-only operations). Since the Lambda already validated the token, it can forward user identity as request context headers. The app must read those headers instead of re-parsing the JWT.

Current auth flow:
```
Client → Lambda Authorizer (JWT validate) → API Gateway → [App: JWT re-validate + session check] → Handler
```

Target auth flow:
```
Client → Lambda Authorizer (JWT validate + inject headers) → API Gateway → [App: role check from headers] → Handler
```

## Goals / Non-Goals

**Goals:**
- Remove all JWT parsing, signature validation, and session persistence from the app.
- Remove the session lifecycle endpoints (`/sessions/*`).
- Introduce a lightweight `RoleRequired(role)` middleware that reads user role from a Lambda-forwarded request header.
- Ensure admin users can always access role-restricted endpoints (admin bypass).
- Simplify `NewRouter` signature by removing `sessionService` and auth middleware dependencies.

**Non-Goals:**
- Changing any business logic in existing handlers.
- Implementing token issuance or refresh — that belongs to the upstream auth service.
- Modifying the Lambda authorizer itself.
- Adding per-user resource scoping beyond role checks.

## Decisions

### Decision 1: Identity propagation via request headers

**Choice**: Lambda authorizer injects identity into request headers before forwarding to the app. The app reads `X-User-ID`, `X-User-Email`, and `X-User-Role` headers.

**Rationale**: AWS API Gateway Lambda authorizers can inject context into headers via the `context` object returned from the authorizer. This is the standard pattern and requires zero in-app crypto dependencies.

**Alternative considered**: Read from `requestContext` via a custom middleware parsing the API Gateway request context JSON — more complex, tightly coupled to AWS, harder to test locally.

---

### Decision 2: New `RoleRequired(role string)` middleware reads from headers

**Choice**: Replace the old JWT-based `RoleRequired` with a new middleware that reads `X-User-Role` from the request header. Access is granted if the role matches OR the role is `admin`.

**Rationale**: Simple, stateless, testable without JWT infrastructure. The Lambda guarantees the header is authentic — no need to re-verify.

**Alternative considered**: Keep using the old `RoleRequired` but skip JWT validation — would leave dead code (jwt package, session lookup) in place.

---

### Decision 3: Remove session persistence entirely

**Choice**: Drop the `sessions` table, session repository, session service, and session use case.

**Rationale**: Sessions were only used to support the in-app token validation flow (checking if a token was revoked). With Lambda owning auth, this state is no longer meaningful or maintained by this service.

**Alternative considered**: Keep sessions for audit logging — out of scope for this change; can be added later if needed.

---

### Decision 4: Remove `pkg/jwt` package

**Choice**: Delete `pkg/jwt/` entirely once no other package depends on it.

**Rationale**: The package's sole purpose was in-app token signing and validation. Removing the auth middleware removes all callers.

**Verification required**: Grep for `pkg/jwt` imports before deleting to confirm no other package depends on it.

## Risks / Trade-offs

- **Header spoofing if Lambda is bypassed** → Any direct call to the app (bypassing API Gateway) would have no role headers, defaulting to no role. Mitigation: ensure the app is not publicly accessible except through API Gateway; use security groups / VPC to enforce this.
- **`X-User-Role` header name must match Lambda output** → If the Lambda uses a different header name, role checks silently fail (empty role = access denied). Mitigation: agree on header names with the team deploying the Lambda and document them in `.env.example` or config.
- **Session endpoint removal is a breaking API change** → Any client calling `/sessions/*` will get 404. Mitigation: document in release notes; clients must migrate to the upstream auth service.
- **Tests referencing auth middleware or session use case must be deleted or rewritten** → Leaves some coverage gaps temporarily. Mitigation: covered by e2e tests that run against the full stack (Lambda included in staging).

## Migration Plan

1. Deploy the Lambda authorizer to staging and verify it injects `X-User-Role`, `X-User-ID`, `X-User-Email` headers.
2. Merge this change to a feature branch and deploy to staging behind the Lambda.
3. Run e2e tests confirming protected routes work with Lambda-forwarded headers.
4. Run a database migration to drop the `sessions` table.
5. Remove JWT env vars (`JWT_SECRET`, `JWT_ACCESS_TOKEN_EXPIRY`, `JWT_REFRESH_TOKEN_EXPIRY`) from all environment configs.
6. Merge to main and deploy to production.

**Rollback**: Revert the commit. Re-add JWT env vars. The `sessions` table drop is destructive — restore from a pre-migration backup if needed (session data is ephemeral so data loss is acceptable).

## Open Questions

- What exact header names does the Lambda authorizer inject? Assumed `X-User-ID`, `X-User-Email`, `X-User-Role` — needs confirmation.
- Should the new `RoleRequired` middleware return `403` when the role header is missing (request bypassed Lambda) vs. `401`? Currently assuming `403 Forbidden`.
- Which specific routes (if any) should use `RoleRequired` in the new setup? The current router only uses `AuthRequired` at the group level — role checks per-route may be a follow-up.
