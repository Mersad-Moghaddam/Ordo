package workspace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
)

func (service *Service) createWorkspaceWithOwner(requestContext context.Context, workspace domainworkspace.Workspace, ownerUserID string) error {
	if existingWorkspace, existingError := service.workspaceRepository.FindWorkspaceByWorkspaceKey(requestContext, workspace.WorkspaceKey); existingError == nil && existingWorkspace.WorkspaceID != "" {
		return domainworkspace.ErrWorkspaceAlreadyExists
	}
	if createWorkspaceError := service.workspaceRepository.CreateWorkspace(requestContext, workspace); createWorkspaceError != nil {
		if errors.Is(createWorkspaceError, domainworkspace.ErrWorkspaceAlreadyExists) {
			return createWorkspaceError
		}
		return fmt.Errorf("create workspace failure: %w", createWorkspaceError)
	}
	ownerMembership := domainworkspace.WorkspaceMembership{WorkspaceID: workspace.WorkspaceID, UserID: ownerUserID, MembershipRole: domainworkspace.MembershipRoleOwner, InvitedByUserID: ownerUserID, JoinedAt: workspace.CreatedAt, LastUpdatedAt: workspace.CreatedAt}
	if createMembershipError := service.membershipRepository.CreateMembership(requestContext, ownerMembership); createMembershipError != nil {
		return fmt.Errorf("create owner membership failure: %w", createMembershipError)
	}
	return nil
}

func canManageMembership(membershipRole domainworkspace.MembershipRole) bool {
	return membershipRole == domainworkspace.MembershipRoleOwner || membershipRole == domainworkspace.MembershipRoleAdmin
}

func canCreateProject(membershipRole domainworkspace.MembershipRole) bool {
	return membershipRole == domainworkspace.MembershipRoleOwner || membershipRole == domainworkspace.MembershipRoleAdmin
}

func defaultIdentifierFunction() (string, error) {
	randomBytes := make([]byte, 16)
	_, randomError := rand.Read(randomBytes)
	if randomError != nil {
		return "", fmt.Errorf("identifier random read failure: %w", randomError)
	}
	return hex.EncodeToString(randomBytes), nil
}
