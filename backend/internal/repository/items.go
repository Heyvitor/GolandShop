package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"goapi/backend/internal/model"
)

var ErrNotFound = errors.New("not found")

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, userID, title, body string) (model.Item, error) {
	const query = `
		INSERT INTO items (user_id, title, body)
		VALUES ($1, $2, $3)
		RETURNING id::text, user_id::text, title, body, created_at, updated_at`

	var item model.Item
	err := r.db.QueryRow(ctx, query, userID, title, body).Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.Body,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (r *ItemRepository) ListByUser(ctx context.Context, userID string, limit int32) ([]model.Item, error) {
	const query = `
		SELECT id::text, user_id::text, title, body, created_at, updated_at
		FROM items
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Item, 0, limit)
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Title,
			&item.Body,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
