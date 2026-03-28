import { FormEvent, useMemo, useState } from 'react'
import { ShieldCheck, LayoutDashboard, ListChecks, MessageCircle } from 'lucide-react'
import { callApi } from './lib/api'
import { RouteKey, SessionState } from './lib/types'
import { Card } from './components/Card'
import { InputField, SelectField, TextareaField } from './components/FormControls'
import { TopNav } from './layout/TopNav'

const DEFAULT_BASE = 'http://localhost:8080/api/v1'

function data(form: HTMLFormElement): Record<string, string> {
  return Object.fromEntries(new FormData(form).entries()) as Record<string, string>
}

export function App() {
  const [route, setRoute] = useState<RouteKey>('auth')
  const [session, setSession] = useState<SessionState>({
    apiBase: localStorage.getItem('apiBase') ?? DEFAULT_BASE,
    accessToken: localStorage.getItem('accessToken') ?? '',
    refreshToken: localStorage.getItem('refreshToken') ?? '',
  })
  const [output, setOutput] = useState('Ready. Start with Sign Up / Sign In.')

  const authReady = useMemo(() => Boolean(session.accessToken), [session.accessToken])

  const saveSession = (next: SessionState) => {
    setSession(next)
    localStorage.setItem('apiBase', next.apiBase)
    localStorage.setItem('accessToken', next.accessToken)
    localStorage.setItem('refreshToken', next.refreshToken)
  }

  const submit = (builder: (values: Record<string, string>) => Promise<unknown>) => async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    try {
      const result = await builder(data(event.currentTarget))
      setOutput(JSON.stringify(result, null, 2))
    } catch (error) {
      setOutput(error instanceof Error ? error.message : String(error))
    }
  }

  return (
    <div className="app-shell">
      <TopNav
        active={route}
        onSwitch={setRoute}
        apiBase={session.apiBase}
        onApiBaseChange={(value) => saveSession({ ...session, apiBase: value.trim() })}
      />

      <main className="content">
        {route === 'auth' && (
          <>
            <Card title="Sign In" subtitle="Authenticate and receive access + refresh token.">
              <form className="form" onSubmit={submit((values) => callApi(session, '/auth/login', 'POST', { email: values.email, password: values.password }))}>
                <InputField label="Email" name="email" type="email" />
                <InputField label="Password" name="password" type="password" />
                <button className="btn primary">Sign In</button>
              </form>
            </Card>
            <Card title="Sign Up" subtitle="Create new account in backend auth service.">
              <form
                className="form"
                onSubmit={submit(async (values) => {
                  const result = (await callApi(session, '/auth/register', 'POST', {
                    email: values.email,
                    password: values.password,
                    role: values.role,
                  })) as { accessToken: string; refreshToken: string }
                  saveSession({ ...session, accessToken: result.accessToken, refreshToken: result.refreshToken })
                  return result
                })}
              >
                <InputField label="Email" name="email" type="email" />
                <InputField label="Password" name="password" type="password" />
                <SelectField label="Role" name="role" options={[{ value: 'member', label: 'Member' }, { value: 'admin', label: 'Admin' }]} />
                <button className="btn primary">Create account</button>
              </form>
            </Card>
            <Card title="Token Management" subtitle="Refresh token or clear current session.">
              <div className="icon-row"><ShieldCheck size={18} /> <span>{authReady ? 'Authenticated' : 'Not authenticated'}</span></div>
              <form
                className="form"
                onSubmit={submit(async () => {
                  const result = (await callApi(session, '/auth/refresh', 'POST', {
                    refreshToken: session.refreshToken,
                  })) as { accessToken: string; refreshToken: string }
                  saveSession({ ...session, accessToken: result.accessToken, refreshToken: result.refreshToken })
                  return result
                })}
              >
                <button className="btn secondary">Refresh Token</button>
              </form>
              <button className="btn ghost" onClick={() => saveSession({ ...session, accessToken: '', refreshToken: '' })}>Sign Out</button>
            </Card>
          </>
        )}

        {route === 'workspace' && (
          <>
            <Card title="Workspace Operations" subtitle="Create workspaces and manage memberships/projects.">
              <div className="icon-row"><LayoutDashboard size={18} /> <span>Workspace + Project orchestration</span></div>
              <form className="form" onSubmit={submit((v) => callApi(session, '/workspaces', 'POST', { workspaceKey: v.workspaceKey, displayName: v.displayName, ownerUserId: v.ownerUserId }))}>
                <InputField label="Workspace Key" name="workspaceKey" />
                <InputField label="Display Name" name="displayName" />
                <InputField label="Owner User ID" name="ownerUserId" />
                <button className="btn primary">Create Workspace</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/workspaces/${v.workspaceId}/memberships`, 'POST', { actorUserId: v.actorUserId, targetUserId: v.targetUserId, targetRole: v.targetRole, invitedByUserId: v.invitedByUserId }))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Actor User ID" name="actorUserId" />
                <InputField label="Target User ID" name="targetUserId" />
                <InputField label="Invited By" name="invitedByUserId" />
                <SelectField label="Target Role" name="targetRole" options={[{ value: 'viewer', label: 'Viewer' }, { value: 'member', label: 'Member' }, { value: 'admin', label: 'Admin' }, { value: 'owner', label: 'Owner' }]} />
                <button className="btn secondary">Add Membership</button>
              </form>
            </Card>
            <Card title="Project Operations" subtitle="Create and list workspace projects.">
              <form className="form" onSubmit={submit((v) => callApi(session, `/workspaces/${v.workspaceId}/projects`, 'POST', { actorUserId: v.actorUserId, projectKey: v.projectKey, displayName: v.displayName, description: v.description }))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Actor User ID" name="actorUserId" />
                <InputField label="Project Key" name="projectKey" />
                <InputField label="Display Name" name="displayName" />
                <TextareaField label="Description" name="description" required={false} />
                <button className="btn primary">Create Project</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/users/${v.userId}/workspaces?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="User ID" name="userId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn secondary">List User Workspaces</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/workspaces/${v.workspaceId}/projects?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn secondary">List Workspace Projects</button>
              </form>
            </Card>
          </>
        )}

        {route === 'tasks' && (
          <>
            <Card title="Task Management" subtitle="Create tasks and update statuses.">
              <div className="icon-row"><ListChecks size={18} /> <span>Task lifecycle</span></div>
              <form className="form" onSubmit={submit((v) => callApi(session, '/tasks', 'POST', { workspaceId: v.workspaceId, projectId: v.projectId, title: v.title, description: v.description, priority: v.priority, assigneeUserId: v.assigneeUserId, createdByUserId: v.createdByUserId }))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Project ID" name="projectId" />
                <InputField label="Title" name="title" />
                <TextareaField label="Description" name="description" required={false} />
                <SelectField label="Priority" name="priority" options={[{ value: 'low', label: 'Low' }, { value: 'medium', label: 'Medium' }, { value: 'high', label: 'High' }]} />
                <InputField label="Assignee User ID" name="assigneeUserId" />
                <InputField label="Created By User ID" name="createdByUserId" />
                <button className="btn primary">Create Task</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/tasks/${v.taskId}/status`, 'PATCH', { status: v.status }))}>
                <InputField label="Task ID" name="taskId" />
                <SelectField label="Status" name="status" options={[{ value: 'todo', label: 'To Do' }, { value: 'in_progress', label: 'In Progress' }, { value: 'done', label: 'Done' }]} />
                <button className="btn secondary">Update Task Status</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/projects/${v.projectId}/tasks?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="Project ID" name="projectId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn ghost">List Project Tasks</button>
              </form>
            </Card>
          </>
        )}

        {route === 'collab' && (
          <>
            <Card title="Collaboration" subtitle="Comments and activity stream across tasks.">
              <div className="icon-row"><MessageCircle size={18} /> <span>Real-time collaboration controls</span></div>
              <form className="form" onSubmit={submit((v) => callApi(session, '/comments', 'POST', { workspaceId: v.workspaceId, projectId: v.projectId, taskId: v.taskId, authorUserId: v.authorUserId, body: v.body }))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Project ID" name="projectId" />
                <InputField label="Task ID" name="taskId" />
                <InputField label="Author User ID" name="authorUserId" />
                <TextareaField label="Comment Body" name="body" required />
                <button className="btn primary">Create Comment</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/comments/${v.commentId}`, 'PATCH', { actorUserId: v.actorUserId, body: v.body }))}>
                <InputField label="Comment ID" name="commentId" />
                <InputField label="Actor User ID" name="actorUserId" />
                <TextareaField label="Updated Body" name="body" required />
                <button className="btn secondary">Update Comment</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/comments/${v.commentId}`, 'DELETE', { actorUserId: v.actorUserId }))}>
                <InputField label="Comment ID" name="commentId" />
                <InputField label="Actor User ID" name="actorUserId" />
                <button className="btn danger">Delete Comment</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/tasks/${v.taskId}/comments?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="Task ID" name="taskId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn ghost">List Comments</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/tasks/${v.taskId}/activities?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="Task ID" name="taskId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn ghost">List Activities</button>
              </form>
            </Card>
          </>
        )}

        {route === 'admin' && (
          <Card title="Admin Operations Center" subtitle="Production-grade control surface for admin-level monitoring.">
            <ul className="admin-list">
              <li>✔ Manage auth lifecycle (register, login, refresh, sign out).</li>
              <li>✔ Provision workspaces and role-based memberships.</li>
              <li>✔ Bootstrap projects and enforce task workflows.</li>
              <li>✔ Moderate collaboration via comment and activity endpoints.</li>
              <li>✔ Audit API responses in JSON console below.</li>
            </ul>
          </Card>
        )}
      </main>

      <section className="result-panel">
        <h2>API Response Console</h2>
        <pre>{output}</pre>
      </section>
    </div>
  )
}
