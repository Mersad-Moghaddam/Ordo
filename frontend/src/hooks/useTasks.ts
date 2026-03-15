import { useMemo, useState } from 'react'
import type { TaskDraft, TaskFilters, TaskItem, TaskStatus } from '../types'

const STORAGE_KEY = 'ordo.tasks.v2'

function createSeedTasks(): TaskItem[] {
  const now = new Date().toISOString()
  return [
    {
      id: crypto.randomUUID(),
      title: 'Finalize Sprint Scope',
      description: 'Review sprint backlog and lock scope with product + engineering.',
      assignee: 'Maya Patel',
      dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24).toISOString(),
      priority: 'high',
      status: 'in_progress',
      createdAt: now,
    },
    {
      id: crypto.randomUUID(),
      title: 'Improve task detail UX',
      description: 'Add quick actions and cleaner visual hierarchy for task cards.',
      assignee: 'Noah Kim',
      dueDate: new Date(Date.now() + 1000 * 60 * 60 * 24 * 3).toISOString(),
      priority: 'medium',
      status: 'todo',
      createdAt: now,
    },
    {
      id: crypto.randomUUID(),
      title: 'Publish release notes',
      description: 'Summarize sprint impact and customer-facing improvements.',
      assignee: 'Olivia Chen',
      dueDate: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
      priority: 'low',
      status: 'todo',
      createdAt: now,
    },
  ]
}

function saveTasks(tasks: TaskItem[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(tasks))
}

function loadInitialTasks(): TaskItem[] {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (!stored) {
    const seeded = createSeedTasks()
    saveTasks(seeded)
    return seeded
  }

  try {
    const parsed = JSON.parse(stored) as unknown
    if (!Array.isArray(parsed)) {
      const seeded = createSeedTasks()
      saveTasks(seeded)
      return seeded
    }

    return parsed as TaskItem[]
  } catch {
    const seeded = createSeedTasks()
    saveTasks(seeded)
    return seeded
  }
}

export function useTasks() {
  const [tasks, setTasks] = useState<TaskItem[]>(() => loadInitialTasks())
  const [filters, setFilters] = useState<TaskFilters>({ query: '', status: 'all', priority: 'all' })

  function persist(next: TaskItem[]) {
    setTasks(next)
    saveTasks(next)
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

  function updateTaskStatus(id: string, status: TaskStatus) {
    persist(tasks.map((task) => (task.id === id ? { ...task, status } : task)))
  }

  function moveTask(id: string, direction: 'forward' | 'backward') {
    const order: TaskStatus[] = ['todo', 'in_progress', 'done']
    persist(
      tasks.map((task) => {
        if (task.id !== id) return task
        const currentIndex = order.indexOf(task.status)
        const nextIndex =
          direction === 'forward'
            ? Math.min(currentIndex + 1, order.length - 1)
            : Math.max(currentIndex - 1, 0)
        return { ...task, status: order[nextIndex] }
      }),
    )
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
    const todo = tasks.filter((task) => task.status === 'todo').length
    const overdue = tasks.filter(
      (task) => task.status !== 'done' && new Date(task.dueDate).getTime() < Date.now(),
    ).length

    return { total, done, progress, overdue, todo }
  }, [tasks])

  const grouped = useMemo(() => {
    return {
      todo: filteredTasks.filter((task) => task.status === 'todo'),
      in_progress: filteredTasks.filter((task) => task.status === 'in_progress'),
      done: filteredTasks.filter((task) => task.status === 'done'),
    }
  }, [filteredTasks])

  return {
    tasks: filteredTasks,
    grouped,
    filters,
    setFilters,
    stats,
    addTask,
    removeTask,
    moveTask,
    updateTaskStatus,
  }
}
