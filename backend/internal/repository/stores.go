package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"goapi/backend/internal/model"
)

var ErrDuplicateSlug = errors.New("slug already exists")

type StoreRepository struct {
	db *pgxpool.Pool
}

func NewStoreRepository(db *pgxpool.Pool) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Create(ctx context.Context, ownerID, name, slug string) (model.Store, error) {
	const query = `
		INSERT INTO stores (owner_id, name, slug)
		VALUES ($1, $2, $3)
		RETURNING id::text, owner_id::text, name, slug, created_at, updated_at`

	var store model.Store
	err := r.db.QueryRow(ctx, query, ownerID, name, slug).Scan(
		&store.ID,
		&store.OwnerID,
		&store.Name,
		&store.Slug,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Store{}, ErrDuplicateSlug
		}
		return model.Store{}, err
	}

	return store, nil
}

func (r *StoreRepository) FindBySlug(ctx context.Context, slug string) (model.Store, error) {
	const query = `
		SELECT id::text, owner_id::text, name, slug, created_at, updated_at
		FROM stores
		WHERE slug = $1`

	var store model.Store
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&store.ID,
		&store.OwnerID,
		&store.Name,
		&store.Slug,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Store{}, ErrNotFound
	}
	if err != nil {
		return model.Store{}, err
	}

	return store, nil
}

func (r *StoreRepository) FindByOwner(ctx context.Context, ownerID string) (model.Store, error) {
	const query = `
		SELECT id::text, owner_id::text, name, slug, created_at, updated_at
		FROM stores
		WHERE owner_id = $1`

	var store model.Store
	err := r.db.QueryRow(ctx, query, ownerID).Scan(
		&store.ID,
		&store.OwnerID,
		&store.Name,
		&store.Slug,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Store{}, ErrNotFound
	}
	if err != nil {
		return model.Store{}, err
	}

	return store, nil
}
