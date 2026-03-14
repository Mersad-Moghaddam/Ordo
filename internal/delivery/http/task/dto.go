package task

type CreateTaskRequest struct {
	WorkspaceID     string  `json:"workspaceId"`
	ProjectID       string  `json:"projectId"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Priority        string  `json:"priority"`
	AssigneeUserID  *string `json:"assigneeUserId"`
	CreatedByUserID string  `json:"createdByUserId"`
}

type UpdateTaskStatusRequest struct {
	TaskStatus string `json:"status"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
