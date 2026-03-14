package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
	"github.com/ordo/backend/internal/repository"
	repositoryworkspace "github.com/ordo/backend/internal/repository/workspace"
	"github.com/ordo/backend/internal/usecase"
)

type Service struct {
	workspaceRepository  repositoryworkspace.WorkspaceRepository
	membershipRepository repositoryworkspace.MembershipRepository
	projectRepository    repositoryworkspace.ProjectRepository
	transactionManager   repository.TransactionManager
	nowFunction          func() time.Time
	identifierFunction   func() (string, error)
}

type Option func(service *Service)

type CreateWorkspaceInput struct {
	WorkspaceKey string
	DisplayName  string
	OwnerUserID  string
}

type AddMembershipInput struct {
	WorkspaceID     string
	ActorUserID     string
	TargetUserID    string
	TargetRole      domainworkspace.MembershipRole
	InvitedByUserID string
}

type CreateProjectInput struct {
	WorkspaceID string
	ActorUserID string
	ProjectKey  string
	DisplayName string
	Description string
}

func NewService(
	workspaceRepository repositoryworkspace.WorkspaceRepository,
	membershipRepository repositoryworkspace.MembershipRepository,
	projectRepository repositoryworkspace.ProjectRepository,
	transactionManager repository.TransactionManager,
	options ...Option,
) *Service {
	service := &Service{
		workspaceRepository:  workspaceRepository,
		membershipRepository: membershipRepository,
		projectRepository:    projectRepository,
		transactionManager:   transactionManager,
		nowFunction:          time.Now,
		identifierFunction:   defaultIdentifierFunction,
	}
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

func (service *Service) CreateWorkspace(requestContext context.Context, createInput CreateWorkspaceInput) (domainworkspace.Workspace, error) {
	workspaceID, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domainworkspace.Workspace{}, fmt.Errorf("workspace identifier generation failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	workspace := domainworkspace.Workspace{WorkspaceID: workspaceID, WorkspaceKey: strings.ToLower(createInput.WorkspaceKey), DisplayName: createInput.DisplayName, CreatedByID: createInput.OwnerUserID, CreatedAt: creationTime, UpdatedAt: creationTime}
	transactionError := service.transactionManager.WithTransaction(requestContext, func(transactionContext context.Context) error {
		return service.createWorkspaceWithOwner(transactionContext, workspace, createInput.OwnerUserID)
	})
	if transactionError != nil {
		return domainworkspace.Workspace{}, transactionError
	}
	return workspace, nil
}

func (service *Service) AddMembership(requestContext context.Context, membershipInput AddMembershipInput) error {
	if membershipInput.TargetRole == "" {
		membershipInput.TargetRole = domainworkspace.MembershipRoleMember
	}
	actorMembership, membershipError := service.membershipRepository.FindMembership(requestContext, membershipInput.WorkspaceID, membershipInput.ActorUserID)
	if membershipError != nil {
		return membershipError
	}
	if !canManageMembership(actorMembership.MembershipRole) {
		return domainworkspace.ErrInsufficientWorkspaceRole
	}
	newMembership := domainworkspace.WorkspaceMembership{WorkspaceID: membershipInput.WorkspaceID, UserID: membershipInput.TargetUserID, MembershipRole: membershipInput.TargetRole, InvitedByUserID: membershipInput.InvitedByUserID, JoinedAt: service.nowFunction(), LastUpdatedAt: service.nowFunction()}
	return service.membershipRepository.CreateMembership(requestContext, newMembership)
}

func (service *Service) ChangeMembershipRole(requestContext context.Context, workspaceID string, actorUserID string, targetUserID string, targetRole domainworkspace.MembershipRole) error {
	actorMembership, actorError := service.membershipRepository.FindMembership(requestContext, workspaceID, actorUserID)
	if actorError != nil {
		return actorError
	}
	if actorMembership.MembershipRole != domainworkspace.MembershipRoleOwner {
		return domainworkspace.ErrInsufficientWorkspaceRole
	}
	return service.membershipRepository.UpdateMembershipRole(requestContext, workspaceID, targetUserID, targetRole)
}

func (service *Service) CreateProject(requestContext context.Context, createInput CreateProjectInput) (domainworkspace.Project, error) {
	actorMembership, membershipError := service.membershipRepository.FindMembership(requestContext, createInput.WorkspaceID, createInput.ActorUserID)
	if membershipError != nil {
		return domainworkspace.Project{}, membershipError
	}
	if !canCreateProject(actorMembership.MembershipRole) {
		return domainworkspace.Project{}, domainworkspace.ErrInsufficientWorkspaceRole
	}
	projectID, identifierError := service.identifierFunction()
	if identifierError != nil {
		return domainworkspace.Project{}, fmt.Errorf("project identifier generation failure: %w", identifierError)
	}
	creationTime := service.nowFunction()
	project := domainworkspace.Project{ProjectID: projectID, WorkspaceID: createInput.WorkspaceID, ProjectKey: strings.ToUpper(createInput.ProjectKey), DisplayName: createInput.DisplayName, Description: createInput.Description, CreatedByID: createInput.ActorUserID, CreatedAt: creationTime, LastUpdatedAt: creationTime}
	if existingProject, existingError := service.projectRepository.FindProjectByWorkspaceAndProjectKey(requestContext, createInput.WorkspaceID, project.ProjectKey); existingError == nil && existingProject.ProjectID != "" {
		return domainworkspace.Project{}, domainworkspace.ErrProjectAlreadyExists
	}
	if createError := service.projectRepository.CreateProject(requestContext, project); createError != nil {
		if errors.Is(createError, domainworkspace.ErrProjectAlreadyExists) {
			return domainworkspace.Project{}, createError
		}
		return domainworkspace.Project{}, fmt.Errorf("project creation failure: %w", createError)
	}
	return project, nil
}

func (service *Service) ListUserWorkspaces(requestContext context.Context, userID string, pageNumber int, pageSize int) (usecase.PageResult[domainworkspace.Workspace], error) {
	workspaceList, totalCount, listError := service.workspaceRepository.ListWorkspacesByUserID(requestContext, userID, pageNumber, pageSize)
	if listError != nil {
		return usecase.PageResult[domainworkspace.Workspace]{}, listError
	}
	return usecase.PageResult[domainworkspace.Workspace]{Items: workspaceList, TotalCount: totalCount, PageNumber: pageNumber, PageSize: pageSize}, nil
}

func (service *Service) ListWorkspaceProjects(requestContext context.Context, workspaceID string, pageNumber int, pageSize int) (usecase.PageResult[domainworkspace.Project], error) {
	projectList, totalCount, listError := service.projectRepository.ListProjectsByWorkspaceID(requestContext, workspaceID, pageNumber, pageSize)
	if listError != nil {
		return usecase.PageResult[domainworkspace.Project]{}, listError
	}
	return usecase.PageResult[domainworkspace.Project]{Items: projectList, TotalCount: totalCount, PageNumber: pageNumber, PageSize: pageSize}, nil
}
