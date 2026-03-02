package rest_test

import (
	"context"
	"fmt"
	"go-app/domain"
	"go-app/internal/repository/postgres"
	"go-app/internal/rest"
	"go-app/service"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewsCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	topicRepo := postgres.NewTopicRepository(kit.DB)
	newsRepo := postgres.NewNewsRepository(kit.DB)
	
	topicSvc := service.NewTopicService(topicRepo)
	newsSvc := service.NewNewsService(newsRepo, topicRepo)
	
	rest.NewTopicHandler(kit.Echo.Group("/api/v1/topics"), topicSvc)
	rest.NewNewsHandler(kit.Echo.Group("/api/v1/news"), newsSvc)

	// Now start the test server
	kit.Start(t)

	// 1. Create a Topic first
	topicName := fmt.Sprintf("Topic-%s", uuid.New().String()[:8])
	topicReq := domain.CreateTopicRequest{Name: topicName}
	type TopicRes domain.ResponseSingleData[domain.Topic]
	topicRes, code := doRequest[TopicRes](t, http.MethodPost, kit.BaseURL+"/api/v1/topics", topicReq)
	require.Equal(t, http.StatusCreated, code)
	topicID := topicRes.Data.ID

	newsTitle := fmt.Sprintf("News-%s", uuid.New().String()[:8])
	newsReq := domain.CreateNewsRequest{
		Title:    newsTitle,
		Content:  "Some content here",
		TopicIDs: []uuid.UUID{topicID},
	}
	type NewsRes domain.ResponseSingleData[domain.News]
	newsRes, code := doRequest[NewsRes](t, http.MethodPost, kit.BaseURL+"/api/v1/news", newsReq)
	require.Equal(t, http.StatusCreated, code)
	news := newsRes.Data
	require.NotEmpty(t, news.ID)
	require.Equal(t, newsTitle, news.Title)

	// 3. Get News List
	type NewsListRes domain.ResponseMultipleData[domain.News]
	list, code := doRequest[NewsListRes](t, http.MethodGet, kit.BaseURL+"/api/v1/news", nil)
	require.Equal(t, http.StatusOK, code)
	require.True(t, len(list.Data) >= 1)

	// 4. Update News
	updatedNewsTitle := fmt.Sprintf("Updated-%s", uuid.New().String()[:8])
	updReq := domain.UpdateNewsRequest{
		Title:    updatedNewsTitle,
		Content:  "Updated content here",
		Status:   domain.NewsStatusPublished,
		TopicIDs: []uuid.UUID{topicID},
	}
	newsRes, code = doRequest[NewsRes](t, http.MethodPut, fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID), updReq)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, updatedNewsTitle, newsRes.Data.Title)

	// 5. Delete News
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 6. Get after delete
	_, code = doRequest[NewsRes](t, http.MethodGet, fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID), nil)
	require.Equal(t, http.StatusNotFound, code)

	// Cleanup
	_, err = kit.DB.Exec(context.Background(), "DELETE FROM news WHERE id = $1", news.ID)
	require.NoError(t, err)
	_, err = kit.DB.Exec(context.Background(), "DELETE FROM topic WHERE id = $1", topicID)
	require.NoError(t, err)
}
