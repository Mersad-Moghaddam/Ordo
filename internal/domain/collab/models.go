package collab

import "time"

type Comment struct {
	CommentID    string
	WorkspaceID  string
	ProjectID    string
	TaskID       string
	AuthorUserID string
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type ActivityType string

const (
	ActivityTypeCommentCreated ActivityType = "comment.created"
	ActivityTypeCommentUpdated ActivityType = "comment.updated"
	ActivityTypeCommentDeleted ActivityType = "comment.deleted"
)

type ActivityLog struct {
	ActivityID   string
	WorkspaceID  string
	ProjectID    string
	TaskID       string
	ActorUserID  string
	ActivityType ActivityType
	Payload      string
	CreatedAt    time.Time
}
