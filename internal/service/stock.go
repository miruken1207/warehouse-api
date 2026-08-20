package service

import (
	"context"
	"fmt"

	"github.com/miruken1207/warehouse-api/internal/model"
	"github.com/miruken1207/warehouse-api/internal/repository"
)

type StockRepository interface {
	GetAll(ctx context.Context) ([]model.Stock, error)
	GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error)
	GetStockByItemID(ctx context.Context, id int) ([]model.Stock, error)
	CreateStock(ctx context.Context, stock *model.Stock) error
	UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error
	BeginTx(ctx context.Context) (repository.StockTx, error)
}

type StockService struct {
	repo StockRepository
}

func NewStockService(r StockRepository) *StockService {
	return &StockService{repo: r}
}

func (s *StockService) GetAll(ctx context.Context) ([]model.Stock, error) {
	return s.repo.GetAll(ctx)
}

func (s *StockService) GetStockByWarehouseID(ctx context.Context, id int) ([]model.Stock, error) {
	return s.repo.GetStockByWarehouseID(ctx, id)
}

func (s *StockService) GetStockByItemID(ctx context.Context, id int) ([]model.Stock, error) {
	return s.repo.GetStockByItemID(ctx, id)
}

func (s *StockService) CreateStock(ctx context.Context, stock *model.Stock) error {
	return s.repo.CreateStock(ctx, stock)
}

func (s *StockService) UpdateStock(ctx context.Context, updateReq *model.UpdateStockRequest) error {
	return s.repo.UpdateStock(ctx, updateReq)
}

func (s *StockService) TransferStock(ctx context.Context, transferReq *model.TransferStockRequest) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("StockService.TransferStock: %w", err)
	}
	defer tx.Rollback()

	if err := tx.UpdateStock(ctx, &model.UpdateStockRequest{
		WarehouseID: transferReq.FromWarehouseID,
		ItemID:      transferReq.ItemID,
		Delta:       -transferReq.Quantity,
	}); err != nil {
		return fmt.Errorf("StockService.TransferStock: %w", err)
	}

	if err := tx.UpdateStock(ctx, &model.UpdateStockRequest{
		WarehouseID: transferReq.ToWarehouseID,
		ItemID:      transferReq.ItemID,
		Delta:       +transferReq.Quantity,
	}); err != nil {
		return fmt.Errorf("StockService.TransferStock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("StockService.TransferStock: %w", err)
	}

	return nil
}
