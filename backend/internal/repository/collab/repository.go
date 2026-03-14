package collab

import (
	"context"

	domaincollab "github.com/ordo/backend/internal/domain/collab"
)

type CommentRepository interface {
	CreateComment(requestContext context.Context, comment domaincollab.Comment) error
	FindCommentByCommentID(requestContext context.Context, commentID string) (domaincollab.Comment, error)
	UpdateCommentBody(requestContext context.Context, commentID string, body string) error
	SoftDeleteComment(requestContext context.Context, commentID string) error
	ListCommentsByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.Comment, int64, error)
}

type ActivityRepository interface {
	CreateActivity(requestContext context.Context, activityLog domaincollab.ActivityLog) error
	ListActivitiesByTaskID(requestContext context.Context, taskID string, pageNumber int, pageSize int) ([]domaincollab.ActivityLog, int64, error)
}
