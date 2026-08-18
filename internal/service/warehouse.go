package service

import (
	"context"

	"github.com/miruken1207/warehouse-api/internal/model"
)

type WarehouseRepository interface {
	GetAll(ctx context.Context) ([]model.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id int) (*model.Warehouse, error)
}

type WarehouseService struct {
	repo WarehouseRepository
}

func NewWarehouseService(repository WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repository}
}

func (s *WarehouseService) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	return s.repo.GetAll(ctx)
}

func (s *WarehouseService) GetWarehouseByID(ctx context.Context, id int) (*model.Warehouse, error) {
	return s.repo.GetWarehouseByID(ctx, id)
}
