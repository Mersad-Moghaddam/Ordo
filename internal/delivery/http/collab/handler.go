package collab

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	domaincollab "github.com/ordo/backend/internal/domain/collab"
	usecasecollab "github.com/ordo/backend/internal/usecase/collab"
)

type Handler struct {
	collabService *usecasecollab.Service
}

func NewHandler(collabService *usecasecollab.Service) *Handler {
	return &Handler{collabService: collabService}
}

func (handler *Handler) RegisterRoutes(fiberRouter fiber.Router) {
	fiberRouter.Post("/comments", handler.createComment)
	fiberRouter.Patch("/comments/:commentId", handler.updateComment)
	fiberRouter.Delete("/comments/:commentId", handler.deleteComment)
	fiberRouter.Get("/tasks/:taskId/comments", handler.listComments)
	fiberRouter.Get("/tasks/:taskId/activities", handler.listActivities)
}

func (handler *Handler) createComment(fiberContext *fiber.Ctx) error {
	var createRequest CreateCommentRequest
	if parseError := fiberContext.BodyParser(&createRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	comment, createError := handler.collabService.CreateComment(fiberContext.UserContext(), usecasecollab.CreateCommentInput{WorkspaceID: createRequest.WorkspaceID, ProjectID: createRequest.ProjectID, TaskID: createRequest.TaskID, AuthorUserID: createRequest.AuthorUserID, Body: createRequest.Body})
	if createError != nil {
		return respondWithCollabError(fiberContext, createError)
	}
	return fiberContext.Status(fiber.StatusCreated).JSON(comment)
}

func (handler *Handler) updateComment(fiberContext *fiber.Ctx) error {
	commentID := fiberContext.Params("commentId")
	var updateRequest UpdateCommentRequest
	if parseError := fiberContext.BodyParser(&updateRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	updateError := handler.collabService.UpdateComment(fiberContext.UserContext(), commentID, updateRequest.ActorUserID, updateRequest.Body)
	if updateError != nil {
		return respondWithCollabError(fiberContext, updateError)
	}
	return fiberContext.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deleteComment(fiberContext *fiber.Ctx) error {
	commentID := fiberContext.Params("commentId")
	var deleteRequest DeleteCommentRequest
	if parseError := fiberContext.BodyParser(&deleteRequest); parseError != nil {
		return fiberContext.Status(fiber.StatusBadRequest).JSON(ErrorResponse{ErrorMessage: "invalid payload"})
	}
	deleteError := handler.collabService.DeleteComment(fiberContext.UserContext(), commentID, deleteRequest.ActorUserID)
	if deleteError != nil {
		return respondWithCollabError(fiberContext, deleteError)
	}
	return fiberContext.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) listComments(fiberContext *fiber.Ctx) error {
	taskID := fiberContext.Params("taskId")
	pageNumber, pageSize := parsePagination(fiberContext)
	commentPage, listError := handler.collabService.ListTaskComments(fiberContext.UserContext(), taskID, pageNumber, pageSize)
	if listError != nil {
		return respondWithCollabError(fiberContext, listError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(commentPage)
}

func (handler *Handler) listActivities(fiberContext *fiber.Ctx) error {
	taskID := fiberContext.Params("taskId")
	pageNumber, pageSize := parsePagination(fiberContext)
	activityPage, listError := handler.collabService.ListTaskActivities(fiberContext.UserContext(), taskID, pageNumber, pageSize)
	if listError != nil {
		return respondWithCollabError(fiberContext, listError)
	}
	return fiberContext.Status(fiber.StatusOK).JSON(activityPage)
}

func parsePagination(fiberContext *fiber.Ctx) (int, int) {
	pageNumber := 1
	pageSize := 20
	if pageValue := fiberContext.Query("page"); pageValue != "" {
		if parsedPage, parseError := strconv.Atoi(pageValue); parseError == nil && parsedPage > 0 {
			pageNumber = parsedPage
		}
	}
	if pageSizeValue := fiberContext.Query("pageSize"); pageSizeValue != "" {
		if parsedPageSize, parseError := strconv.Atoi(pageSizeValue); parseError == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}
	return pageNumber, pageSize
}

func respondWithCollabError(fiberContext *fiber.Ctx, responseError error) error {
	if errors.Is(responseError, domaincollab.ErrCommentForbidden) {
		return fiberContext.Status(fiber.StatusForbidden).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domaincollab.ErrCommentNotFound) {
		return fiberContext.Status(fiber.StatusNotFound).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	if errors.Is(responseError, domaincollab.ErrCommentAlreadyDeleted) {
		return fiberContext.Status(fiber.StatusConflict).JSON(ErrorResponse{ErrorMessage: responseError.Error()})
	}
	return fiberContext.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{ErrorMessage: "internal error"})
}
