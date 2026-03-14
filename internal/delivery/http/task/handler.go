package task

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	domaintask "github.com/ordo/backend/internal/domain/task"
	usecasetask "github.com/ordo/backend/internal/usecase/task"
)

type Handler struct {
	taskService *usecasetask.Service
}

func NewHandler(taskService *usecasetask.Service) *Handler {
	return &Handler{taskService: taskService}
}

func (handler *Handler) RegisterRoutes(fiberRouter fiber.Router) {
	fiberRouter.Post("/tasks", handler.createTask)
	fiberRouter.Patch("/tasks/:taskId/status", handler.updateTaskStatus)
	fiberRouter.Get("/projects/:projectId/tasks", handler.listProjectTasks)
}

func (handler *Handler) createTask(fiberContext *fiber.Ctx) error {
	var createTaskRequest CreateTaskRequest
	if parseError := fiberContext.BodyParser(&createTaskRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	task, createError := handler.taskService.CreateTask(fiberContext.UserContext(), usecasetask.CreateTaskInput{WorkspaceID: createTaskRequest.WorkspaceID, ProjectID: createTaskRequest.ProjectID, Title: createTaskRequest.Title, Description: createTaskRequest.Description, Priority: domaintask.TaskPriority(createTaskRequest.Priority), AssigneeUserID: createTaskRequest.AssigneeUserID, CreatedByUserID: createTaskRequest.CreatedByUserID})
	if createError != nil {
		return respondWithTaskError(fiberContext, createError)
	}
	return fiberContext.Status(fiber.StatusCreated).JSON(task)
}

func (handler *Handler) updateTaskStatus(fiberContext *fiber.Ctx) error {
	taskID := fiberContext.Params("taskId")
	var updateRequest UpdateTaskStatusRequest
	if parseError := fiberContext.BodyParser(&updateRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	updateError := handler.taskService.UpdateTaskStatus(fiberContext.UserContext(), taskID, domaintask.TaskStatus(updateRequest.TaskStatus))
	if updateError != nil {
		return respondWithTaskError(fiberContext, updateError)
	}
	return fiberContext.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) listProjectTasks(fiberContext *fiber.Ctx) error {
	projectID := fiberContext.Params("projectId")
	pageNumber, pageSize := parsePagination(fiberContext)
	taskPage, listError := handler.taskService.ListProjectTasks(fiberContext.UserContext(), projectID, pageNumber, pageSize)
	if listError != nil {
		return respondWithTaskError(fiberContext, listError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(taskPage)
}

func parsePagination(fiberContext *fiber.Ctx) (int, int) {
	pageNumber := 1
	pageSize := 20
	queryPage := fiberContext.Query("page")
	if queryPage != "" {
		parsedPageNumber, parseError := strconv.Atoi(queryPage)
		if parseError == nil && parsedPageNumber > 0 {
			pageNumber = parsedPageNumber
		}
	}
	queryPageSize := fiberContext.Query("pageSize")
	if queryPageSize != "" {
		parsedPageSize, parseError := strconv.Atoi(queryPageSize)
		if parseError == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}
	return pageNumber, pageSize
}

func respondWithTaskError(fiberContext *fiber.Ctx, responseError error) error {
	if errors.Is(responseError, domaintask.ErrTaskNotFound) {
		return fiberContext.Status(fiber.StatusNotFound).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domaintask.ErrTaskAlreadyExists) || errors.Is(responseError, domaintask.ErrInvalidTaskStatusTransition) {
		return fiberContext.Status(fiber.StatusConflict).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	return fiberContext.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "internal error"})
}
