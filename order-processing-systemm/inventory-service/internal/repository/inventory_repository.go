package repository

import (
	"order-processing-system/inventory-service/internal/model"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	DB *gorm.DB
}

func NewInventoryRepository(
	db *gorm.DB,
) *InventoryRepository {

	return &InventoryRepository{
		DB: db,
	}
}

func (r *InventoryRepository) CheckStock(
	productID int,
	quantity int,
) bool {

	var inventory model.Inventory

	err := r.DB.Where(
		"product_id = ?",
		productID,
	).First(&inventory).Error

	if err != nil {
		return false
	}

	return inventory.Stock >= quantity
}
