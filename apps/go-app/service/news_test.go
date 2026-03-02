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

func TestNewsService_CreateNews(t *testing.T) {
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockTopicRepo := new(mocks.MockTopicRepository)
	svc := service.NewNewsService(mockNewsRepo, mockTopicRepo)
	ctx := context.Background()

	newsID := uuid.New()
	topicIDs := []uuid.UUID{uuid.New()}
	req := &domain.CreateNewsRequest{
		Title:    "Title",
		Content:  "Content",
		TopicIDs: topicIDs,
	}
	news := &domain.News{ID: newsID, Title: "Title"}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("CreateNews", ctx, req).Return(news, nil).Once()
		mockNewsRepo.On("AssignTopics", ctx, newsID, topicIDs).Return(nil).Once()
		mockNewsRepo.On("GetNews", ctx, newsID).Return(news, nil).Once()

		res, err := svc.CreateNews(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, news, res)
		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		mockNewsRepo.On("CreateNews", ctx, req).Return(nil, errors.New("db error")).Once()
		res, err := svc.CreateNews(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestNewsService_GetNewsList(t *testing.T) {
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockTopicRepo := new(mocks.MockTopicRepository)
	svc := service.NewNewsService(mockNewsRepo, mockTopicRepo)
	ctx := context.Background()

	filter := &domain.NewsFilter{}
	expected := []domain.News{{ID: uuid.New(), Title: "News 1"}}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("GetNewsList", ctx, filter).Return(expected, nil).Once()
		res, err := svc.GetNewsList(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestNewsService_GetNews(t *testing.T) {
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockTopicRepo := new(mocks.MockTopicRepository)
	svc := service.NewNewsService(mockNewsRepo, mockTopicRepo)
	ctx := context.Background()

	id := uuid.New()
	expected := &domain.News{ID: id, Title: "News 1"}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("GetNews", ctx, id).Return(expected, nil).Once()
		res, err := svc.GetNews(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestNewsService_UpdateNews(t *testing.T) {
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockTopicRepo := new(mocks.MockTopicRepository)
	svc := service.NewNewsService(mockNewsRepo, mockTopicRepo)
	ctx := context.Background()

	id := uuid.New()
	topicIDs := []uuid.UUID{uuid.New()}
	req := &domain.UpdateNewsRequest{
		Title:    "Updated",
		Content:  "Updated content",
		Status:   domain.NewsStatusPublished,
		TopicIDs: topicIDs,
	}
	news := &domain.News{ID: id, Title: "Updated"}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("UpdateNews", ctx, id, req).Return(news, nil).Once()
		mockNewsRepo.On("AssignTopics", ctx, id, topicIDs).Return(nil).Once()
		mockNewsRepo.On("GetNews", ctx, id).Return(news, nil).Once()

		res, err := svc.UpdateNews(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, news, res)
	})
}

func TestNewsService_DeleteNews(t *testing.T) {
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockTopicRepo := new(mocks.MockTopicRepository)
	svc := service.NewNewsService(mockNewsRepo, mockTopicRepo)
	ctx := context.Background()

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("DeleteNews", ctx, id).Return(nil).Once()
		err := svc.DeleteNews(ctx, id)
		assert.NoError(t, err)
	})
}
