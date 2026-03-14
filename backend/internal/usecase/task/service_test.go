package task

import (
	"context"
	"errors"
	"testing"
	"time"

	domaintask "github.com/ordo/backend/internal/domain/task"
)

type mockTransactionManager struct{}

func (mockTransactionManager mockTransactionManager) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	return transactionWorkload(requestContext)
}

type mockTaskRepository struct {
	tasksByID map[string]domaintask.Task
}

func (mockRepository *mockTaskRepository) CreateTask(requestContext context.Context, task domaintask.Task) error {
	if _, hasTask := mockRepository.tasksByID[task.TaskID]; hasTask {
		return domaintask.ErrTaskAlreadyExists
	}
	mockRepository.tasksByID[task.TaskID] = task
	return nil
}

func (mockRepository *mockTaskRepository) FindTaskByTaskID(requestContext context.Context, taskID string) (domaintask.Task, error) {
	task, hasTask := mockRepository.tasksByID[taskID]
	if !hasTask {
		return domaintask.Task{}, domaintask.ErrTaskNotFound
	}
	return task, nil
}

func (mockRepository *mockTaskRepository) UpdateTaskStatus(requestContext context.Context, taskID string, taskStatus domaintask.TaskStatus) error {
	task, hasTask := mockRepository.tasksByID[taskID]
	if !hasTask {
		return domaintask.ErrTaskNotFound
	}
	task.Status = taskStatus
	task.UpdatedAt = time.Now()
	mockRepository.tasksByID[taskID] = task
	return nil
}

func (mockRepository *mockTaskRepository) ListTasksByProjectID(requestContext context.Context, projectID string, pageNumber int, pageSize int) ([]domaintask.Task, int64, error) {
	taskList := make([]domaintask.Task, 0)
	for _, task := range mockRepository.tasksByID {
		if task.ProjectID == projectID {
			taskList = append(taskList, task)
		}
	}
	return taskList, int64(len(taskList)), nil
}

type mockOutboxRepository struct {
	outboxEvents []domaintask.OutboxEvent
}

func (mockRepository *mockOutboxRepository) CreateOutboxEvent(requestContext context.Context, event domaintask.OutboxEvent) error {
	mockRepository.outboxEvents = append(mockRepository.outboxEvents, event)
	return nil
}

func (mockRepository *mockOutboxRepository) ListPendingOutboxEvents(requestContext context.Context, batchSize int) ([]domaintask.OutboxEvent, error) {
	return mockRepository.outboxEvents, nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventPublished(requestContext context.Context, eventID string) error {
	return nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventRetry(requestContext context.Context, eventID string, attempts int, nextRetryUnixTimestamp int64) error {
	return nil
}

func TestCreateTaskAndStatusUpdate(testingSuite *testing.T) {
	nowValue := time.Unix(1700000200, 0)
	taskRepository := &mockTaskRepository{tasksByID: map[string]domaintask.Task{}}
	outboxRepository := &mockOutboxRepository{outboxEvents: make([]domaintask.OutboxEvent, 0)}
	service := NewService(taskRepository, outboxRepository, mockTransactionManager{}, WithNowFunction(func() time.Time { return nowValue }), WithIdentifierFunction(func() (string, error) { return "task-identifier", nil }))

	task, createError := service.CreateTask(context.Background(), CreateTaskInput{WorkspaceID: "workspace-id", ProjectID: "project-id", Title: "build task core", Description: "implement phase three", Priority: domaintask.TaskPriorityHigh, CreatedByUserID: "owner-user"})
	if createError != nil {
		testingSuite.Fatalf("create task failure: %v", createError)
	}
	if task.Status != domaintask.TaskStatusTodo {
		testingSuite.Fatalf("expected todo status")
	}
	if len(outboxRepository.outboxEvents) != 1 {
		testingSuite.Fatalf("expected one outbox event after creation")
	}

	statusError := service.UpdateTaskStatus(context.Background(), task.TaskID, domaintask.TaskStatusInProgress)
	if statusError != nil {
		testingSuite.Fatalf("status update failure: %v", statusError)
	}
	if len(outboxRepository.outboxEvents) != 2 {
		testingSuite.Fatalf("expected outbox event for status update")
	}

	listResult, listError := service.ListProjectTasks(context.Background(), "project-id", 1, 20)
	if listError != nil {
		testingSuite.Fatalf("list tasks failure: %v", listError)
	}
	if listResult.TotalCount != 1 {
		testingSuite.Fatalf("expected one task in listing")
	}
}

func TestInvalidTransitionAndIdentifierFailure(testingSuite *testing.T) {
	taskRepository := &mockTaskRepository{tasksByID: map[string]domaintask.Task{"task-id": {TaskID: "task-id", ProjectID: "project-id", Status: domaintask.TaskStatusDone}}}
	outboxRepository := &mockOutboxRepository{outboxEvents: make([]domaintask.OutboxEvent, 0)}
	service := NewService(taskRepository, outboxRepository, mockTransactionManager{}, WithIdentifierFunction(func() (string, error) { return "", errors.New("identifier failure") }))

	_, createError := service.CreateTask(context.Background(), CreateTaskInput{WorkspaceID: "workspace-id", ProjectID: "project-id", Title: "build", Description: "desc", Priority: domaintask.TaskPriorityLow, CreatedByUserID: "owner-user"})
	if createError == nil {
		testingSuite.Fatalf("expected create identifier failure")
	}

	statusError := service.UpdateTaskStatus(context.Background(), "task-id", domaintask.TaskStatusTodo)
	if !errors.Is(statusError, domaintask.ErrInvalidTaskStatusTransition) {
		testingSuite.Fatalf("expected invalid transition error, got: %v", statusError)
	}
}
