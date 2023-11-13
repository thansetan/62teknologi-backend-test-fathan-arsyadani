package usecase

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/repository"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
)

type TransactionUsecase interface {
	GetAll(ctx context.Context) ([]dto.Transaction, *helpers.ResponseError)
}

type transactionUsecase struct {
	repo   repository.TransactionRepository
	logger *slog.Logger
}

func NewTransactionUsecase(repo repository.TransactionRepository, logger *slog.Logger) TransactionUsecase {
	return &transactionUsecase{repo, logger}
}

func (uc *transactionUsecase) GetAll(ctx context.Context) ([]dto.Transaction, *helpers.ResponseError) {
	transactions, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("Transaction Usecase", "function", helpers.GetFunctionName(), "err", err)
		return nil, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	data := make([]dto.Transaction, len(transactions))
	for i, transaction := range transactions {
		data[i] = dto.Transaction{
			Type: transaction.Type,
		}
	}

	return data, nil
}
