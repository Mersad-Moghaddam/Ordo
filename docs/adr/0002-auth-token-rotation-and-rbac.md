# ADR 0002: Auth Token Rotation and RBAC Baseline

## Status
Accepted

## Context
Phase 1 introduces user authentication, refresh token lifecycle, and role-aware authorization while preserving clean architecture boundaries and testability.

## Decision
Implement an auth usecase service responsible for registration, login, refresh rotation, and role authorization.

Persist user records in `users` and refresh token sessions in `refresh_sessions`.

Use stateless signed tokens for access and refresh claims, plus server-side hashed opaque refresh token material bound to a rotating session version.

Adopt refresh token rotation by revoking the current session and issuing a replacement session for every successful refresh.

Enforce RBAC through usecase-level `AuthorizeRole` logic and HTTP middleware that verifies bearer tokens and required roles.

## Consequences
Positive outcomes:
- Rotation reduces replay impact for stolen refresh tokens.
- Session table enables explicit revocation and auditability.
- Authorization logic remains independent from delivery adapters.

Tradeoffs:
- Additional persistence operations on refresh flow.
- More moving parts than a simple non-rotating token model.
