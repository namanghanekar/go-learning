package service

import (
	"ticket-system/inventory-service/models"
	"ticket-system/inventory-service/repository"
)

type Service struct {
	repo *repository.Repo
}

func NewService(r *repository.Repo) *Service {
	return &Service{repo: r}
}

// GET ALL
func (s *Service) GetSeats() []models.Seat {
	return s.repo.GetAll()
}

// LOCK
func (s *Service) Lock(id int, user string) error {
	return s.repo.LockSeat(id, user)
}

// UNLOCK
func (s *Service) Unlock(id int) error {
	return s.repo.UnlockSeat(id)
}

// CONFIRM
func (s *Service) Confirm(seatID int, userID string) error {
	if err := s.repo.ConfirmSeat(seatID, userID); err != nil {
		return err
	}
	return nil
}

func (s *Service) SeedSeats(total int) error {
	return s.repo.SeedSeats(total)
}
