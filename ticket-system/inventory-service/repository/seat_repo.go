package repository

import (
	"ticket-system/inventory-service/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// LOCK
func (r *Repo) LockSeat(id int, user string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var seat models.Seat

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&seat, id).Error; err != nil {
			return err
		}

		if seat.Status != "available" {
			return gorm.ErrInvalidData
		}

		expiresAt := time.Now().Add(5 * time.Minute)
		seat.Status = "locked"
		seat.UserID = user
		seat.LockExpiresAt = &expiresAt

		return tx.Save(&seat).Error
	})
}

// GET BY ID
func (r *Repo) GetByID(id int) (models.Seat, error) {
	var seat models.Seat
	err := r.db.First(&seat, id).Error
	return seat, err
}

// confirm seat
func (r *Repo) ConfirmSeat(id int, user string) error {
	now := time.Now()
	return r.db.Model(&models.Seat{}).
		Where("id = ? AND status = ? AND user_id = ? AND lock_expires_at > ?", id, "locked", user, now).
		Updates(map[string]interface{}{
			"status":          "booked",
			"lock_expires_at": nil,
		}).Error
}

// UNLOCK
func (r *Repo) UnlockSeat(id int) error {
	return r.db.Model(&models.Seat{}).
		Where("id = ? AND status = ?", id, "locked").
		Updates(map[string]interface{}{
			"status":          "available",
			"user_id":         "",
			"lock_expires_at": nil,
		}).Error
}
func (r *Repo) GetSeat(id int) models.Seat {
	var seat models.Seat
	r.db.First(&seat, id)
	return seat
}

func (r *Repo) Update(seat models.Seat) {
	r.db.Save(&seat)
}

// GET ALL
func (r *Repo) GetAll() []models.Seat {
	var seats []models.Seat
	r.ReleaseExpiredLocks()
	r.db.Find(&seats)
	return seats
}

func (r *Repo) ReleaseExpiredLocks() error {
	now := time.Now()
	return r.db.Model(&models.Seat{}).
		Where("status = ? AND lock_expires_at IS NOT NULL AND lock_expires_at <= ?", "locked", now).
		Updates(map[string]interface{}{
			"status":          "available",
			"user_id":         "",
			"lock_expires_at": nil,
		}).Error
}

func (r *Repo) SeedSeats(total int) error {
	var count int64
	if err := r.db.Model(&models.Seat{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seats := make([]models.Seat, 0, total)
	for i := 1; i <= total; i++ {
		seats = append(seats, models.Seat{
			ID:     i,
			Status: "available",
		})
	}

	return r.db.Create(&seats).Error
}
