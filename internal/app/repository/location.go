package repository

import (
	"context"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"gorm.io/gorm"
)

type LocationRepository interface {
	Create(context.Context, *gorm.DB, domain.Location) (domain.Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{
		db: db,
	}
}

func (repo *locationRepository) Create(ctx context.Context, tx *gorm.DB, data domain.Location) (domain.Location, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
