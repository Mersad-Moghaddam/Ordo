import { useMemo, useState } from 'react'
import type { TaskDraft, TaskFilters, TaskItem } from '../types'

const STORAGE_KEY = 'ordo.tasks'

function createSeedTasks(): TaskItem[] {
  const now = new Date().toISOString()
  return [
    {
      id: crypto.randomUUID(),
      title: 'Define MVP delivery plan',
      description: 'Break Ordo launch into milestones for auth, workspace, and collaboration.',
      assignee: 'Product Team',
      dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24 * 2).toISOString(),
      priority: 'high',
      status: 'in_progress',
      createdAt: now,
    },
    {
      id: crypto.randomUUID(),
      title: 'Ship onboarding flow',
      description: 'Improve first-run guidance with empty states, hints, and quick actions.',
      assignee: 'Frontend Team',
      dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24 * 5).toISOString(),
      priority: 'medium',
      status: 'todo',
      createdAt: now,
    },
  ]
}

function loadInitialTasks(): TaskItem[] {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (!stored) {
    const seeded = createSeedTasks()
    localStorage.setItem(STORAGE_KEY, JSON.stringify(seeded))
    return seeded
  }

  try {
    const parsed = JSON.parse(stored) as TaskItem[]
    return Array.isArray(parsed) ? parsed : createSeedTasks()
  } catch {
    return createSeedTasks()
  }
}

export function useTasks() {
  const [tasks, setTasks] = useState<TaskItem[]>(() => loadInitialTasks())
  const [filters, setFilters] = useState<TaskFilters>({ query: '', status: 'all', priority: 'all' })

  function persist(next: TaskItem[]) {
    setTasks(next)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
  }

  function addTask(draft: TaskDraft) {
    persist([
      {
        id: crypto.randomUUID(),
        createdAt: new Date().toISOString(),
        ...draft,
      },
      ...tasks,
    ])
  }

  function updateTaskStatus(id: string, status: TaskItem['status']) {
    persist(tasks.map((task) => (task.id === id ? { ...task, status } : task)))
  }

  function removeTask(id: string) {
    persist(tasks.filter((task) => task.id !== id))
  }

  const filteredTasks = useMemo(() => {
    return tasks.filter((task) => {
      const queryMatched = [task.title, task.description, task.assignee]
        .join(' ')
        .toLowerCase()
        .includes(filters.query.trim().toLowerCase())
      const statusMatched = filters.status === 'all' || task.status === filters.status
      const priorityMatched = filters.priority === 'all' || task.priority === filters.priority
      return queryMatched && statusMatched && priorityMatched
    })
  }, [filters, tasks])

  const stats = useMemo(() => {
    const total = tasks.length
    const done = tasks.filter((task) => task.status === 'done').length
    const progress = tasks.filter((task) => task.status === 'in_progress').length
    const overdue = tasks.filter(
      (task) => task.status !== 'done' && new Date(task.dueDate).getTime() < Date.now(),
    ).length

    return { total, done, progress, overdue }
  }, [tasks])

  return {
    tasks: filteredTasks,
    filters,
    setFilters,
    stats,
    addTask,
    removeTask,
    updateTaskStatus,
  }
}
