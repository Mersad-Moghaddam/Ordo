import { FormEvent, useMemo, useState } from 'react'
import { useTasks } from './hooks/useTasks'
import type { TaskDraft, TaskPriority, TaskStatus, ViewMode } from './types'

const priorityOptions: TaskPriority[] = ['low', 'medium', 'high', 'critical']
const statusOptions: TaskStatus[] = ['todo', 'in_progress', 'done']

const initialDraft: TaskDraft = {
  title: '',
  description: '',
  assignee: '',
  dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24 * 7).toISOString().slice(0, 10),
  priority: 'medium',
  status: 'todo',
}

function prettyStatus(status: TaskStatus) {
  return status.replace('_', ' ')
}

function dueTone(dueDate: string, status: TaskStatus) {
  if (status === 'done') return 'neutral'
  return new Date(dueDate).getTime() < Date.now() ? 'overdue' : 'upcoming'
}

export default function App() {
  const { tasks, grouped, filters, setFilters, stats, addTask, removeTask, moveTask, updateTaskStatus } = useTasks()
  const [draft, setDraft] = useState<TaskDraft>(initialDraft)
  const [view, setView] = useState<ViewMode>('board')

  const completion = useMemo(() => {
    if (stats.total === 0) return 0
    return Math.round((stats.done / stats.total) * 100)
  }, [stats.done, stats.total])

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!draft.title.trim() || !draft.assignee.trim()) return

    addTask({
      ...draft,
      title: draft.title.trim(),
      assignee: draft.assignee.trim(),
      description: draft.description.trim(),
      dueDate: new Date(`${draft.dueDate}T12:00:00`).toISOString(),
    })
    setDraft({ ...initialDraft, dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24 * 7).toISOString().slice(0, 10) })
  }

  return (
    <div className="shell">
      <aside className="sidebar">
        <h2>Ordo</h2>
        <p>Task Manager</p>
        <nav>
          <a className="active">Board</a>
          <a>Projects</a>
          <a>Reports</a>
          <a>Settings</a>
        </nav>
      </aside>

      <main className="page">
        <header className="topbar">
          <div>
            <p className="eyebrow">Sprint workspace</p>
            <h1>Jira-style Task Command Center</h1>
            <p className="subtitle">Track work across backlog, active delivery, and done columns.</p>
          </div>
          <div className="view-toggle">
            <button className={view === 'board' ? 'active' : ''} onClick={() => setView('board')} type="button">Board</button>
            <button className={view === 'list' ? 'active' : ''} onClick={() => setView('list')} type="button">List</button>
          </div>
        </header>

        <section className="stats-grid">
          <article className="stat-card"><p>Total</p><h2>{stats.total}</h2></article>
          <article className="stat-card"><p>To Do</p><h2>{stats.todo}</h2></article>
          <article className="stat-card"><p>In Progress</p><h2>{stats.progress}</h2></article>
          <article className="stat-card"><p>Done</p><h2>{stats.done}</h2></article>
          <article className="stat-card warning"><p>Overdue</p><h2>{stats.overdue}</h2></article>
        </section>

        <section className="progress-card">
          <div className="progress-labels">
            <span>Sprint completion</span>
            <strong>{completion}%</strong>
          </div>
          <div className="progress-track"><div className="progress-fill" style={{ width: `${completion}%` }} /></div>
        </section>

        <section className="workspace-grid">
          <article className="panel form-panel">
            <h3>Create task</h3>
            <form onSubmit={handleSubmit}>
              <label>Title</label>
              <input required value={draft.title} onChange={(event) => setDraft((current) => ({ ...current, title: event.target.value }))} placeholder="Implement workspace analytics" />

              <label>Description</label>
              <textarea rows={3} value={draft.description} onChange={(event) => setDraft((current) => ({ ...current, description: event.target.value }))} placeholder="Write scope and acceptance criteria" />

              <div className="form-row">
                <div>
                  <label>Assignee</label>
                  <input required value={draft.assignee} onChange={(event) => setDraft((current) => ({ ...current, assignee: event.target.value }))} placeholder="Jordan Lee" />
                </div>
                <div>
                  <label>Due date</label>
                  <input type="date" value={draft.dueDate.slice(0, 10)} onChange={(event) => setDraft((current) => ({ ...current, dueDate: event.target.value }))} />
                </div>
              </div>

              <div className="form-row">
                <div>
                  <label>Priority</label>
                  <select value={draft.priority} onChange={(event) => setDraft((current) => ({ ...current, priority: event.target.value as TaskPriority }))}>
                    {priorityOptions.map((priority) => <option key={priority} value={priority}>{priority}</option>)}
                  </select>
                </div>
                <div>
                  <label>Status</label>
                  <select value={draft.status} onChange={(event) => setDraft((current) => ({ ...current, status: event.target.value as TaskStatus }))}>
                    {statusOptions.map((status) => <option key={status} value={status}>{prettyStatus(status)}</option>)}
                  </select>
                </div>
              </div>

              <button type="submit">Create task</button>
            </form>
          </article>

          <article className="panel">
            <div className="toolbar">
              <h3>Task explorer</h3>
              <input placeholder="Search tasks, assignees, descriptions" value={filters.query} onChange={(event) => setFilters((current) => ({ ...current, query: event.target.value }))} />
            </div>

            <div className="filters">
              <select value={filters.status} onChange={(event) => setFilters((current) => ({ ...current, status: event.target.value as TaskStatus | 'all' }))}>
                <option value="all">All status</option>
                {statusOptions.map((status) => <option key={status} value={status}>{prettyStatus(status)}</option>)}
              </select>
              <select value={filters.priority} onChange={(event) => setFilters((current) => ({ ...current, priority: event.target.value as TaskPriority | 'all' }))}>
                <option value="all">All priority</option>
                {priorityOptions.map((priority) => <option key={priority} value={priority}>{priority}</option>)}
              </select>
            </div>

            {view === 'board' ? (
              <div className="board">
                {statusOptions.map((status) => (
                  <section key={status} className="column">
                    <header>
                      <h4>{prettyStatus(status)}</h4>
                      <span>{grouped[status].length}</span>
                    </header>
                    <div className="column-list">
                      {grouped[status].length === 0 ? (
                        <p className="empty">No tasks</p>
                      ) : (
                        grouped[status].map((task) => (
                          <article key={task.id} className="task-card">
                            <div className="task-head">
                              <h5>{task.title}</h5>
                              <span className={`pill ${task.priority}`}>{task.priority}</span>
                            </div>
                            <p className="task-text">{task.description || 'No description provided.'}</p>
                            <div className="task-meta">
                              <span>{task.assignee}</span>
                              <span className={dueTone(task.dueDate, task.status)}>{new Date(task.dueDate).toLocaleDateString()}</span>
                            </div>
                            <div className="task-actions">
                              <button type="button" onClick={() => moveTask(task.id, 'backward')}>←</button>
                              <button type="button" onClick={() => moveTask(task.id, 'forward')}>→</button>
                              <button className="danger" type="button" onClick={() => removeTask(task.id)}>Delete</button>
                            </div>
                          </article>
                        ))
                      )}
                    </div>
                  </section>
                ))}
              </div>
            ) : (
              <div className="task-list">
                {tasks.length === 0 ? (
                  <p className="empty">No tasks match your filters.</p>
                ) : (
                  tasks.map((task) => (
                    <article key={task.id} className="task-card list-card">
                      <div className="task-head">
                        <div>
                          <h5>{task.title}</h5>
                          <p className="task-text">{task.description || 'No description provided.'}</p>
                        </div>
                        <span className={`pill ${task.priority}`}>{task.priority}</span>
                      </div>
                      <div className="task-meta">
                        <span>Assignee: {task.assignee}</span>
                        <span className={dueTone(task.dueDate, task.status)}>Due: {new Date(task.dueDate).toLocaleDateString()}</span>
                      </div>
                      <div className="task-actions">
                        <select value={task.status} onChange={(event) => updateTaskStatus(task.id, event.target.value as TaskStatus)}>
                          {statusOptions.map((status) => <option key={status} value={status}>{prettyStatus(status)}</option>)}
                        </select>
                        <button className="danger" type="button" onClick={() => removeTask(task.id)}>Delete</button>
                      </div>
                    </article>
                  ))
                )}
              </div>
            )}
          </article>
        </section>
      </main>
    </div>
  )
}
