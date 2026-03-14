package task

import "time"

type TaskStatus string

type TaskPriority string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	TaskID          string
	WorkspaceID     string
	ProjectID       string
	Title           string
	Description     string
	Status          TaskStatus
	Priority        TaskPriority
	AssigneeUserID  *string
	CreatedByUserID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OutboxEvent struct {
	EventID        string
	AggregateType  string
	AggregateID    string
	EventType      string
	Payload        string
	Status         string
	Attempts       int
	NextRetryAt    *time.Time
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
