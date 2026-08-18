package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/miruken1207/warehouse-api/internal/model"
)

type StockRepository struct {
	db *sqlx.DB
}

func NewStockRepository(database *sqlx.DB) *StockRepository {
	return &StockRepository{db: database}
}

func (r *StockRepository) GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error) {
	var stock []model.Stock
	query := `SELECT id, warehouse_id, item_id, quantity FROM stock WHERE warehouse_id = $1`
	if err := r.db.SelectContext(ctx, &stock, query, id); err != nil {
		return nil, fmt.Errorf("StockRepository.GetStockByWarehouseID: %w", err)
	}
	return stock, nil
}
