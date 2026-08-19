package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/miruken1207/warehouse-api/internal/model"
)

type WarehouseRepository struct {
	db *sqlx.DB
}

func NewWarehouseRepository(database *sqlx.DB) *WarehouseRepository {
	return &WarehouseRepository{db: database}
}

func (r *WarehouseRepository) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	var warehouses []model.Warehouse

	query := `SELECT id, name, location FROM warehouses ORDER BY id`

	if err := r.db.SelectContext(ctx, &warehouses, query); err != nil {
		return nil, fmt.Errorf("WarehouseRepository.GetALL: %w", err)
	}

	return warehouses, nil
}

func (r *WarehouseRepository) GetWarehouseByID(ctx context.Context, id int) (*model.Warehouse, error) {
	var warehouse model.Warehouse

	query := `SELECT id, name, location FROM warehouses WHERE id = $1"`
	if err := r.db.GetContext(ctx, &warehouse, query, id); err != nil {
		return nil, fmt.Errorf("WarehouseRepository.GetWarehouseByID: %w", err)
	}

	return &warehouse, nil
}

func (r *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	query := `INSERT INTO warehouses(name, location) VALUES ($1, $2) RETURNING id`
	if err := r.db.GetContext(ctx, &warehouse.ID, query, warehouse.Name, warehouse.Location); err != nil {
		return fmt.Errorf("WarehouseRepository.CreateWarehouse: %w", err)
	}
	return nil
}
