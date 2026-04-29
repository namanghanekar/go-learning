package store

import (
	"errors"
	"fmt"
	"time"
	"worldtour-tickets/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) SeedRoom(roomID string, total int) error {
	for i := 1; i <= total; i++ {
		number := seatNumber(i)
		seat := domain.Seat{
			ID:     roomID + "-" + number,
			RoomID: roomID,
			Number: number,
			Status: domain.SeatStatusAvailable,
		}

		if err := r.db.Where("room_id = ? AND number = ?", roomID, seat.Number).FirstOrCreate(&seat).Error; err != nil {
			return err
		}
	}
	return nil
}

func seatNumber(n int) string {
	row := 'A' + rune((n-1)/10)
	col := ((n - 1) % 10) + 1
	return fmt.Sprintf("%c%02d", row, col)
}

func (r *SeatRepository) List(roomID string) ([]domain.Seat, error) {
	var seats []domain.Seat
	err := r.db.Where("room_id = ?", roomID).Order("number asc").Find(&seats).Error
	return seats, err
}

func (r *SeatRepository) HoldSeat(roomID, seatNumber, userID, holdToken string, expiresAt time.Time) (domain.Seat, error) {
	var seat domain.Seat
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("room_id = ? AND number = ?", roomID, seatNumber).
			First(&seat).Error; err != nil {
			return err
		}

		if seat.Status == domain.SeatStatusSold {
			return errors.New("seat already sold")
		}

		if seat.Status == domain.SeatStatusHeld && seat.ExpiresAt != nil && seat.ExpiresAt.After(time.Now()) {
			return errors.New("seat is already held")
		}

		seat.Status = domain.SeatStatusHeld
		seat.HeldBy = userID
		seat.HoldToken = holdToken
		seat.ExpiresAt = &expiresAt

		return tx.Save(&seat).Error
	})
	return seat, err
}

func (r *SeatRepository) ConfirmSeat(seatID, holdToken string) (domain.Seat, error) {
	var seat domain.Seat
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", seatID).
			First(&seat).Error; err != nil {
			return err
		}

		if seat.HoldToken != holdToken {
			return errors.New("hold token mismatch")
		}

		seat.Status = domain.SeatStatusSold
		seat.ExpiresAt = nil

		return tx.Save(&seat).Error
	})
	return seat, err
}

func (r *SeatRepository) ReleaseSeat(seatID, holdToken string) (domain.Seat, error) {
	var seat domain.Seat
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", seatID).
			First(&seat).Error; err != nil {
			return err
		}

		if holdToken != "" && seat.HoldToken != holdToken {
			return errors.New("hold token mismatch")
		}

		if seat.Status == domain.SeatStatusSold {
			return errors.New("sold seats cannot be released")
		}

		seat.Status = domain.SeatStatusAvailable
		seat.HeldBy = ""
		seat.HoldToken = ""
		seat.ExpiresAt = nil

		return tx.Save(&seat).Error
	})
	return seat, err
}
