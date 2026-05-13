package service

import (
	"order-processing-system/inventory-service/internal/repository"
)

type InventoryService struct {
	Repo *repository.InventoryRepository
}

func NewInventoryService(
	repo *repository.InventoryRepository,
) *InventoryService {

	return &InventoryService{
		Repo: repo,
	}
}

func (s *InventoryService) CheckStock(
	productID int,
	quantity int,
) bool {

	return s.Repo.CheckStock(
		productID,
		quantity,
	)
}
