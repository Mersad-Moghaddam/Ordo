# Ordo Frontend (Production-style UI)

A **React + TypeScript + Vite** frontend that fully maps to the current backend API surface:

- Auth (register/login/refresh)
- Workspace and memberships
- Projects
- Tasks and status transitions
- Comments and activity feed
- Admin operations overview

## Tech stack

- React 18
- TypeScript (strict)
- Vite
- ESLint
- Lucide icons
- Responsive CSS design system (custom)

## Quick start

```bash
cd frontend
npm install
npm run dev
```

App runs on `http://localhost:4173` by default.

## Production build

```bash
npm run build
npm run preview
```

## Backend mapping

The UI includes forms wired to all backend endpoints exposed in:

- `/auth/register`
- `/auth/login`
- `/auth/refresh`
- `/workspaces`
- `/workspaces/:workspaceId/memberships`
- `/workspaces/:workspaceId/projects`
- `/users/:userId/workspaces`
- `/workspaces/:workspaceId/projects` (GET)
- `/tasks`
- `/tasks/:taskId/status`
- `/projects/:projectId/tasks`
- `/comments`
- `/comments/:commentId` (PATCH, DELETE)
- `/tasks/:taskId/comments`
- `/tasks/:taskId/activities`

## UX highlights

- API base URL switcher in top navigation
- Session token persistence
- Dedicated sections for Auth / Workspace / Tasks / Collaboration / Admin
- Structured JSON response console for all actions
- Mobile-friendly adaptive layout
