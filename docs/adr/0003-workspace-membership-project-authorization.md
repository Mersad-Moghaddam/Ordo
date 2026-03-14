# ADR 0003: Workspace, Membership, and Project Authorization Boundaries

## Status
Accepted

## Context
Phase 2 introduces collaborative organization structures with workspaces, memberships, and projects. Access control decisions must remain deterministic and reusable across delivery adapters.

## Decision
Implement workspace core orchestration in the usecase layer with explicit role checks:
- owner/admin can invite members.
- owner/admin can create projects.
- owner can change member role.

Persist workspaces, memberships, and projects in dedicated MySQL tables with unique keys and foreign keys.

Expose list operations through paginated usecase responses using the generic page result model.

## Consequences
Positive outcomes:
- Deterministic authorization behavior independent of HTTP handlers.
- Clear data model for future task and activity entities tied to projects.
- SQL schema constraints enforce uniqueness and relational integrity.

Tradeoffs:
- Additional join/query complexity for workspace listing.
- Increased migration and repository surface area.
