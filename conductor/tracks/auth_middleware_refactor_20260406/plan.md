# Implementation Plan: Refactor Auth & Role Middlewares for Header-based Context

## Phase 1: Refactor `auth.go` Middleware [checkpoint: 2eaabe9]
- [x] Task: Write failing tests for `auth.go` to cover header extraction (`X-User-Id`, `X-User-Email`, `X-User-Role`), validation (presence, emptiness, ID format), and context setting (using string keys). [ede03d2]
- [x] Task: Implement header extraction and validation in `auth.go` to pass the tests, completely replacing the old logic. [1809ea7]
- [x] Task: Conductor - User Manual Verification 'Phase 1: Refactor auth.go Middleware' (Protocol in workflow.md) [1809ea7]

## Phase 2: Refactor `role.go` Middleware [checkpoint: 314058c]
- [x] Task: Write failing tests for `role.go` to cover reading the user role from the new Go request context string key (`user_role`). [03f2995]
- [x] Task: Implement the updated context reading logic in `role.go` to pass the tests. [f2655fa]
- [x] Task: Conductor - User Manual Verification 'Phase 2: Refactor role.go Middleware' (Protocol in workflow.md) [f2655fa]

## Phase 3: Refactor Handlers [checkpoint: 0c01269]
- [x] Task: Identify all handlers that currently extract user information from the context using the old authentication logic. [1809ea7]
- [x] Task: Update the unit tests for the identified handlers to mock the new context string keys. [0c01269]
- [x] Task: Refactor the identified handlers to retrieve `user_id`, `user_email`, and `user_role` using the new string keys from the request context. [0c01269]
- [x] Task: Conductor - User Manual Verification 'Phase 3: Refactor Handlers' (Protocol in workflow.md) [0c01269]