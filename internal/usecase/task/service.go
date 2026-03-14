package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domaintask "github.com/ordo/backend/internal/domain/task"
	"github.com/ordo/backend/internal/repository"
	repositorytask "github.com/ordo/backend/internal/repository/task"
	"github.com/ordo/backend/internal/usecase"
)

type Service struct {
	taskRepository     repositorytask.TaskRepository
	outboxRepository   repositorytask.OutboxRepository
	transactionManager repository.TransactionManager
	nowFunction        func() time.Time
	identifierFunction func() (string, error)
}

type Option func(service *Service)

type CreateTaskInput struct {
	WorkspaceID     string
	ProjectID       string
	Title           string
	Description     string
	Priority        domaintask.TaskPriority
	AssigneeUserID  *string
	CreatedByUserID string
}

func NewService(
	taskRepository repositorytask.TaskRepository,
	outboxRepository repositorytask.OutboxRepository,
	transactionManager repository.TransactionManager,
	options ...Option,
) *Service {
	service := &Service{
		taskRepository:     taskRepository,
		outboxRepository:   outboxRepository,
		transactionManager: transactionManager,
		nowFunction:        time.Now,
		identifierFunction: defaultIdentifierFunction,
	}
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}
	return service
}

func WithNowFunction(nowFunction func() time.Time) Option {
	return func(service *Service) {
		if nowFunction != nil {
			service.nowFunction = nowFunction
		}
	}
}

func WithIdentifierFunction(identifierFunction func() (string, error)) Option {
	return func(service *Service) {
		if identifierFunction != nil {
			service.identifierFunction = identifierFunction
		}
	}
}

func (service *Service) CreateTask(requestContext context.Context, createTaskInput CreateTaskInput) (domaintask.Task, error) {
	creationTime := service.nowFunction()
	taskID, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domaintask.Task{}, fmt.Errorf("task identifier generation failure: %w", identifierError)
	}
	task := domaintask.Task{TaskID: taskID, WorkspaceID: createTaskInput.WorkspaceID, ProjectID: createTaskInput.ProjectID, Title: createTaskInput.Title, Description: createTaskInput.Description, Status: domaintask.TaskStatusTodo, Priority: createTaskInput.Priority, AssigneeUserID: createTaskInput.AssigneeUserID, CreatedByUserID: createTaskInput.CreatedByUserID, CreatedAt: creationTime, UpdatedAt: creationTime}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		return service.createTaskAndOutbox(transactionContext, task)
	})
	if transactionError != nil {
		return domaintask.Task{}, transactionError
	}
	return task, nil
}

func (service *Service) UpdateTaskStatus(requestContext context.Context, taskID string, nextStatus domaintask.TaskStatus) error {
	existingTask, findError := service.taskRepository.FindTaskByTaskID(requestContext, taskID)
	if findError != nil {
		return findError
	}
	if !isTransitionValid(existingTask.Status, nextStatus) {
		return domaintask.ErrInvalidTaskStatusTransition
	}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		if updateError := service.taskRepository.UpdateTaskStatus(transactionContext, taskID, nextStatus); updateError != nil {
			return updateError
		}
		return service.persistStatusOutbox(transactionContext, existingTask, nextStatus)
	})
	if transactionError != nil {
		return transactionError
	}
	return nil
}

func (service *Service) ListProjectTasks(requestContext context.Context, projectID string, pageNumber int, pageSize int) (usecase.PageResult[domaintask.Task], error) {
	taskList, totalCount, listError := service.taskRepository.ListTasksByProjectID(requestContext, projectID, pageNumber, pageSize)
	if listError != nil {
		return usecase.PageResult[domaintask.Task]{}, listError
	}
	return usecase.PageResult[domaintask.Task]{Items: taskList, TotalCount: totalCount, PageNumber: pageNumber, PageSize: pageSize}, nil
}

func (service *Service) createTaskAndOutbox(requestContext context.Context, task domaintask.Task) error {
	if createError := service.taskRepository.CreateTask(requestContext, task); createError != nil {
		return createError
	}
	payloadMap := map[string]any{"taskId": task.TaskID, "workspaceId": task.WorkspaceID, "projectId": task.ProjectID, "title": task.Title, "status": task.Status, "priority": task.Priority, "createdByUserId": task.CreatedByUserID}
	payloadBytes, payloadError := json.Marshal(payloadMap)
	if payloadError != nil {
		return fmt.Errorf("task create payload encoding failure: %w", payloadError)
	}
	outboxEvent, outboxError := service.newOutboxEvent(task.TaskID, "task.created", string(payloadBytes))
	if outboxError != nil {
		return outboxError
	}
	if persistError := service.outboxRepository.CreateOutboxEvent(requestContext, outboxEvent); persistError != nil {
		return domaintask.ErrOutboxPersistFailure
	}
	return nil
}

func (service *Service) persistStatusOutbox(requestContext context.Context, existingTask domaintask.Task, nextStatus domaintask.TaskStatus) error {
	payloadMap := map[string]any{"taskId": existingTask.TaskID, "projectId": existingTask.ProjectID, "previousStatus": existingTask.Status, "nextStatus": nextStatus}
	payloadBytes, payloadError := json.Marshal(payloadMap)
	if payloadError != nil {
		return fmt.Errorf("task status payload encoding failure: %w", payloadError)
	}
	outboxEvent, outboxError := service.newOutboxEvent(existingTask.TaskID, "task.status.updated", string(payloadBytes))
	if outboxError != nil {
		return outboxError
	}
	if persistError := service.outboxRepository.CreateOutboxEvent(requestContext, outboxEvent); persistError != nil {
		return domaintask.ErrOutboxPersistFailure
	}
	return nil
}

func (service *Service) newOutboxEvent(aggregateID string, eventType string, payload string) (domaintask.OutboxEvent, error) {
	eventIdentifier, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domaintask.OutboxEvent{}, fmt.Errorf("outbox identifier generation failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	return domaintask.OutboxEvent{EventID: eventIdentifier, AggregateType: "task", AggregateID: aggregateID, EventType: eventType, Payload: payload, Status: "pending", Attempts: 0, IdempotencyKey: aggregateID + ":" + eventType + ":" + creationTime.Format(time.RFC3339Nano), CreatedAt: creationTime, UpdatedAt: creationTime}, nil
}

func isTransitionValid(currentStatus domaintask.TaskStatus, nextStatus domaintask.TaskStatus) bool {
	if currentStatus == nextStatus {
		return true
	}
	if currentStatus == domaintask.TaskStatusTodo && (nextStatus == domaintask.TaskStatusInProgress || nextStatus == domaintask.TaskStatusDone) {
		return true
	}
	if currentStatus == domaintask.TaskStatusInProgress && nextStatus == domaintask.TaskStatusDone {
		return true
	}
	return false
}
