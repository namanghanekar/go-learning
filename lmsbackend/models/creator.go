package models

type CreatorProfile struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"unique"` // one-to-one
	Name        string
	Platform    string
	Followers   string
	ProfileLink string
}
