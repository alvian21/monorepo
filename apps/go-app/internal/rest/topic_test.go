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

func TestTopicCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	topicRepo := postgres.NewTopicRepository(kit.DB)
	topicSvc := service.NewTopicService(topicRepo)
	rest.NewTopicHandler(kit.Echo.Group("/api/v1/topics"), topicSvc)

	// Now start the test server
	kit.Start(t)

	topicName := fmt.Sprintf("Topic-%s", uuid.New().String()[:8])
	createReq := domain.CreateTopicRequest{
		Name: topicName,
	}
	type CreateType domain.ResponseSingleData[domain.Topic]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/topics",
		createReq,
	)
	require.Equal(t, http.StatusCreated, code)
	topic := cre.Data
	require.NotEmpty(t, topic.ID)
	require.Equal(t, topicName, topic.Name)

	// Get List
	type ListType domain.ResponseMultipleData[domain.Topic]
	list, code := doRequest[ListType](
		t, http.MethodGet,
		kit.BaseURL+"/api/v1/topics",
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.True(t, len(list.Data) >= 1)

	// Get Single
	type GetType domain.ResponseSingleData[domain.Topic]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, topic.ID, getE.Data.ID)

	// Update
	updatedTopicName := fmt.Sprintf("Updated-%s", uuid.New().String()[:8])
	updPayload := domain.UpdateTopicRequest{
		Name: updatedTopicName,
	}
	type UpdType domain.ResponseSingleData[domain.Topic]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		updPayload,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, updatedTopicName, updE.Data.Name)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		nil,
	)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Get after delete
	_, code = doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		nil,
	)
	require.Equal(t, http.StatusNotFound, code)

	// Hard delete cleanup
	_, err = kit.DB.Exec(context.Background(), "DELETE FROM topic WHERE id = $1", topic.ID)
	require.NoError(t, err)
}
