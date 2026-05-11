package models

type CreatorProfile struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"unique"`
	Name        string
	Platform    string
	Followers   string
	ProfileLink string
}
