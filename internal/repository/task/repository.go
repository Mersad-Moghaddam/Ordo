package task

import (
	"context"

	domaintask "github.com/ordo/backend/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(requestContext context.Context, task domaintask.Task) error
	FindTaskByTaskID(requestContext context.Context, taskID string) (domaintask.Task, error)
	UpdateTaskStatus(requestContext context.Context, taskID string, taskStatus domaintask.TaskStatus) error
	ListTasksByProjectID(requestContext context.Context, projectID string, pageNumber int, pageSize int) ([]domaintask.Task, int64, error)
}

type OutboxRepository interface {
	CreateOutboxEvent(requestContext context.Context, event domaintask.OutboxEvent) error
	ListPendingOutboxEvents(requestContext context.Context, batchSize int) ([]domaintask.OutboxEvent, error)
	MarkOutboxEventPublished(requestContext context.Context, eventID string) error
	MarkOutboxEventRetry(requestContext context.Context, eventID string, attempts int, nextRetryUnixTimestamp int64) error
}
