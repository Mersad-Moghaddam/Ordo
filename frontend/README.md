# Ordo Frontend

Single-page frontend for the Ordo backend API.

## Features covered

- **Authentication pages**: Sign In, Sign Up, token refresh.
- **Main dashboard**: workspace/project/task/comment/activity forms.
- **Admin page**: centralized admin overview and operations guidance.
- **API response viewer**: live JSON output for every action.

## Run locally

This frontend is static and does not require a build step.

```bash
cd frontend
python3 -m http.server 4173
```

Then open: <http://localhost:4173>

By default, API calls go to:

- `http://localhost:8080/api/v1`

You can change the API base URL in the top navbar.

## Backend endpoints wired

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
