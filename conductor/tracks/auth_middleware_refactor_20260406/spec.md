# Specification: Refactor Auth & Role Middlewares for API Gateway Authorizer Headers

## Overview
Refactor the `auth.go` middleware to extract user information (`X-User-Id`, `X-User-Email`, `X-User-Role`) directly from incoming HTTP headers instead of existing authentication mechanisms. Set these values in the Go request context. Subsequently, refactor the `role.go` middleware to utilize this new context data for authorization. Finally, refactor any handlers that currently rely on the old authentication logic to use the new context values. This change relies on the API Gateway acting as the authorizer and passing the verified headers.

## Functional Requirements
1. **Header Extraction (`auth.go`):**
   - Completely replace the existing authentication logic (e.g., JWT extraction).
   - Read the following headers passed by the API Gateway: `X-User-Id`, `X-User-Email`, `X-User-Role`.
   - Validate that all three headers are present and not empty. If any are missing or empty, immediately return a `401 Unauthorized`.
   - Validate data types:
     - `X-User-Id` must be a valid number (e.g., parsable to integer).
     - `X-User-Email` must be a non-empty string.
     - `X-User-Role` must be a non-empty string.
   - Set the extracted values into the Go request context using standard string keys (e.g., `user_id`, `user_email`, `user_role`).

2. **Authorization Updates (`role.go`):**
   - Update the role middleware to read the user's role from the updated Go request context (using the string key `user_role`) instead of its previous source.
   - Maintain the existing role verification logic.

3. **Handler Refactoring:**
   - Identify and refactor any handlers that currently extract user information (ID, Email, Role) using the old authentication logic (e.g., parsing tokens directly from the context).
   - Update these handlers to retrieve the information using the new context string keys (`user_id`, `user_email`, `user_role`).

4. **Testing:**
   - Update or rewrite unit tests for both `auth.go` and `role.go` middlewares to reflect the new header-based extraction and context setting.
   - Cover success scenarios (valid headers) and failure scenarios (missing headers, empty headers, invalid ID format).
   - Update any handler tests affected by the context changes.

## Out of Scope
- Implementing new authentication logic (e.g., JWT validation) within the application itself, as this is now delegated to the API Gateway.