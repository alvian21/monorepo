package rest

import (
	"context"
	"go-app/domain"
	"go-app/internal/logging"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TopicService interface {
	CreateTopic(ctx context.Context, req *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, req *domain.UpdateTopicRequest) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
}

type TopicHandler struct {
	Service TopicService
}

func NewTopicHandler(e *echo.Group, svc TopicService) {
	handler := &TopicHandler{
		Service: svc,
	}

	e.GET("", handler.GetTopicList)
	e.GET("/:id", handler.GetTopic)
	e.POST("", handler.CreateTopic)
	e.PUT("/:id", handler.UpdateTopic)
	e.DELETE("/:id", handler.DeleteTopic)
}

func (h *TopicHandler) GetTopicList(c echo.Context) error {
	ctx := c.Request().Context()
	filter := new(domain.TopicFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind topic filter")
	}

	topics, err := h.Service.GetTopicList(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.Topic]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved topics",
		Data:    topics,
	})
}

func (h *TopicHandler) GetTopic(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	topic, err := h.Service.GetTopic(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusNotFound,
			Message: "Topic not found",
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved topic",
		Data:    *topic,
	})
}

func (h *TopicHandler) CreateTopic(c echo.Context) error {
	var req domain.CreateTopicRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	topic, err := h.Service.CreateTopic(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.Topic]{
		Code:    http.StatusCreated,
		Message: "Successfully created topic",
		Data:    *topic,
	})
}

func (h *TopicHandler) UpdateTopic(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var req domain.UpdateTopicRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	topic, err := h.Service.UpdateTopic(c.Request().Context(), id, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Topic]{
		Code:    http.StatusOK,
		Message: "Successfully updated topic",
		Data:    *topic,
	})
}

func (h *TopicHandler) DeleteTopic(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	err = h.Service.DeleteTopic(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusOK,
		Message: "Successfully deleted topic",
	})
}
