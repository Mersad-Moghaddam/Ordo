package workspace

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	domainworkspace "github.com/ordo/backend/internal/domain/workspace"
	usecaseworkspace "github.com/ordo/backend/internal/usecase/workspace"
)

type Handler struct {
	workspaceService *usecaseworkspace.Service
}

func NewHandler(workspaceService *usecaseworkspace.Service) *Handler {
	return &Handler{workspaceService: workspaceService}
}

func (handler *Handler) RegisterRoutes(fiberRouter fiber.Router) {
	fiberRouter.Post("/workspaces", handler.createWorkspace)
	fiberRouter.Post("/workspaces/:workspaceId/memberships", handler.addMembership)
	fiberRouter.Post("/workspaces/:workspaceId/projects", handler.createProject)
	fiberRouter.Get("/users/:userId/workspaces", handler.listWorkspaces)
	fiberRouter.Get("/workspaces/:workspaceId/projects", handler.listProjects)
}

func (handler *Handler) createWorkspace(fiberContext *fiber.Ctx) error {
	var createWorkspaceRequest CreateWorkspaceRequest
	if parseError := fiberContext.BodyParser(&createWorkspaceRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	workspace, createError := handler.workspaceService.CreateWorkspace(fiberContext.UserContext(), usecaseworkspace.CreateWorkspaceInput{WorkspaceKey: createWorkspaceRequest.WorkspaceKey, DisplayName: createWorkspaceRequest.DisplayName, OwnerUserID: createWorkspaceRequest.OwnerUserID})
	if createError != nil {
		return respondWithWorkspaceError(fiberContext, createError)
	}
	return fiberContext.Status(fiber.StatusCreated).JSON(workspace)
}

func (handler *Handler) addMembership(fiberContext *fiber.Ctx) error {
	workspaceID := fiberContext.Params("workspaceId")
	var addMembershipRequest AddMembershipRequest
	if parseError := fiberContext.BodyParser(&addMembershipRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	addError := handler.workspaceService.AddMembership(fiberContext.UserContext(), usecaseworkspace.AddMembershipInput{WorkspaceID: workspaceID, ActorUserID: addMembershipRequest.ActorUserID, TargetUserID: addMembershipRequest.TargetUserID, TargetRole: domainworkspace.MembershipRole(addMembershipRequest.TargetRole), InvitedByUserID: addMembershipRequest.InvitedByUserID})
	if addError != nil {
		return respondWithWorkspaceError(fiberContext, addError)
	}
	return fiberContext.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) createProject(fiberContext *fiber.Ctx) error {
	workspaceID := fiberContext.Params("workspaceId")
	var createProjectRequest CreateProjectRequest
	if parseError := fiberContext.BodyParser(&createProjectRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	project, createError := handler.workspaceService.CreateProject(fiberContext.UserContext(), usecaseworkspace.CreateProjectInput{WorkspaceID: workspaceID, ActorUserID: createProjectRequest.ActorUserID, ProjectKey: createProjectRequest.ProjectKey, DisplayName: createProjectRequest.DisplayName, Description: createProjectRequest.Description})
	if createError != nil {
		return respondWithWorkspaceError(fiberContext, createError)
	}
	return fiberContext.Status(fiber.StatusCreated).JSON(project)
}

func (handler *Handler) listWorkspaces(fiberContext *fiber.Ctx) error {
	userID := fiberContext.Params("userId")
	pageNumber, pageSize := parsePagination(fiberContext)
	workspacePage, listError := handler.workspaceService.ListUserWorkspaces(fiberContext.UserContext(), userID, pageNumber, pageSize)
	if listError != nil {
		return respondWithWorkspaceError(fiberContext, listError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(workspacePage)
}

func (handler *Handler) listProjects(fiberContext *fiber.Ctx) error {
	workspaceID := fiberContext.Params("workspaceId")
	pageNumber, pageSize := parsePagination(fiberContext)
	projectPage, listError := handler.workspaceService.ListWorkspaceProjects(fiberContext.UserContext(), workspaceID, pageNumber, pageSize)
	if listError != nil {
		return respondWithWorkspaceError(fiberContext, listError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(projectPage)
}

func parsePagination(fiberContext *fiber.Ctx) (int, int) {
	pageNumber := 1
	pageSize := 20
	if queryPage := fiberContext.Query("page"); queryPage != "" {
		if parsedPageNumber, parseError := strconv.Atoi(queryPage); parseError == nil && parsedPageNumber > 0 {
			pageNumber = parsedPageNumber
		}
	}
	if queryPageSize := fiberContext.Query("pageSize"); queryPageSize != "" {
		if parsedPageSize, parseError := strconv.Atoi(queryPageSize); parseError == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}
	return pageNumber, pageSize
}

func respondWithWorkspaceError(fiberContext *fiber.Ctx, responseError error) error {
	if errors.Is(responseError, domainworkspace.ErrInsufficientWorkspaceRole) {
		return fiberContext.Status(fiber.StatusForbidden).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domainworkspace.ErrWorkspaceAlreadyExists) || errors.Is(responseError, domainworkspace.ErrProjectAlreadyExists) || errors.Is(responseError, domainworkspace.ErrMembershipAlreadyExists) {
		return fiberContext.Status(fiber.StatusConflict).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domainworkspace.ErrWorkspaceNotFound) || errors.Is(responseError, domainworkspace.ErrProjectNotFound) || errors.Is(responseError, domainworkspace.ErrMembershipNotFound) {
		return fiberContext.Status(fiber.StatusNotFound).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	return fiberContext.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "internal error"})
}
