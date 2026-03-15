import { FormEvent, useState } from 'react'
import { useTasks } from './hooks/useTasks'
import type { TaskDraft, TaskPriority, TaskStatus } from './types'

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

export default function App() {
  const { tasks, filters, setFilters, stats, addTask, removeTask, updateTaskStatus } = useTasks()
  const [draft, setDraft] = useState<TaskDraft>(initialDraft)

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!draft.title.trim() || !draft.assignee.trim()) return

    addTask({
      ...draft,
      title: draft.title.trim(),
      assignee: draft.assignee.trim(),
      description: draft.description.trim(),
      dueDate: new Date(draft.dueDate).toISOString(),
    })
    setDraft(initialDraft)
  }

  return (
    <div className="page">
      <header className="topbar">
        <div>
          <p className="eyebrow">Ordo</p>
          <h1>Task Operations Command Center</h1>
          <p className="subtitle">Plan work, focus your team, and close delivery loops from one workspace.</p>
        </div>
      </header>

      <section className="stats-grid">
        <article className="stat-card"><p>Total Tasks</p><h2>{stats.total}</h2></article>
        <article className="stat-card"><p>In Progress</p><h2>{stats.progress}</h2></article>
        <article className="stat-card"><p>Completed</p><h2>{stats.done}</h2></article>
        <article className="stat-card"><p>Overdue</p><h2>{stats.overdue}</h2></article>
      </section>

      <section className="workspace-grid">
        <article className="panel form-panel">
          <h3>Create Task</h3>
          <form onSubmit={handleSubmit}>
            <label>Title</label>
            <input
              required
              value={draft.title}
              onChange={(event) => setDraft((current) => ({ ...current, title: event.target.value }))}
              placeholder="Launch team dashboard"
            />

            <label>Description</label>
            <textarea
              rows={3}
              value={draft.description}
              onChange={(event) => setDraft((current) => ({ ...current, description: event.target.value }))}
              placeholder="Define acceptance criteria and rollout steps"
            />

            <div className="form-row">
              <div>
                <label>Assignee</label>
                <input
                  required
                  value={draft.assignee}
                  onChange={(event) => setDraft((current) => ({ ...current, assignee: event.target.value }))}
                  placeholder="Alex Morgan"
                />
              </div>
              <div>
                <label>Due Date</label>
                <input
                  type="date"
                  value={draft.dueDate.slice(0, 10)}
                  onChange={(event) => setDraft((current) => ({ ...current, dueDate: event.target.value }))}
                />
              </div>
            </div>

            <div className="form-row">
              <div>
                <label>Priority</label>
                <select
                  value={draft.priority}
                  onChange={(event) => setDraft((current) => ({ ...current, priority: event.target.value as TaskPriority }))}
                >
                  {priorityOptions.map((priority) => (
                    <option key={priority} value={priority}>{priority}</option>
                  ))}
                </select>
              </div>
              <div>
                <label>Status</label>
                <select
                  value={draft.status}
                  onChange={(event) => setDraft((current) => ({ ...current, status: event.target.value as TaskStatus }))}
                >
                  {statusOptions.map((status) => (
                    <option key={status} value={status}>{prettyStatus(status)}</option>
                  ))}
                </select>
              </div>
            </div>

            <button type="submit">Add Task</button>
          </form>
        </article>

        <article className="panel">
          <div className="toolbar">
            <h3>Backlog</h3>
            <input
              placeholder="Search title, assignee, description"
              value={filters.query}
              onChange={(event) => setFilters((current) => ({ ...current, query: event.target.value }))}
            />
          </div>

          <div className="filters">
            <select
              value={filters.status}
              onChange={(event) => setFilters((current) => ({ ...current, status: event.target.value as TaskStatus | 'all' }))}
            >
              <option value="all">All status</option>
              {statusOptions.map((status) => <option key={status} value={status}>{prettyStatus(status)}</option>)}
            </select>
            <select
              value={filters.priority}
              onChange={(event) => setFilters((current) => ({ ...current, priority: event.target.value as TaskPriority | 'all' }))}
            >
              <option value="all">All priority</option>
              {priorityOptions.map((priority) => <option key={priority} value={priority}>{priority}</option>)}
            </select>
          </div>

          <div className="task-list">
            {tasks.length === 0 ? (
              <p className="empty">No tasks match your current filters.</p>
            ) : (
              tasks.map((task) => (
                <article key={task.id} className="task-card">
                  <div className="task-head">
                    <div>
                      <h4>{task.title}</h4>
                      <p>{task.description || 'No description provided.'}</p>
                    </div>
                    <span className={`pill ${task.priority}`}>{task.priority}</span>
                  </div>
                  <div className="task-meta">
                    <span>Assignee: {task.assignee}</span>
                    <span>Due: {new Date(task.dueDate).toLocaleDateString()}</span>
                  </div>
                  <div className="task-actions">
                    <select
                      value={task.status}
                      onChange={(event) => updateTaskStatus(task.id, event.target.value as TaskStatus)}
                    >
                      {statusOptions.map((status) => <option key={status} value={status}>{prettyStatus(status)}</option>)}
                    </select>
                    <button className="danger" type="button" onClick={() => removeTask(task.id)}>Delete</button>
                  </div>
                </article>
              ))
            )}
          </div>
        </article>
      </section>
    </div>
  )
}
