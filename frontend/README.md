# Ordo Frontend (Jira-style flow)

This frontend now follows a clear page flow:

1. **Login page** (`/login`)
2. **Create Account page** (`/signup`) from "Create Account" button
3. **Main page** (`/main`) after successful login/signup

The main page includes forms for **all backend routes/features**: auth refresh, workspaces, memberships, projects, tasks, comments, and activities.

## Tech

- React 18 + TypeScript
- Vite
- Simple SPA path router (no extra dependency)
- Reusable component-based UI
- Responsive production-style CSS

## Run

```bash
cd frontend
npm install
npm run dev
```

## Build

```bash
npm run build
```

## Route coverage in UI

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /workspaces`
- `POST /workspaces/:workspaceId/memberships`
- `POST /workspaces/:workspaceId/projects`
- `GET /users/:userId/workspaces`
- `GET /workspaces/:workspaceId/projects`
- `POST /tasks`
- `PATCH /tasks/:taskId/status`
- `GET /projects/:projectId/tasks`
- `POST /comments`
- `PATCH /comments/:commentId`
- `DELETE /comments/:commentId`
- `GET /tasks/:taskId/comments`
- `GET /tasks/:taskId/activities`
