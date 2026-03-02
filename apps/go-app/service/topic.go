package service

import (
	"context"
	"go-app/domain"

	"github.com/google/uuid"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, req *domain.CreateTopicRequest) (*domain.Topic, error)
	GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error)
	GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id uuid.UUID, req *domain.UpdateTopicRequest) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id uuid.UUID) error
	GetTopicsByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Topic, error)
}

type TopicService struct {
	repo TopicRepository
}

func NewTopicService(repo TopicRepository) *TopicService {
	return &TopicService{
		repo: repo,
	}
}

func (s *TopicService) CreateTopic(ctx context.Context, req *domain.CreateTopicRequest) (*domain.Topic, error) {
	return s.repo.CreateTopic(ctx, req)
}

func (s *TopicService) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	return s.repo.GetTopicList(ctx, filter)
}

func (s *TopicService) GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	return s.repo.GetTopic(ctx, id)
}

func (s *TopicService) UpdateTopic(ctx context.Context, id uuid.UUID, req *domain.UpdateTopicRequest) (*domain.Topic, error) {
	return s.repo.UpdateTopic(ctx, id, req)
}

func (s *TopicService) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTopic(ctx, id)
}
