package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/miruken1207/warehouse-api/internal/model"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(database *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: database}
}

func (r *ItemRepository) GetAll(ctx context.Context) ([]model.Item, error) {
	var items []model.Item
	query := `SELECT id, name, category FROM items`
	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("ItemRepository.GetAll: %w", err)
	}

	return items, nil
}
