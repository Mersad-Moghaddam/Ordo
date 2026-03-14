package workspace

import (
	"context"
	"errors"
	"testing"
	"time"

	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
)

type mockTransactionManager struct{}

func (mockTransactionManager mockTransactionManager) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	return transactionWorkload(requestContext)
}

type mockWorkspaceRepository struct {
	workspacesByKey  map[string]domainworkspace.Workspace
	workspacesByID   map[string]domainworkspace.Workspace
	workspacesByUser map[string][]domainworkspace.Workspace
}

func (mockRepository *mockWorkspaceRepository) CreateWorkspace(requestContext context.Context, workspace domainworkspace.Workspace) error {
	if _, hasWorkspace := mockRepository.workspacesByKey[workspace.WorkspaceKey]; hasWorkspace {
		return domainworkspace.ErrWorkspaceAlreadyExists
	}
	mockRepository.workspacesByKey[workspace.WorkspaceKey] = workspace
	mockRepository.workspacesByID[workspace.WorkspaceID] = workspace
	mockRepository.workspacesByUser[workspace.CreatedByID] = append(mockRepository.workspacesByUser[workspace.CreatedByID], workspace)
	return nil
}

func (mockRepository *mockWorkspaceRepository) FindWorkspaceByWorkspaceID(requestContext context.Context, workspaceID string) (domainworkspace.Workspace, error) {
	workspace, hasWorkspace := mockRepository.workspacesByID[workspaceID]
	if !hasWorkspace {
		return domainworkspace.Workspace{}, domainworkspace.ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (mockRepository *mockWorkspaceRepository) FindWorkspaceByWorkspaceKey(requestContext context.Context, workspaceKey string) (domainworkspace.Workspace, error) {
	workspace, hasWorkspace := mockRepository.workspacesByKey[workspaceKey]
	if !hasWorkspace {
		return domainworkspace.Workspace{}, domainworkspace.ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (mockRepository *mockWorkspaceRepository) ListWorkspacesByUserID(requestContext context.Context, userID string, pageNumber int, pageSize int) ([]domainworkspace.Workspace, int64, error) {
	workspaceList := mockRepository.workspacesByUser[userID]
	return workspaceList, int64(len(workspaceList)), nil
}

type mockMembershipRepository struct {
	memberships map[string]domainworkspace.WorkspaceMembership
}

func (mockRepository *mockMembershipRepository) CreateMembership(requestContext context.Context, membership domainworkspace.WorkspaceMembership) error {
	membershipKey := membership.WorkspaceID + ":" + membership.UserID
	if _, hasMembership := mockRepository.memberships[membershipKey]; hasMembership {
		return domainworkspace.ErrMembershipAlreadyExists
	}
	mockRepository.memberships[membershipKey] = membership
	return nil
}

func (mockRepository *mockMembershipRepository) FindMembership(requestContext context.Context, workspaceID string, userID string) (domainworkspace.WorkspaceMembership, error) {
	membershipKey := workspaceID + ":" + userID
	membership, hasMembership := mockRepository.memberships[membershipKey]
	if !hasMembership {
		return domainworkspace.WorkspaceMembership{}, domainworkspace.ErrMembershipNotFound
	}
	return membership, nil
}

func (mockRepository *mockMembershipRepository) UpdateMembershipRole(requestContext context.Context, workspaceID string, userID string, membershipRole domainworkspace.MembershipRole) error {
	membershipKey := workspaceID + ":" + userID
	membership, hasMembership := mockRepository.memberships[membershipKey]
	if !hasMembership {
		return domainworkspace.ErrMembershipNotFound
	}
	membership.MembershipRole = membershipRole
	membership.LastUpdatedAt = time.Now()
	mockRepository.memberships[membershipKey] = membership
	return nil
}

type mockProjectRepository struct {
	projectsByWorkspace map[string][]domainworkspace.Project
}

func (mockRepository *mockProjectRepository) CreateProject(requestContext context.Context, project domainworkspace.Project) error {
	projectList := mockRepository.projectsByWorkspace[project.WorkspaceID]
	for _, existingProject := range projectList {
		if existingProject.ProjectKey == project.ProjectKey {
			return domainworkspace.ErrProjectAlreadyExists
		}
	}
	mockRepository.projectsByWorkspace[project.WorkspaceID] = append(projectList, project)
	return nil
}

func (mockRepository *mockProjectRepository) FindProjectByProjectID(requestContext context.Context, projectID string) (domainworkspace.Project, error) {
	for _, projectList := range mockRepository.projectsByWorkspace {
		for _, project := range projectList {
			if project.ProjectID == projectID {
				return project, nil
			}
		}
	}
	return domainworkspace.Project{}, domainworkspace.ErrProjectNotFound
}

func (mockRepository *mockProjectRepository) FindProjectByWorkspaceAndProjectKey(requestContext context.Context, workspaceID string, projectKey string) (domainworkspace.Project, error) {
	for _, project := range mockRepository.projectsByWorkspace[workspaceID] {
		if project.ProjectKey == projectKey {
			return project, nil
		}
	}
	return domainworkspace.Project{}, domainworkspace.ErrProjectNotFound
}

func (mockRepository *mockProjectRepository) ListProjectsByWorkspaceID(requestContext context.Context, workspaceID string, pageNumber int, pageSize int) ([]domainworkspace.Project, int64, error) {
	projectList := mockRepository.projectsByWorkspace[workspaceID]
	return projectList, int64(len(projectList)), nil
}

func TestWorkspaceLifecycle(testingSuite *testing.T) {
	nowValue := time.Unix(1700000100, 0)
	workspaceRepository := &mockWorkspaceRepository{workspacesByKey: map[string]domainworkspace.Workspace{}, workspacesByID: map[string]domainworkspace.Workspace{}, workspacesByUser: map[string][]domainworkspace.Workspace{}}
	membershipRepository := &mockMembershipRepository{memberships: map[string]domainworkspace.WorkspaceMembership{}}
	projectRepository := &mockProjectRepository{projectsByWorkspace: map[string][]domainworkspace.Project{}}
	service := NewService(workspaceRepository, membershipRepository, projectRepository, mockTransactionManager{}, WithNowFunction(func() time.Time { return nowValue }), WithIdentifierFunction(func() (string, error) { return "workspace-identifier", nil }))

	workspaceResult, workspaceError := service.CreateWorkspace(context.Background(), CreateWorkspaceInput{WorkspaceKey: "platform", DisplayName: "Platform Workspace", OwnerUserID: "owner-user"})
	if workspaceError != nil {
		testingSuite.Fatalf("workspace creation failure: %v", workspaceError)
	}
	if workspaceResult.WorkspaceKey != "platform" {
		testingSuite.Fatalf("unexpected workspace key")
	}

	membershipAddError := service.AddMembership(context.Background(), AddMembershipInput{WorkspaceID: workspaceResult.WorkspaceID, ActorUserID: "owner-user", TargetUserID: "admin-user", TargetRole: domainworkspace.MembershipRoleAdmin, InvitedByUserID: "owner-user"})
	if membershipAddError != nil {
		testingSuite.Fatalf("membership add failure: %v", membershipAddError)
	}

	projectResult, projectError := service.CreateProject(context.Background(), CreateProjectInput{WorkspaceID: workspaceResult.WorkspaceID, ActorUserID: "admin-user", ProjectKey: "api", DisplayName: "API", Description: "Backend API"})
	if projectError != nil {
		testingSuite.Fatalf("project creation failure: %v", projectError)
	}
	if projectResult.ProjectKey != "API" {
		testingSuite.Fatalf("project key should be uppercased")
	}

	workspacePage, workspaceListError := service.ListUserWorkspaces(context.Background(), "owner-user", 1, 20)
	if workspaceListError != nil || workspacePage.TotalCount != 1 {
		testingSuite.Fatalf("workspace list failure: %v", workspaceListError)
	}

	projectPage, projectListError := service.ListWorkspaceProjects(context.Background(), workspaceResult.WorkspaceID, 1, 20)
	if projectListError != nil || projectPage.TotalCount != 1 {
		testingSuite.Fatalf("project list failure: %v", projectListError)
	}
}

func TestWorkspaceAuthorization(testingSuite *testing.T) {
	workspaceRepository := &mockWorkspaceRepository{workspacesByKey: map[string]domainworkspace.Workspace{}, workspacesByID: map[string]domainworkspace.Workspace{}, workspacesByUser: map[string][]domainworkspace.Workspace{}}
	membershipRepository := &mockMembershipRepository{memberships: map[string]domainworkspace.WorkspaceMembership{"workspace-id:member-user": {WorkspaceID: "workspace-id", UserID: "member-user", MembershipRole: domainworkspace.MembershipRoleMember}}}
	projectRepository := &mockProjectRepository{projectsByWorkspace: map[string][]domainworkspace.Project{}}
	service := NewService(workspaceRepository, membershipRepository, projectRepository, mockTransactionManager{})

	membershipError := service.AddMembership(context.Background(), AddMembershipInput{WorkspaceID: "workspace-id", ActorUserID: "member-user", TargetUserID: "other-user", TargetRole: domainworkspace.MembershipRoleMember, InvitedByUserID: "member-user"})
	if !errors.Is(membershipError, domainworkspace.ErrInsufficientWorkspaceRole) {
		testingSuite.Fatalf("expected insufficient role error, got: %v", membershipError)
	}

	_, projectError := service.CreateProject(context.Background(), CreateProjectInput{WorkspaceID: "workspace-id", ActorUserID: "member-user", ProjectKey: "core", DisplayName: "Core", Description: "Core project"})
	if !errors.Is(projectError, domainworkspace.ErrInsufficientWorkspaceRole) {
		testingSuite.Fatalf("expected insufficient role for project creation, got: %v", projectError)
	}
}

func TestWorkspaceConflictAndRoleChange(testingSuite *testing.T) {
	workspaceRepository := &mockWorkspaceRepository{workspacesByKey: map[string]domainworkspace.Workspace{"platform": {WorkspaceID: "existing-workspace", WorkspaceKey: "platform", CreatedByID: "owner-user"}}, workspacesByID: map[string]domainworkspace.Workspace{"existing-workspace": {WorkspaceID: "existing-workspace", WorkspaceKey: "platform", CreatedByID: "owner-user"}}, workspacesByUser: map[string][]domainworkspace.Workspace{}}
	membershipRepository := &mockMembershipRepository{memberships: map[string]domainworkspace.WorkspaceMembership{"existing-workspace:owner-user": {WorkspaceID: "existing-workspace", UserID: "owner-user", MembershipRole: domainworkspace.MembershipRoleOwner}, "existing-workspace:member-user": {WorkspaceID: "existing-workspace", UserID: "member-user", MembershipRole: domainworkspace.MembershipRoleMember}}}
	projectRepository := &mockProjectRepository{projectsByWorkspace: map[string][]domainworkspace.Project{"existing-workspace": {{ProjectID: "project-id", WorkspaceID: "existing-workspace", ProjectKey: "CORE"}}}}
	service := NewService(workspaceRepository, membershipRepository, projectRepository, mockTransactionManager{}, WithIdentifierFunction(func() (string, error) { return "existing-workspace", nil }))

	_, workspaceError := service.CreateWorkspace(context.Background(), CreateWorkspaceInput{WorkspaceKey: "platform", DisplayName: "Duplicate", OwnerUserID: "owner-user"})
	if !errors.Is(workspaceError, domainworkspace.ErrWorkspaceAlreadyExists) {
		testingSuite.Fatalf("expected workspace conflict, got: %v", workspaceError)
	}

	roleChangeError := service.ChangeMembershipRole(context.Background(), "existing-workspace", "owner-user", "member-user", domainworkspace.MembershipRoleAdmin)
	if roleChangeError != nil {
		testingSuite.Fatalf("role change failure: %v", roleChangeError)
	}

	_, projectError := service.CreateProject(context.Background(), CreateProjectInput{WorkspaceID: "existing-workspace", ActorUserID: "owner-user", ProjectKey: "core", DisplayName: "Core", Description: "Duplicate"})
	if !errors.Is(projectError, domainworkspace.ErrProjectAlreadyExists) {
		testingSuite.Fatalf("expected project conflict, got: %v", projectError)
	}
}

func TestWorkspaceIdentifierFailure(testingSuite *testing.T) {
	workspaceRepository := &mockWorkspaceRepository{workspacesByKey: map[string]domainworkspace.Workspace{}, workspacesByID: map[string]domainworkspace.Workspace{}, workspacesByUser: map[string][]domainworkspace.Workspace{}}
	membershipRepository := &mockMembershipRepository{memberships: map[string]domainworkspace.WorkspaceMembership{}}
	projectRepository := &mockProjectRepository{projectsByWorkspace: map[string][]domainworkspace.Project{}}
	service := NewService(workspaceRepository, membershipRepository, projectRepository, mockTransactionManager{}, WithIdentifierFunction(func() (string, error) { return "", errors.New("identifier failure") }))

	_, workspaceError := service.CreateWorkspace(context.Background(), CreateWorkspaceInput{WorkspaceKey: "platform", DisplayName: "Platform", OwnerUserID: "owner-user"})
	if workspaceError == nil {
		testingSuite.Fatalf("expected identifier failure error")
	}
}
