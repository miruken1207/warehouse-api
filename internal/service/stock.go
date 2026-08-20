package service

import (
	"context"

	"github.com/miruken1207/warehouse-api/internal/model"
)

type StockRepository interface {
	GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error)
	GetStockByItemID(ctx context.Context, id int) ([]model.Stock, error)
	UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error
}

type StockService struct {
	repo StockRepository
}

func NewStockService(r StockRepository) *StockService {
	return &StockService{repo: r}
}

func (s *StockService) GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error) {
	return s.repo.GetStockByWarehouseID(ctx, id)
}

func (s *StockService) GetStockByItemID(ctx context.Context, id int) ([]model.Stock, error) {
	return s.repo.GetStockByItemID(ctx, id)
}

func (s *StockService) UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error {
	return s.repo.UpdateStock(ctx, updateReq)
}
