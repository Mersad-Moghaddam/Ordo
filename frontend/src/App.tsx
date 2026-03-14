import { useMemo, useState } from 'react'

type ApiResponse<T> = {
  status: number
  body: T | string
}

type TokenBody = {
  accessToken: string
  refreshToken: string
}

type WorkspaceBody = {
  WorkspaceID: string
}

type ProjectBody = {
  ProjectID: string
}

type TaskBody = {
  TaskID: string
}

type CommentBody = {
  CommentID: string
}

const defaultBaseUrl = 'http://127.0.0.1:8080'

async function callApi<T>(baseUrl: string, path: string, method: string, payload?: unknown): Promise<ApiResponse<T>> {
  const response = await fetch(`${baseUrl}${path}`, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: payload ? JSON.stringify(payload) : undefined,
  })
  const responseText = await response.text()
  try {
    return { status: response.status, body: JSON.parse(responseText) as T }
  } catch {
    return { status: response.status, body: responseText }
  }
}

export default function App() {
  const [baseUrl, setBaseUrl] = useState(defaultBaseUrl)
  const [ownerEmail, setOwnerEmail] = useState('owner@ordo.dev')
  const [ownerPassword, setOwnerPassword] = useState('secret')
  const [refreshToken, setRefreshToken] = useState('')
  const [workspaceId, setWorkspaceId] = useState('')
  const [projectId, setProjectId] = useState('')
  const [taskId, setTaskId] = useState('')
  const [commentId, setCommentId] = useState('')
  const [feed, setFeed] = useState<string[]>([])

  const logger = useMemo(() => ({
    push(entry: string) {
      setFeed((existing) => [entry, ...existing].slice(0, 80))
    }
  }), [])

  async function runFullFlow() {
    logger.push('Starting full API flow…')
    const register = await callApi<TokenBody>(baseUrl, '/api/v1/auth/register', 'POST', { email: ownerEmail, password: ownerPassword, role: 'owner' })
    logger.push(`register -> ${register.status}`)

    const login = await callApi<TokenBody>(baseUrl, '/api/v1/auth/login', 'POST', { email: ownerEmail, password: ownerPassword })
    logger.push(`login -> ${login.status}`)
    if (typeof login.body !== 'string') {
      setRefreshToken(login.body.refreshToken)
    }

    if (typeof login.body !== 'string' && login.body.refreshToken) {
      const refresh = await callApi<TokenBody>(baseUrl, '/api/v1/auth/refresh', 'POST', { refreshToken: login.body.refreshToken })
      logger.push(`refresh -> ${refresh.status}`)
    }

    const workspace = await callApi<WorkspaceBody>(baseUrl, '/api/v1/workspaces', 'POST', { workspaceKey: 'platform', displayName: 'Platform Workspace', ownerUserId: 'owner-user' })
    logger.push(`create workspace -> ${workspace.status}`)
    if (typeof workspace.body !== 'string') {
      setWorkspaceId(workspace.body.WorkspaceID)
      const membership = await callApi<unknown>(baseUrl, `/api/v1/workspaces/${workspace.body.WorkspaceID}/memberships`, 'POST', { actorUserId: 'owner-user', targetUserId: 'admin-user', targetRole: 'admin', invitedByUserId: 'owner-user' })
      logger.push(`add membership -> ${membership.status}`)

      const project = await callApi<ProjectBody>(baseUrl, `/api/v1/workspaces/${workspace.body.WorkspaceID}/projects`, 'POST', { actorUserId: 'admin-user', projectKey: 'api', displayName: 'API Project', description: 'Backend and integration project' })
      logger.push(`create project -> ${project.status}`)
      if (typeof project.body !== 'string') {
        setProjectId(project.body.ProjectID)

        const task = await callApi<TaskBody>(baseUrl, '/api/v1/tasks', 'POST', {
          workspaceId: workspace.body.WorkspaceID,
          projectId: project.body.ProjectID,
          title: 'Validate full flow',
          description: 'Smoke validate all endpoint groups',
          priority: 'high',
          createdByUserId: 'owner-user',
        })
        logger.push(`create task -> ${task.status}`)
        if (typeof task.body !== 'string') {
          setTaskId(task.body.TaskID)

          const updateTask = await callApi<unknown>(baseUrl, `/api/v1/tasks/${task.body.TaskID}/status`, 'PATCH', { status: 'in_progress' })
          logger.push(`update task status -> ${updateTask.status}`)

          const comment = await callApi<CommentBody>(baseUrl, '/api/v1/comments', 'POST', {
            workspaceId: workspace.body.WorkspaceID,
            projectId: project.body.ProjectID,
            taskId: task.body.TaskID,
            authorUserId: 'owner-user',
            body: 'Initial comment from UI flow',
          })
          logger.push(`create comment -> ${comment.status}`)
          if (typeof comment.body !== 'string') {
            setCommentId(comment.body.CommentID)
            const updateComment = await callApi<unknown>(baseUrl, `/api/v1/comments/${comment.body.CommentID}`, 'PATCH', { actorUserId: 'owner-user', body: 'Updated comment from UI flow' })
            logger.push(`update comment -> ${updateComment.status}`)
          }
        }
      }
    }

    const health = await callApi<unknown>(baseUrl, '/health', 'GET')
    const metrics = await callApi<unknown>(baseUrl, '/metrics', 'GET')
    logger.push(`health -> ${health.status}, metrics -> ${metrics.status}`)
  }

  async function listCurrent() {
    if (!workspaceId || !projectId || !taskId) {
      logger.push('Need workspace/project/task IDs first (run flow).')
      return
    }
    const userWorkspaces = await callApi<unknown>(baseUrl, '/api/v1/users/owner-user/workspaces', 'GET')
    const workspaceProjects = await callApi<unknown>(baseUrl, `/api/v1/workspaces/${workspaceId}/projects`, 'GET')
    const projectTasks = await callApi<unknown>(baseUrl, `/api/v1/projects/${projectId}/tasks`, 'GET')
    const taskComments = await callApi<unknown>(baseUrl, `/api/v1/tasks/${taskId}/comments`, 'GET')
    const taskActivities = await callApi<unknown>(baseUrl, `/api/v1/tasks/${taskId}/activities`, 'GET')
    logger.push(`lists -> workspaces:${userWorkspaces.status} projects:${workspaceProjects.status} tasks:${projectTasks.status} comments:${taskComments.status} activities:${taskActivities.status}`)
  }

  return (
    <div className="page">
      <header className="hero">
        <h1>Ordo Command Center</h1>
        <p>Modern full-stack workspace/task orchestration UI wired to your backend endpoints.</p>
      </header>

      <section className="panel">
        <h2>Connection</h2>
        <label>Backend Base URL</label>
        <input value={baseUrl} onChange={(event) => setBaseUrl(event.target.value)} />
        <div className="grid-two">
          <div>
            <label>Owner Email</label>
            <input value={ownerEmail} onChange={(event) => setOwnerEmail(event.target.value)} />
          </div>
          <div>
            <label>Owner Password</label>
            <input value={ownerPassword} onChange={(event) => setOwnerPassword(event.target.value)} />
          </div>
        </div>
        <div className="actions">
          <button onClick={() => void runFullFlow()}>Run Full Register → Task → Comment Flow</button>
          <button className="secondary" onClick={() => void listCurrent()}>Validate List Endpoints</button>
        </div>
      </section>

      <section className="panel">
        <h2>Live IDs & Tokens</h2>
        <p><strong>Refresh Token:</strong> <span className="mono">{refreshToken || '—'}</span></p>
        <p><strong>Workspace ID:</strong> <span className="mono">{workspaceId || '—'}</span></p>
        <p><strong>Project ID:</strong> <span className="mono">{projectId || '—'}</span></p>
        <p><strong>Task ID:</strong> <span className="mono">{taskId || '—'}</span></p>
        <p><strong>Comment ID:</strong> <span className="mono">{commentId || '—'}</span></p>
      </section>

      <section className="panel">
        <h2>Execution Log</h2>
        <div className="feed">
          {feed.length === 0 ? <p>No events yet. Run the full flow.</p> : feed.map((entry, index) => <p key={`${entry}-${index}`}>{entry}</p>)}
        </div>
      </section>
    </div>
  )
}
