package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *StockRepository) GetAll(ctx context.Context) ([]model.Stock, error) {
	stock := []model.Stock{}
	query := `SELECT id, warehouse_id, item_id, quantity FROM stock`
	if err := r.db.SelectContext(ctx, &stock, query); err != nil {
		return nil, fmt.Errorf("StockRepository.GetAll: %w", err)
	}

	return stock, nil
}

func (r *StockRepository) GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error) {
	var exists bool
	if err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM warehouses WHERE id = $1)`, id); err != nil {
		return nil, fmt.Errorf("StockRepository.GetStockByWarehouseID: %w", err)
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	stock := []model.Stock{}
	query := `SELECT id, warehouse_id, item_id, quantity FROM stock WHERE warehouse_id = $1`
	if err := r.db.SelectContext(ctx, &stock, query, id); err != nil {
		return nil, fmt.Errorf("StockRepository.GetStockByWarehouseID: %w", err)
	}

	return stock, nil
}

func (r *StockRepository) GetStockByItemID(ctx context.Context, id int) ([]model.Stock, error) {
	var exists bool
	if err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)`, id); err != nil {
		return nil, fmt.Errorf("StockRepository.GetStockByItemID: %w", err)
	}
	if !exists {
		return nil, sql.ErrNoRows
	}

	stock := []model.Stock{}
	query := `SELECT id, warehouse_id, item_id, quantity FROM stock WHERE item_id = $1`
	if err := r.db.SelectContext(ctx, &stock, query, id); err != nil {
		return nil, fmt.Errorf("StockRepository.GetStockByItemID: %w", err)
	}

	return stock, nil
}

var (
	ErrStockNotFound     = errors.New("stock record not found")
	ErrInsufficientStock = errors.New("insufficient stock quantity")
)

func (r *StockRepository) CreateStock(ctx context.Context, stock *model.Stock) error {
	query := `INSERT INTO stock (warehouse_id, item_id, quantity) VALUES ($1, $2, $3) RETURNING id`
	if err := r.db.GetContext(ctx, &stock.ID, query, stock.WarehouseID, stock.ItemID, stock.Quantity); err != nil {
		return fmt.Errorf("StockRepository.CreateStock: %w", err)
	}

	return nil
}

func (r *StockRepository) UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error {
	query := `UPDATE stock SET quantity = quantity + $1 WHERE warehouse_id = $2 AND item_id = $3 AND quantity + $1 >= 0`

	res, err := r.db.ExecContext(ctx, query, updateReq.Delta, updateReq.WarehouseID, updateReq.ItemID)
	if err != nil {
		return fmt.Errorf("StockRepository.UpdateStock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("StockRepository.UpdateStock: %w", err)
	}

	if rowsAffected == 0 {
		query = `SELECT id, warehouse_id, item_id, quantity FROM stock WHERE warehouse_id = $1 AND item_id = $2`
		var stock model.Stock
		if err := r.db.GetContext(ctx, &stock, query, updateReq.WarehouseID, updateReq.ItemID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrStockNotFound
			}
			return fmt.Errorf("StockRepository.UpdateStock: %w", err)
		}
		return ErrInsufficientStock
	}

	return nil
}
