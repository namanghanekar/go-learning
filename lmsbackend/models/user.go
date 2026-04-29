package models

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Role       string `gorm:"not null"`
	IsVerified bool   `gorm:"default:false"`

	// Relations
	CreatorProfile CreatorProfile `gorm:"foreignKey:UserID"`
	BrandProfile   BrandProfile   `gorm:"foreignKey:UserID"`
}
