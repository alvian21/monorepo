package service_test

import (
	"context"
	"errors"
	"go-app/domain"
	"go-app/service"
	"go-app/service/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTopicService_CreateTopic(t *testing.T) {
	mockRepo := new(mocks.MockTopicRepository)
	svc := service.NewTopicService(mockRepo)
	ctx := context.Background()

	req := &domain.CreateTopicRequest{Name: "General"}
	expected := &domain.Topic{ID: uuid.New(), Name: "General"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreateTopic", ctx, req).Return(expected, nil).Once()
		res, err := svc.CreateTopic(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("CreateTopic", ctx, req).Return(nil, errors.New("db error")).Once()
		res, err := svc.CreateTopic(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTopicService_GetTopicList(t *testing.T) {
	mockRepo := new(mocks.MockTopicRepository)
	svc := service.NewTopicService(mockRepo)
	ctx := context.Background()

	filter := &domain.TopicFilter{}
	expected := []domain.Topic{{ID: uuid.New(), Name: "Topic 1"}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetTopicList", ctx, filter).Return(expected, nil).Once()
		res, err := svc.GetTopicList(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTopicService_GetTopic(t *testing.T) {
	mockRepo := new(mocks.MockTopicRepository)
	svc := service.NewTopicService(mockRepo)
	ctx := context.Background()

	id := uuid.New()
	expected := &domain.Topic{ID: id, Name: "Topic 1"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetTopic", ctx, id).Return(expected, nil).Once()
		res, err := svc.GetTopic(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTopicService_UpdateTopic(t *testing.T) {
	mockRepo := new(mocks.MockTopicRepository)
	svc := service.NewTopicService(mockRepo)
	ctx := context.Background()

	id := uuid.New()
	req := &domain.UpdateTopicRequest{Name: "Updated"}
	expected := &domain.Topic{ID: id, Name: "Updated"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateTopic", ctx, id, req).Return(expected, nil).Once()
		res, err := svc.UpdateTopic(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTopicService_DeleteTopic(t *testing.T) {
	mockRepo := new(mocks.MockTopicRepository)
	svc := service.NewTopicService(mockRepo)
	ctx := context.Background()

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteTopic", ctx, id).Return(nil).Once()
		err := svc.DeleteTopic(ctx, id)
		assert.NoError(t, err)
	})
}
