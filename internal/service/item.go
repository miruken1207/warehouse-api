package service

import (
	"context"

	"github.com/miruken1207/warehouse-api/internal/model"
)

type ItemRepository interface {
	GetAll(ctx context.Context) ([]model.Item, error)
}

type ItemService struct {
	repo ItemRepository
}

func NewItemService(r ItemRepository) *ItemService {
	return &ItemService{repo: r}
}

func (s *ItemService) GetAll(ctx context.Context) ([]model.Item, error) {
	return s.repo.GetAll(ctx)
}
