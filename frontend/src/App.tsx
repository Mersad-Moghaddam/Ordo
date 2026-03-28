import { FormEvent, useEffect, useState } from 'react'
import { ShieldCheck, LayoutDashboard, ListChecks, MessageCircle } from 'lucide-react'
import { callApi } from './lib/api'
import { SessionState } from './lib/types'
import { Card } from './components/Card'
import { InputField, SelectField, TextareaField } from './components/FormControls'
import { TopNav } from './layout/TopNav'

type AppPath = '/login' | '/signup' | '/main'
const DEFAULT_BASE = 'http://localhost:8080/api/v1'

function data(form: HTMLFormElement): Record<string, string> {
  return Object.fromEntries(new FormData(form).entries()) as Record<string, string>
}

function readPath(): AppPath {
  const current = window.location.pathname
  if (current === '/signup') return '/signup'
  if (current === '/main') return '/main'
  return '/login'
}

function useAppPath() {
  const [path, setPath] = useState<AppPath>(readPath())

  const navigate = (next: AppPath) => {
    window.history.pushState({}, '', next)
    setPath(next)
  }

  useEffect(() => {
    const onPop = () => setPath(readPath())
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  return { path, navigate }
}

export function App() {
  const [session, setSession] = useState<SessionState>({
    apiBase: localStorage.getItem('apiBase') ?? DEFAULT_BASE,
    accessToken: localStorage.getItem('accessToken') ?? '',
    refreshToken: localStorage.getItem('refreshToken') ?? '',
  })
  const [output, setOutput] = useState('Ready.')
  const { path, navigate } = useAppPath()

  const saveSession = (next: SessionState) => {
    setSession(next)
    localStorage.setItem('apiBase', next.apiBase)
    localStorage.setItem('accessToken', next.accessToken)
    localStorage.setItem('refreshToken', next.refreshToken)
  }

  const submit =
    (builder: (values: Record<string, string>) => Promise<unknown>) => async (event: FormEvent<HTMLFormElement>) => {
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
        apiBase={session.apiBase}
        onApiBaseChange={(value) => saveSession({ ...session, apiBase: value.trim() })}
        onNavigate={navigate}
      />

      <main className="content">
        {path === '/login' && (
          <Card title="Login" subtitle="Sign in to access Ordo workspace control center.">
            <form
              className="form"
              onSubmit={submit(async (values) => {
                const result = (await callApi(session, '/auth/login', 'POST', {
                  email: values.email,
                  password: values.password,
                })) as { accessToken: string; refreshToken: string }
                saveSession({ ...session, accessToken: result.accessToken, refreshToken: result.refreshToken })
                navigate('/main')
                return result
              })}
            >
              <InputField label="Email" name="email" type="email" />
              <InputField label="Password" name="password" type="password" />
              <button className="btn primary">Login</button>
            </form>
            <button className="btn ghost" onClick={() => navigate('/signup')}>Create Account</button>
          </Card>
        )}

        {path === '/signup' && (
          <Card title="Create Account" subtitle="Register and automatically continue to main page.">
            <form
              className="form"
              onSubmit={submit(async (values) => {
                const result = (await callApi(session, '/auth/register', 'POST', {
                  email: values.email,
                  password: values.password,
                  role: values.role,
                })) as { accessToken: string; refreshToken: string }
                saveSession({ ...session, accessToken: result.accessToken, refreshToken: result.refreshToken })
                navigate('/main')
                return result
              })}
            >
              <InputField label="Email" name="email" type="email" />
              <InputField label="Password" name="password" type="password" />
              <SelectField label="Role" name="role" options={[{ value: 'member', label: 'Member' }, { value: 'admin', label: 'Admin' }]} />
              <button className="btn primary">Create Account</button>
            </form>
            <button className="btn ghost" onClick={() => navigate('/login')}>Back to Login</button>
          </Card>
        )}

        {path === '/main' && (
          <>
            <Card title="Session" subtitle="Current auth lifecycle">
              <div className="icon-row"><ShieldCheck size={18} /> <span>{session.accessToken ? 'Authenticated' : 'Not Authenticated'}</span></div>
              <form className="form" onSubmit={submit(async () => {
                const result = (await callApi(session, '/auth/refresh', 'POST', { refreshToken: session.refreshToken })) as { accessToken: string; refreshToken: string }
                saveSession({ ...session, accessToken: result.accessToken, refreshToken: result.refreshToken })
                return result
              })}>
                <button className="btn secondary">Refresh Token</button>
              </form>
            </Card>

            <Card title="Workspace + Project" subtitle="All workspace and project routes">
              <div className="icon-row"><LayoutDashboard size={18} /> <span>Workspaces, memberships, projects</span></div>
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
                <InputField label="Invited By User ID" name="invitedByUserId" />
                <SelectField label="Target Role" name="targetRole" options={[{ value: 'viewer', label: 'Viewer' }, { value: 'member', label: 'Member' }, { value: 'admin', label: 'Admin' }, { value: 'owner', label: 'Owner' }]} />
                <button className="btn secondary">Add Membership</button>
              </form>
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
                <button className="btn ghost">List User Workspaces</button>
              </form>
              <form className="form" onSubmit={submit((v) => callApi(session, `/workspaces/${v.workspaceId}/projects?page=${v.page}&pageSize=${v.pageSize}`))}>
                <InputField label="Workspace ID" name="workspaceId" />
                <InputField label="Page" name="page" type="number" defaultValue="1" />
                <InputField label="Page Size" name="pageSize" type="number" defaultValue="20" />
                <button className="btn ghost">List Workspace Projects</button>
              </form>
            </Card>

            <Card title="Tasks" subtitle="Create/update/list task routes">
              <div className="icon-row"><ListChecks size={18} /> <span>Task lifecycle operations</span></div>
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

            <Card title="Collaboration + Admin" subtitle="Comments, activity feed, moderation actions">
              <div className="icon-row"><MessageCircle size={18} /> <span>Collaboration and moderation</span></div>
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
      </main>

      <section className="result-panel">
        <h2>API Response Console</h2>
        <pre>{output}</pre>
      </section>
    </div>
  )
}
