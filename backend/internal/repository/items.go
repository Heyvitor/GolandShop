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

func (r *ItemRepository) Create(ctx context.Context, item model.Item) (model.Item, error) {
	const query = `
		INSERT INTO items (user_id, store_id, title, body, name, description, price, variant, variant_price, shipping_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id::text, user_id::text, store_id::text, name, description, price, variant, variant_price, shipping_type, created_at, updated_at`

	var created model.Item
	err := r.db.QueryRow(
		ctx,
		query,
		item.UserID,
		item.StoreID,
		item.Name,
		item.Description,
		item.Name,
		item.Description,
		item.Price,
		item.Variant,
		item.VariantPrice,
		item.ShippingType,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.StoreID,
		&created.Name,
		&created.Description,
		&created.Price,
		&created.Variant,
		&created.VariantPrice,
		&created.ShippingType,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return model.Item{}, err
	}

	return created, nil
}

func (r *ItemRepository) ListByUser(ctx context.Context, userID string, limit int32) ([]model.Item, error) {
	const query = `
		SELECT id::text, user_id::text, store_id::text, name, description, price, variant, variant_price, shipping_type, created_at, updated_at
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
			&item.StoreID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Variant,
			&item.VariantPrice,
			&item.ShippingType,
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
