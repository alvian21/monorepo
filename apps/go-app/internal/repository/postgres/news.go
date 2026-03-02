package postgres

import (
	"context"
	"fmt"
	"go-app/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	Conn *pgxpool.Pool
}

func NewNewsRepository(conn *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{
		Conn: conn,
	}
}

func (r *NewsRepository) CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error) {
	query := `
		INSERT INTO news (title, content, status, created_at, updated_at)
		VALUES ($1, $2, 'DRAFT', NOW(), NOW())
		RETURNING id, title, content, status, created_at, updated_at`

	var news domain.News
	err := r.Conn.QueryRow(ctx, query, req.Title, req.Content).Scan(
		&news.ID,
		&news.Title,
		&news.Content,
		&news.Status,
		&news.CreatedAt,
		&news.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *NewsRepository) GetNewsList(ctx context.Context, filter *domain.NewsFilter) ([]domain.News, error) {
	query := `
		SELECT n.id, n.title, n.content, n.status, n.created_at, n.updated_at
		FROM news n
		WHERE n.deleted_at IS NULL`

	var args []interface{}
	idx := 1

	if filter != nil {
		if filter.Status != "" {
			query += fmt.Sprintf(" AND n.status = $%d", idx)
			args = append(args, filter.Status)
			idx++
		}
		if filter.Search != "" {
			query += fmt.Sprintf(" AND (n.title ILIKE $%d OR n.content ILIKE $%d)", idx, idx)
			args = append(args, "%"+filter.Search+"%")
			idx++
		}
	}

	rows, err := r.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []domain.News
	for rows.Next() {
		var n domain.News
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Status, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}

	return newsList, nil
}

func (r *NewsRepository) GetNews(ctx context.Context, id uuid.UUID) (*domain.News, error) {
	query := `
		SELECT id, title, content, status, created_at, updated_at
		FROM news
		WHERE id = $1 AND deleted_at IS NULL`

	var n domain.News
	err := r.Conn.QueryRow(ctx, query, id).Scan(&n.ID, &n.Title, &n.Content, &n.Status, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Load topics
	topicsQuery := `
		SELECT t.id, t.name, t.created_at, t.updated_at
		FROM topic t
		JOIN news_topics nt ON t.id = nt.topic_id
		WHERE nt.news_id = $1`

	rows, err := r.Conn.Query(ctx, topicsQuery, id)
	if err != nil {
		return &n, nil // Return news even if topics fail
	}
	defer rows.Close()

	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err == nil {
			n.Topics = append(n.Topics, t)
		}
	}

	return &n, nil
}

func (r *NewsRepository) UpdateNews(ctx context.Context, id uuid.UUID, req *domain.UpdateNewsRequest) (*domain.News, error) {
	query := `
		UPDATE news
		SET title = $1, content = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, title, content, status, created_at, updated_at`

	var n domain.News
	err := r.Conn.QueryRow(ctx, query, req.Title, req.Content, req.Status, id).Scan(
		&n.ID,
		&n.Title,
		&n.Content,
		&n.Status,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NewsRepository) DeleteNews(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE news
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	res, err := r.Conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("news not found")
	}

	return nil
}

func (r *NewsRepository) AssignTopics(ctx context.Context, newsID uuid.UUID, topicIDs []uuid.UUID) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing associations
	_, err = tx.Exec(ctx, "DELETE FROM news_topics WHERE news_id = $1", newsID)
	if err != nil {
		return err
	}

	if len(topicIDs) == 0 {
		return tx.Commit(ctx)
	}

	// Insert new associations
	query := "INSERT INTO news_topics (news_id, topic_id, created_at, updated_at) VALUES "
	vals := []interface{}{}
	for i, tid := range topicIDs {
		p1 := i*2 + 1
		p2 := i*2 + 2
		query += fmt.Sprintf("($%d, $%d, NOW(), NOW()),", p1, p2)
		vals = append(vals, newsID, tid)
	}
	query = strings.TrimSuffix(query, ",")

	_, err = tx.Exec(ctx, query, vals...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
