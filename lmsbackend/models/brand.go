package models

type BrandProfile struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"unique"` // one-to-one
	CompanyName string
	Website     string
	GST         string
	Description string
}
