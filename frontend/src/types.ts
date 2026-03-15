export type TaskStatus = 'todo' | 'in_progress' | 'done'
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical'

export type TaskItem = {
  id: string
  title: string
  description: string
  assignee: string
  dueDate: string
  status: TaskStatus
  priority: TaskPriority
  createdAt: string
}

export type TaskDraft = Omit<TaskItem, 'id' | 'createdAt'>

export type TaskFilters = {
  query: string
  status: TaskStatus | 'all'
  priority: TaskPriority | 'all'
}
