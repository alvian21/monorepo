package rest

import (
	"context"
	"go-app/domain"
	"go-app/internal/logging"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NewsService interface {
	CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error)
	GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error)
	GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error)
	UpdateNews(ctx context.Context, id uuid.UUID, req *domain.UpdateNewsRequest) (*domain.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
}

type NewsHandler struct {
	Service NewsService
}

func NewNewsHandler(e *echo.Group, svc NewsService) {
	handler := &NewsHandler{
		Service: svc,
	}

	e.GET("", handler.GetNewsList)
	e.GET("/:id", handler.GetNews)
	e.POST("", handler.CreateNews)
	e.PUT("/:id", handler.UpdateNews)
	e.DELETE("/:id", handler.DeleteNews)
}

func (h *NewsHandler) GetNewsList(c echo.Context) error {
	ctx := c.Request().Context()
	filter := new(domain.NewsFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind news filter")
	}

	newsList, err := h.Service.GetNewsList(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.News]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved news",
		Data:    newsList,
	})
}

func (h *NewsHandler) GetNews(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	news, err := h.Service.GetNews(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusNotFound,
			Message: "News not found",
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.News]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved news",
		Data:    *news,
	})
}

func (h *NewsHandler) CreateNews(c echo.Context) error {
	var req domain.CreateNewsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: FormatValidationError(err),
		})
	}

	news, err := h.Service.CreateNews(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.News]{
		Code:    http.StatusCreated,
		Message: "Successfully created news",
		Data:    *news,
	})
}

func (h *NewsHandler) UpdateNews(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var req domain.UpdateNewsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: FormatValidationError(err),
		})
	}

	news, err := h.Service.UpdateNews(c.Request().Context(), id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.News]{
		Code:    http.StatusOK,
		Message: "Successfully updated news",
		Data:    *news,
	})
}

func (h *NewsHandler) DeleteNews(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	err = h.Service.DeleteNews(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusOK,
		Message: "Successfully deleted news",
	})
}
