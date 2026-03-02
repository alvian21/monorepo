package service

import (
	"context"
	"go-app/domain"

	"github.com/google/uuid"
)

type NewsRepository interface {
	CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error)
	GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error)
	GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error)
	UpdateNews(ctx context.Context, id uuid.UUID, req *domain.UpdateNewsRequest) (*domain.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error
	AssignTopics(ctx context.Context, newsID uuid.UUID, topicIDs []uuid.UUID) error
}

type NewsService struct {
	newsRepo  NewsRepository
	topicRepo TopicRepository
}

func NewNewsService(newsRepo NewsRepository, topicRepo TopicRepository) *NewsService {
	return &NewsService{
		newsRepo:  newsRepo,
		topicRepo: topicRepo,
	}
}

func (s *NewsService) CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error) {
	news, err := s.newsRepo.CreateNews(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(req.TopicIDs) > 0 {
		err = s.newsRepo.AssignTopics(ctx, news.ID, req.TopicIDs)
		if err != nil {
			return nil, err
		}
	}

	return s.GetNews(ctx, news.ID)
}

func (s *NewsService) GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error) {
	return s.newsRepo.GetNewsList(ctx, filter)
}

func (s *NewsService) GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error) {
	return s.newsRepo.GetNews(ctx, id)
}

func (s *NewsService) UpdateNews(ctx context.Context, id uuid.UUID, req *domain.UpdateNewsRequest) (*domain.News, error) {
	_, err := s.newsRepo.UpdateNews(ctx, id, req)
	if err != nil {
		return nil, err
	}

	err = s.newsRepo.AssignTopics(ctx, id, req.TopicIDs)
	if err != nil {
		return nil, err
	}

	return s.GetNews(ctx, id)
}

func (s *NewsService) DeleteNews(ctx context.Context, id uuid.UUID) error {
	return s.newsRepo.DeleteNews(ctx, id)
}
