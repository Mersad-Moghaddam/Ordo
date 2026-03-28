package collab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domaincollab "github.com/ordo/backend/internal/domain/collab"
	domaintask "github.com/ordo/backend/internal/domain/task"
	"github.com/ordo/backend/internal/repository"
	repositorycollab "github.com/ordo/backend/internal/repository/collab"
	repositorytask "github.com/ordo/backend/internal/repository/task"
	"github.com/ordo/backend/internal/usecase"
)

type Service struct {
	commentRepository  repositorycollab.CommentRepository
	activityRepository repositorycollab.ActivityRepository
	outboxRepository   repositorytask.OutboxRepository
	transactionManager repository.TransactionManager
	nowFunction        func() time.Time
	identifierFunction func() (string, error)
}

type Option func(service *Service)

type CreateCommentInput struct {
	WorkspaceID  string
	ProjectID    string
	TaskID       string
	AuthorUserID string
	Body         string
}

func NewService(
	commentRepository repositorycollab.CommentRepository,
	activityRepository repositorycollab.ActivityRepository,
	outboxRepository repositorytask.OutboxRepository,
	transactionManager repository.TransactionManager,
	options ...Option,
) *Service {
	service := &Service{commentRepository: commentRepository, activityRepository: activityRepository, outboxRepository: outboxRepository, transactionManager: transactionManager, nowFunction: time.Now, identifierFunction: defaultIdentifierFunction}
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

func (service *Service) CreateComment(requestContext context.Context, createInput CreateCommentInput) (domaincollab.Comment, error) {
	commentIdentifier, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domaincollab.Comment{}, fmt.Errorf("comment identifier generation failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	comment := domaincollab.Comment{CommentID: commentIdentifier, WorkspaceID: createInput.WorkspaceID, ProjectID: createInput.ProjectID, TaskID: createInput.TaskID, AuthorUserID: createInput.AuthorUserID, Body: createInput.Body, CreatedAt: creationTime, UpdatedAt: creationTime}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		return service.persistCommentCreated(transactionContext, comment)
	})
	if transactionError != nil {
		return domaincollab.Comment{}, transactionError
	}
	return comment, nil
}

func (service *Service) UpdateComment(requestContext context.Context, commentID string, actorUserID string, body string) error {
	existingComment, findError := service.commentRepository.FindCommentByCommentID(requestContext, commentID)
	if findError != nil {
		return findError
	}
	if existingComment.AuthorUserID != actorUserID {
		return domaincollab.ErrCommentForbidden
	}
	if existingComment.DeletedAt != nil {
		return domaincollab.ErrCommentAlreadyDeleted
	}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		if updateError := service.commentRepository.UpdateCommentBody(transactionContext, commentID, body); updateError != nil {
			return updateError
		}
		return service.persistActivityAndOutbox(transactionContext, existingComment, domaincollab.ActivityTypeCommentUpdated, map[string]any{"commentId": commentID, "body": body, "actorUserId": actorUserID})
	})
	return transactionError
}

func (service *Service) DeleteComment(requestContext context.Context, commentID string, actorUserID string) error {
	existingComment, findError := service.commentRepository.FindCommentByCommentID(requestContext, commentID)
	if findError != nil {
		return findError
	}
	if existingComment.AuthorUserID != actorUserID {
		return domaincollab.ErrCommentForbidden
	}
	if existingComment.DeletedAt != nil {
		return domaincollab.ErrCommentAlreadyDeleted
	}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		if deleteError := service.commentRepository.SoftDeleteComment(transactionContext, commentID); deleteError != nil {
			return deleteError
		}
		return service.persistActivityAndOutbox(transactionContext, existingComment, domaincollab.ActivityTypeCommentDeleted, map[string]any{"commentId": commentID, "actorUserId": actorUserID})
	})
	return transactionError
}

func (service *Service) ListTaskComments(requestContext context.Context, taskID string, pageNumber int, pageSize int) (usecase.PageResult[domaincollab.Comment], error) {
	commentList, totalCount, listError := service.commentRepository.ListCommentsByTaskID(requestContext, taskID, pageNumber, pageSize)
	if listError != nil {
		return usecase.PageResult[domaincollab.Comment]{}, listError
	}
	return usecase.PageResult[domaincollab.Comment]{Items: commentList, TotalCount: totalCount, PageNumber: pageNumber, PageSize: pageSize}, nil
}

func (service *Service) ListTaskActivities(requestContext context.Context, taskID string, pageNumber int, pageSize int) (usecase.PageResult[domaincollab.ActivityLog], error) {
	activityList, totalCount, listError := service.activityRepository.ListActivitiesByTaskID(requestContext, taskID, pageNumber, pageSize)
	if listError != nil {
		return usecase.PageResult[domaincollab.ActivityLog]{}, listError
	}
	return usecase.PageResult[domaincollab.ActivityLog]{Items: activityList, TotalCount: totalCount, PageNumber: pageNumber, PageSize: pageSize}, nil
}
func (service *Service) persistCommentCreated(requestContext context.Context, comment domaincollab.Comment) error {
	if createError := service.commentRepository.CreateComment(requestContext, comment); createError != nil {
		return createError
	}
	payloadMap := map[string]any{"commentId": comment.CommentID, "taskId": comment.TaskID, "authorUserId": comment.AuthorUserID, "body": comment.Body}
	return service.persistActivityAndOutbox(requestContext, comment, domaincollab.ActivityTypeCommentCreated, payloadMap)
}

func (service *Service) persistActivityAndOutbox(requestContext context.Context, comment domaincollab.Comment, activityType domaincollab.ActivityType, payloadMap map[string]any) error {
	payloadBytes, payloadError := json.Marshal(payloadMap)
	if payloadError != nil {
		return fmt.Errorf("payload encoding failure: %w", payloadError)
	}
	activityIdentifier, activityIdentifierError := service.identifierFunction()
	if activityIdentifierError != nil {
		return fmt.Errorf("activity identifier generation failure: %w", activityIdentifierError)
	}
	activityLog := domaincollab.ActivityLog{ActivityID: activityIdentifier, WorkspaceID: comment.WorkspaceID, ProjectID: comment.ProjectID, TaskID: comment.TaskID, ActorUserID: comment.AuthorUserID, ActivityType: activityType, Payload: string(payloadBytes), CreatedAt: service.nowFunction()}
	if createError := service.activityRepository.CreateActivity(requestContext, activityLog); createError != nil {
		return domaincollab.ErrActivityLogWriteFailure
	}
	outboxEvent, outboxError := service.newOutboxEvent(comment.TaskID, string(activityType), string(payloadBytes))
	if outboxError != nil {
		return outboxError
	}
	if createError := service.outboxRepository.CreateOutboxEvent(requestContext, outboxEvent); createError != nil {
		return domaintask.ErrOutboxPersistFailure
	}
	return nil
}

func (service *Service) newOutboxEvent(taskID string, eventType string, payload string) (domaintask.OutboxEvent, error) {
	eventIdentifier, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domaintask.OutboxEvent{}, fmt.Errorf("outbox identifier generation failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	return domaintask.OutboxEvent{EventID: eventIdentifier, AggregateType: "comment", AggregateID: taskID, EventType: eventType, Payload: payload, Status: "pending", Attempts: 0, IdempotencyKey: taskID + ":" + eventType + ":" + creationTime.Format(time.RFC3339Nano), CreatedAt: creationTime, UpdatedAt: creationTime}, nil
}

func ensureCommentError(applicationError error) error {
	if errors.Is(applicationError, domaincollab.ErrCommentNotFound) {
		return domaincollab.ErrCommentNotFound
	}
	return applicationError
}
