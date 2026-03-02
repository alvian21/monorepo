package postgres

import (
	"context"
	"fmt"
	"go-app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicRepository struct {
	Conn *pgxpool.Pool
}

func NewTopicRepository(conn *pgxpool.Pool) *TopicRepository {
	return &TopicRepository{
		Conn: conn,
	}
}

func (r *TopicRepository) CreateTopic(ctx context.Context, req *domain.CreateTopicRequest) (*domain.Topic, error) {
	query := `
		INSERT INTO topic (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, name, created_at, updated_at`

	var topic domain.Topic
	err := r.Conn.QueryRow(ctx, query, req.Name).Scan(
		&topic.ID,
		&topic.Name,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func (r *TopicRepository) GetTopicList(ctx context.Context, filter *domain.TopicFilter) ([]domain.Topic, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM topic
		WHERE deleted_at IS NULL`

	var args []interface{}
	if filter != nil && filter.Search != "" {
		query += " AND name ILIKE $1"
		args = append(args, "%"+filter.Search+"%")
	}

	rows, err := r.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

func (r *TopicRepository) GetTopic(ctx context.Context, id uuid.UUID) (*domain.Topic, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM topic
		WHERE id = $1 AND deleted_at IS NULL`

	var t domain.Topic
	err := r.Conn.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TopicRepository) UpdateTopic(ctx context.Context, id uuid.UUID, req *domain.UpdateTopicRequest) (*domain.Topic, error) {
	query := `
		UPDATE topic
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, name, created_at, updated_at`

	var t domain.Topic
	err := r.Conn.QueryRow(ctx, query, req.Name, id).Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TopicRepository) DeleteTopic(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE topic
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	res, err := r.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("topic not found")
	}

	return nil
}

func (r *TopicRepository) GetTopicsByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Topic, error) {
	if len(ids) == 0 {
		return []domain.Topic{}, nil
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM topic
		WHERE id = ANY($1) AND deleted_at IS NULL`

	rows, err := r.Conn.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}
