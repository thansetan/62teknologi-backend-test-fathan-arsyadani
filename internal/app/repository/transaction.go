package repository

import (
	"context"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindAll(context.Context) ([]domain.Transaction, error)
	FindSome(context.Context, []string) ([]domain.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (repo *transactionRepository) FindAll(ctx context.Context) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := repo.db.WithContext(ctx).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *transactionRepository) FindSome(ctx context.Context, types []string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := repo.db.WithContext(ctx).Where("type IN ?", types).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
