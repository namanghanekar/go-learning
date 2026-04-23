package models

type Task struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `json:"title"`
	UserID uint
}
