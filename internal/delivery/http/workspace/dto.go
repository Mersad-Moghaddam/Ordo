package workspace

type CreateWorkspaceRequest struct {
	WorkspaceKey string `json:"workspaceKey"`
	DisplayName  string `json:"displayName"`
	OwnerUserID  string `json:"ownerUserId"`
}

type AddMembershipRequest struct {
	ActorUserID     string `json:"actorUserId"`
	TargetUserID    string `json:"targetUserId"`
	TargetRole      string `json:"targetRole"`
	InvitedByUserID string `json:"invitedByUserId"`
}

type CreateProjectRequest struct {
	ActorUserID string `json:"actorUserId"`
	ProjectKey  string `json:"projectKey"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}
