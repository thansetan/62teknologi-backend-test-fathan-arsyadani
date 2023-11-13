package repository

import (
	"context"
	"fmt"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	Create(context.Context, domain.Business) error
	FindByIDOrAlias(context.Context, string) (domain.Business, error)
	Update(context.Context, domain.Business) error
	Delete(context.Context, domain.Business) error
	FindByParams(context.Context, domain.BusinessQuery) ([]domain.Business, int64, error)
}

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{db}
}

func (repo *businessRepository) Create(ctx context.Context, data domain.Business) error {
	err := repo.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *businessRepository) FindByIDOrAlias(ctx context.Context, IDorAlias string) (domain.Business, error) {
	var business domain.Business
	err := repo.db.WithContext(ctx).Preload("Categories").Preload("Transactions").Preload("Location").First(&business, "id = ? OR alias = ?", IDorAlias, IDorAlias).Error
	if err != nil {
		return business, err
	}

	return business, nil
}

func (repo *businessRepository) Update(ctx context.Context, data domain.Business) error {
	tx := repo.db.Begin().WithContext(ctx)
	err := tx.Updates(&data.Location).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Updates(&data).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&data).Association("Categories").Replace(data.Categories)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&data).Association("Transactions").Replace(data.Transactions)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (repo *businessRepository) Delete(ctx context.Context, data domain.Business) error {
	return repo.db.Delete(&data).Error
}

func (repo *businessRepository) FindByParams(ctx context.Context, params domain.BusinessQuery) ([]domain.Business, int64, error) {
	var (
		businesses []domain.Business
		total      int64
	)

	db := repo.db.WithContext(ctx)
	if params.Categories != nil {
		db = db.Where("c.alias IN ?", params.Categories)
	}

	if params.TransactionType != nil {
		db = db.Where("t.type IN ?", params.TransactionType)
	}

	if params.OpenNow != "" || params.OpenAt != "" {
		db = db.Where("? BETWEEN open_at AND close_at", func(a, b string) string {
			if a == "" {
				return b
			}
			return a
		}(params.OpenNow, params.OpenAt))
	}

	err := db.Preload("Location").Preload("Categories").Preload("Transactions").
		Joins("LEFT JOIN businesses_categories bc ON businesses.id=bc.business_id").
		Joins("FULL JOIN categories c ON c.id=bc.category_id").
		Joins("LEFT JOIN businesses_transactions bt ON businesses.id=bt.business_id").
		Joins("FULL JOIN transactions t ON t.id=bt.transaction_id").
		Joins("LEFT JOIN locations l ON businesses.id=l.business_id").
		Where("LIKE(LOWER(CONCAT(l.address1, ' ', l.address2, ' ', l.address3, ' ', l.city, ' ', l.state, ' ', l.country)), ?)", fmt.
			Sprintf("%%%s%%", params.Location)).
		Group("businesses.id").Order("businesses.id ASC").Model(&domain.Business{}).Count(&total).Offset(params.Offset).Limit(params.Limit).
		Find(&businesses).Error

	if err != nil {
		return nil, total, err
	}

	return businesses, total, nil
}
