package model

type Inventory struct {
	ID        uint `gorm:"primaryKey"`
	ProductID int
	Stock     int
}
