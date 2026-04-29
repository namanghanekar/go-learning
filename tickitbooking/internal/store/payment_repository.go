package store

import (
	"time"
	"worldtour-tickets/internal/domain"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) UpsertPending(record domain.PaymentRecord) error {
	record.RequestedAt = time.Now().UTC()
	return r.db.Where("id = ?", record.ID).Assign(record).FirstOrCreate(&record).Error
}

func (r *PaymentRepository) MarkProcessed(id, status, reason string) error {
	now := time.Now().UTC()
	return r.db.Model(&domain.PaymentRecord{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":         status,
			"failure_reason": reason,
			"processed_at":   &now,
		}).Error
}
