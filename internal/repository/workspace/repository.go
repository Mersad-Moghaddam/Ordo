package workspace

import (
	"context"

	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
)

type WorkspaceRepository interface {
	CreateWorkspace(requestContext context.Context, workspace domainworkspace.Workspace) error
	FindWorkspaceByWorkspaceID(requestContext context.Context, workspaceID string) (domainworkspace.Workspace, error)
	FindWorkspaceByWorkspaceKey(requestContext context.Context, workspaceKey string) (domainworkspace.Workspace, error)
	ListWorkspacesByUserID(requestContext context.Context, userID string, pageNumber int, pageSize int) ([]domainworkspace.Workspace, int64, error)
}

type MembershipRepository interface {
	CreateMembership(requestContext context.Context, membership domainworkspace.WorkspaceMembership) error
	FindMembership(requestContext context.Context, workspaceID string, userID string) (domainworkspace.WorkspaceMembership, error)
	UpdateMembershipRole(requestContext context.Context, workspaceID string, userID string, membershipRole domainworkspace.MembershipRole) error
}

type ProjectRepository interface {
	CreateProject(requestContext context.Context, project domainworkspace.Project) error
	FindProjectByProjectID(requestContext context.Context, projectID string) (domainworkspace.Project, error)
	FindProjectByWorkspaceAndProjectKey(requestContext context.Context, workspaceID string, projectKey string) (domainworkspace.Project, error)
	ListProjectsByWorkspaceID(requestContext context.Context, workspaceID string, pageNumber int, pageSize int) ([]domainworkspace.Project, int64, error)
}
