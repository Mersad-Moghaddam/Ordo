package collab

type CreateCommentRequest struct {
	WorkspaceID  string `json:"workspaceId"`
	ProjectID    string `json:"projectId"`
	TaskID       string `json:"taskId"`
	AuthorUserID string `json:"authorUserId"`
	Body         string `json:"body"`
}

type UpdateCommentRequest struct {
	ActorUserID string `json:"actorUserId"`
	Body        string `json:"body"`
}

type DeleteCommentRequest struct {
	ActorUserID string `json:"actorUserId"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
