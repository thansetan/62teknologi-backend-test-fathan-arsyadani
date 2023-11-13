package repository

import (
	"context"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(context.Context) ([]domain.Category, error)
	FindSomeByAlias(context.Context, []string) ([]domain.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (repo *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	err := repo.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *categoryRepository) FindSomeByAlias(ctx context.Context, aliases []string) ([]domain.Category, error) {
	var categories []domain.Category
	err := repo.db.WithContext(ctx).Where("alias IN ?", aliases).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
