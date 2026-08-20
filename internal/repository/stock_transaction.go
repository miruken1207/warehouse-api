package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/miruken1207/warehouse-api/internal/model"
)

type StockTx interface {
	UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error
	Commit() error
	Rollback() error
}

func (r *StockRepository) BeginTx(ctx context.Context) (StockTx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("StockRepository.BeginTx: %w", err)
	}
	return &stockTx{tx: tx}, err
}

type stockTx struct {
	tx *sqlx.Tx
}

func (t *stockTx) UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error {
	query := `UPDATE stock SET quantity = quantity + $1 WHERE warehouse_id = $2 AND item_id = $3 AND quantity + $1 >= 0`

	res, err := t.tx.ExecContext(ctx, query, updateReq.Delta, updateReq.WarehouseID, updateReq.ItemID)
	if err != nil {
		return fmt.Errorf("stockTx.UpdateStock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("stockTx.UpdateStock: %w", err)
	}

	if rowsAffected == 0 {
		query = `SELECT id, warehouse_id, item_id, quantity FROM stock WHERE warehouse_id = $1 AND item_id = $2`
		var stock model.Stock
		if err := t.tx.GetContext(ctx, &stock, query, updateReq.WarehouseID, updateReq.ItemID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrStockNotFound
			}
			return fmt.Errorf("stockTx.UpdateStock: %w", err)
		}
		return ErrInsufficientStock
	}

	return nil
}

func (t *stockTx) Commit() error {
	return t.tx.Commit()
}

func (t *stockTx) Rollback() error {
	return t.tx.Rollback()
}
