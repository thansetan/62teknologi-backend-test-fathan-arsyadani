package usecase

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/repository"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
)

type CategoryUsecase interface {
	GetAll(context.Context) ([]dto.Category, error)
}

type categoryUsecase struct {
	repo   repository.CategoryRepository
	logger *slog.Logger
}

func NewCategoryUsecase(repo repository.CategoryRepository, logger *slog.Logger) CategoryUsecase {
	return &categoryUsecase{repo, logger}
}

func (uc *categoryUsecase) GetAll(ctx context.Context) ([]dto.Category, error) {
	categories, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error("Category Usecase", "function", helpers.GetFunctionName(), "err", err)
		return nil, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	data := make([]dto.Category, len(categories))
	for i, category := range categories {
		data[i] = dto.Category{
			Title: category.Title,
			Alias: category.Alias,
		}
	}

	return data, nil
}
