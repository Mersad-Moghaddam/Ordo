package collab

import (
	"context"
	"errors"
	"testing"
	"time"

	domaincollab "github.com/ordo/backend/internal/domain/collab"
	domaintask "github.com/ordo/backend/internal/domain/task"
)

type mockTransactionManager struct{}

func (mockTransactionManager mockTransactionManager) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	return transactionWorkload(requestContext)
}

type mockCommentRepository struct {
	commentsByID map[string]domaincollab.Comment
}

func (mockRepository *mockCommentRepository) CreateComment(requestContext context.Context, comment domaincollab.Comment) error {
	mockRepository.commentsByID[comment.CommentID] = comment
	return nil
}

func (mockRepository *mockCommentRepository) FindCommentByCommentID(requestContext context.Context, commentID string) (domaincollab.Comment, error) {
	comment, hasComment := mockRepository.commentsByID[commentID]
	if !hasComment {
		return domaincollab.Comment{}, domaincollab.ErrCommentNotFound
	}
	return comment, nil
}

func (mockRepository *mockCommentRepository) UpdateCommentBody(requestContext context.Context, commentID string, body string) error {
	comment, hasComment := mockRepository.commentsByID[commentID]
	if !hasComment {
		return domaincollab.ErrCommentNotFound
	}
	comment.Body = body
	mockRepository.commentsByID[commentID] = comment
	return nil
}

func (mockRepository *mockCommentRepository) SoftDeleteComment(requestContext context.Context, commentID string) error {
	comment, hasComment := mockRepository.commentsByID[commentID]
	if !hasComment {
		return domaincollab.ErrCommentNotFound
	}
	nowValue := time.Now()
	comment.DeletedAt = &nowValue
	mockRepository.commentsByID[commentID] = comment
	return nil
}

func (mockRepository *mockCommentRepository) ListCommentsByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.Comment, int64, error) {
	commentList := make([]domaincollab.Comment, 0)
	for _, comment := range mockRepository.commentsByID {
		if comment.TaskID == taskID {
			commentList = append(commentList, comment)
		}
	}
	return commentList, int64(len(commentList)), nil
}

type mockActivityRepository struct {
	activityList []domaincollab.ActivityLog
}

func (mockRepository *mockActivityRepository) CreateActivity(requestContext context.Context, activityLog domaincollab.ActivityLog) error {
	mockRepository.activityList = append(mockRepository.activityList, activityLog)
	return nil
}

func (mockRepository *mockActivityRepository) ListActivitiesByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.ActivityLog, int64, error) {
	filteredActivityList := make([]domaincollab.ActivityLog, 0)
	for _, activityLog := range mockRepository.activityList {
		if activityLog.TaskID == taskID {
			filteredActivityList = append(filteredActivityList, activityLog)
		}
	}
	return filteredActivityList, int64(len(filteredActivityList)), nil
}

type mockOutboxRepository struct {
	outboxEventList []domaintask.OutboxEvent
}

func (mockRepository *mockOutboxRepository) CreateOutboxEvent(requestContext context.Context, event domaintask.OutboxEvent) error {
	mockRepository.outboxEventList = append(mockRepository.outboxEventList, event)
	return nil
}

func (mockRepository *mockOutboxRepository) ListPendingOutboxEvents(requestContext context.Context, batchSize int) ([]domaintask.OutboxEvent, error) {
	return mockRepository.outboxEventList, nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventPublished(requestContext context.Context, eventID string) error {
	return nil
}

func (mockRepository *mockOutboxRepository) MarkOutboxEventRetry(requestContext context.Context, eventID string, attempts int, nextRetryUnixTimestamp int64) error {
	return nil
}

func TestCommentLifecycle(testingSuite *testing.T) {
	commentRepository := &mockCommentRepository{commentsByID: map[string]domaincollab.Comment{}}
	activityRepository := &mockActivityRepository{activityList: make([]domaincollab.ActivityLog, 0)}
	outboxRepository := &mockOutboxRepository{outboxEventList: make([]domaintask.OutboxEvent, 0)}
	service := NewService(commentRepository, activityRepository, outboxRepository, mockTransactionManager{}, WithIdentifierFunction(func() (string, error) { return "identifier-value", nil }))

	createdComment, createError := service.CreateComment(context.Background(), CreateCommentInput{WorkspaceID: "workspace-id", ProjectID: "project-id", TaskID: "task-id", AuthorUserID: "user-id", Body: "comment body"})
	if createError != nil {
		testingSuite.Fatalf("comment creation failure: %v", createError)
	}
	if createdComment.CommentID == "" {
		testingSuite.Fatalf("comment id must not be empty")
	}
	if len(activityRepository.activityList) != 1 || len(outboxRepository.outboxEventList) != 1 {
		testingSuite.Fatalf("create should generate one activity and one outbox event")
	}

	updateError := service.UpdateComment(context.Background(), createdComment.CommentID, "user-id", "updated body")
	if updateError != nil {
		testingSuite.Fatalf("comment update failure: %v", updateError)
	}

	deleteError := service.DeleteComment(context.Background(), createdComment.CommentID, "user-id")
	if deleteError != nil {
		testingSuite.Fatalf("comment delete failure: %v", deleteError)
	}

	commentsPage, commentsError := service.ListTaskComments(context.Background(), "task-id", 1, 20)
	if commentsError != nil || commentsPage.TotalCount != 1 {
		testingSuite.Fatalf("comment listing failure: %v", commentsError)
	}

	activityPage, activityError := service.ListTaskActivities(context.Background(), "task-id", 1, 20)
	if activityError != nil || activityPage.TotalCount != 3 {
		testingSuite.Fatalf("activity listing failure: %v", activityError)
	}
}

func TestCommentForbiddenAndIdentifierFailure(testingSuite *testing.T) {
	commentRepository := &mockCommentRepository{commentsByID: map[string]domaincollab.Comment{"comment-id": {CommentID: "comment-id", TaskID: "task-id", AuthorUserID: "owner-user", Body: "body"}}}
	activityRepository := &mockActivityRepository{activityList: make([]domaincollab.ActivityLog, 0)}
	outboxRepository := &mockOutboxRepository{outboxEventList: make([]domaintask.OutboxEvent, 0)}
	service := NewService(commentRepository, activityRepository, outboxRepository, mockTransactionManager{}, WithIdentifierFunction(func() (string, error) { return "", errors.New("identifier failure") }))

	_, createError := service.CreateComment(context.Background(), CreateCommentInput{WorkspaceID: "workspace-id", ProjectID: "project-id", TaskID: "task-id", AuthorUserID: "owner-user", Body: "body"})
	if createError == nil {
		testingSuite.Fatalf("expected create identifier error")
	}

	updateError := service.UpdateComment(context.Background(), "comment-id", "other-user", "updated")
	if !errors.Is(updateError, domaincollab.ErrCommentForbidden) {
		testingSuite.Fatalf("expected forbidden error, got: %v", updateError)
	}
}
