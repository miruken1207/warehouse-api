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
	items := []model.Item{}
	query := `SELECT id, name, category FROM items`
	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("ItemRepository.GetAll: %w", err)
	}

	return items, nil
}

func (r *ItemRepository) GetItemByID(ctx context.Context, id int) (*model.Item, error) {
	var item model.Item
	query := `SELECT id, name, category FROM items WHERE id = $1`
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		return nil, fmt.Errorf("ItemRepository.GetItemByID: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) CreateItem(ctx context.Context, item *model.Item) error {
	query := `INSERT INTO items (name, category) VALUES ($1, $2) RETURNING id`
	if err := r.db.GetContext(ctx, &item.ID, query, item.Name, item.Category); err != nil {
		return fmt.Errorf("ItemRepository.CreateItem: %w", err)
	}

	return nil
}
