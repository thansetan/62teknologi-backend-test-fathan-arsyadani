package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/segmentio/ksuid"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/repository"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
	"gorm.io/gorm"
)

type BusinessUsecase interface {
	Create(context.Context, dto.BusinessRequest) (dto.CreateBusinessResponse, error)
	Update(context.Context, string, dto.BusinessRequest) error
	GetBusiness(context.Context, string) (dto.BusinessResponse, error)
	Delete(context.Context, string) error
	GetBusinesses(context.Context, dto.BusinessQueryParams) (dto.BusinessesResponse, error)
}

type businessUsecase struct {
	businessRepo    repository.BusinessRepository
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	db              *gorm.DB
	logger          *slog.Logger
}

func NewBusinessUsecase(businessRepo repository.BusinessRepository, transactionRepo repository.TransactionRepository, categoryRepo repository.CategoryRepository, db *gorm.DB, logger *slog.Logger) *businessUsecase {
	return &businessUsecase{businessRepo, transactionRepo, categoryRepo, db, logger}
}

func (uc *businessUsecase) Create(ctx context.Context, data dto.BusinessRequest) (dto.CreateBusinessResponse, error) {
	var res dto.CreateBusinessResponse
	if err := data.ValidateCreate(); err != nil {
		return res, helpers.NewError(helpers.NewValdiationError(err), http.StatusBadRequest)
	}

	var (
		categoriesData   []domain.Category
		transactionsData []domain.Transaction
		err              error
	)

	if data.Categories != nil {
		categories := strings.Split(*data.Categories, ",")
		categoriesData, err = uc.categoryRepo.FindSomeByAlias(ctx, categories)
		if err != nil {
			uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
			return res, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
		}

		if *data.Categories != "" && len(categories) != len(categoriesData) {
			return res, helpers.NewError(helpers.ErrCategories, http.StatusBadRequest)
		}

	}

	if data.Transactions != nil {
		transactions := strings.Split(*data.Transactions, ",")
		transactionsData, err = uc.transactionRepo.FindSome(ctx, transactions)
		if err != nil {
			uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
			return res, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
		}

		if *data.Transactions != "" && len(transactions) != len(transactionsData) {
			return res, helpers.NewError(helpers.ErrTransactions, http.StatusBadRequest)
		}
	}

	alias := slug.Make(fmt.Sprintf("%s %s %d", data.Name, data.City, time.Now().UnixMicro()))

	business := domain.Business{
		ID:       ksuid.New().String(),
		Name:     data.Name,
		Alias:    alias,
		ImageURL: data.ImageURL,
		Location: domain.Location{
			Address1:  data.Address1,
			Address2:  data.Address2,
			Address3:  data.Address3,
			City:      data.City,
			ZipCode:   data.ZipCode,
			Country:   data.Country,
			State:     data.State,
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		},
		Price:        data.Price,
		Phone:        data.Phone,
		OpenAt:       data.OpenAt,
		CloseAt:      data.CloseAt,
		Categories:   categoriesData,
		Transactions: transactionsData,
	}

	err = uc.businessRepo.Create(ctx, business)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return res, helpers.NewError(errors.New("the provided phone number is already used by another business"), http.StatusConflict)
		}
		return res, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	res.Alias = business.Alias

	return res, nil
}

func (uc *businessUsecase) Update(ctx context.Context, id string, data dto.BusinessRequest) error {
	if err := data.ValidateUpdate(); err != nil {
		return helpers.NewError(helpers.NewValdiationError(err), http.StatusBadRequest)
	}

	var (
		categoriesData   []domain.Category
		transactionsData []domain.Transaction
		err              error
	)

	if data.Categories != nil {
		categories := strings.Split(*data.Categories, ",")
		categoriesData, err = uc.categoryRepo.FindSomeByAlias(ctx, categories)
		if err != nil {
			uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
			return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
		}

		if *data.Categories != "" && len(categories) != len(categoriesData) {
			return helpers.NewError(helpers.ErrCategories, http.StatusBadRequest)
		}

	}

	if data.Transactions != nil {
		transactions := strings.Split(*data.Transactions, ",")
		transactionsData, err = uc.transactionRepo.FindSome(ctx, transactions)
		if err != nil {
			uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
			return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
		}

		if *data.Transactions != "" && len(transactions) != len(transactionsData) {
			return helpers.NewError(helpers.ErrTransactions, http.StatusBadRequest)
		}
	}

	business, err := uc.businessRepo.FindByIDOrAlias(ctx, id)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.NewError(helpers.ErrBusinessNotFound, http.StatusNotFound)
		}
		return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	business.Name = data.Name
	business.ImageURL = data.ImageURL
	if categoriesData != nil {
		business.Categories = categoriesData
	}
	if transactionsData != nil {
		business.Transactions = transactionsData
	}
	business.Price = data.Price
	business.Location.Latitude = data.Latitude
	business.Location.Longitude = data.Longitude
	business.Location.Address1 = data.Address1
	business.Location.Address2 = data.Address2
	business.Location.Address3 = data.Address3
	business.Location.City = data.City
	business.Location.ZipCode = data.ZipCode
	business.Location.Country = data.Country
	business.Location.State = data.State
	business.Phone = data.Phone
	business.OpenAt = data.OpenAt
	business.CloseAt = data.CloseAt

	err = uc.businessRepo.Update(ctx, business)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (uc *businessUsecase) GetBusiness(ctx context.Context, id string) (dto.BusinessResponse, error) {
	var res dto.BusinessResponse
	business, err := uc.businessRepo.FindByIDOrAlias(ctx, id)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, helpers.NewError(helpers.ErrBusinessNotFound, http.StatusNotFound)
		}
		return res, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	businessModelToDto(ctx, &res, business, false, nil)

	return res, nil
}

func (uc *businessUsecase) Delete(ctx context.Context, id string) error {
	business, err := uc.businessRepo.FindByIDOrAlias(ctx, id)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.NewError(helpers.ErrBusinessNotFound, http.StatusNotFound)
		}
		return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	err = uc.businessRepo.Delete(ctx, business)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		return helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (uc *businessUsecase) GetBusinesses(ctx context.Context, params dto.BusinessQueryParams) (dto.BusinessesResponse, error) {
	var data dto.BusinessesResponse
	if err := params.Validate(); err != nil {
		return data, helpers.NewError(helpers.NewValdiationError(err), http.StatusBadRequest)

	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.PerPage <= 0 {
		params.PerPage = 5
	}

	var query domain.BusinessQuery
	now := time.Now()

	query.Location = strings.ToLower(params.Location)

	if params.OpenNow {
		query.OpenNow = now.Format("1504")
	}

	if params.OpenAt != "" {
		query.OpenAt = params.OpenAt
	}

	if params.Categories != "" {
		query.Categories = params.CategoriesList()
	}

	if params.TransactionType != "" {
		query.TransactionType = params.TransactionsList()
	}

	if params.PerPage > 0 {
		query.Limit = params.PerPage
	}

	if params.Page > 0 {
		query.Offset = (params.Page - 1) * params.PerPage
	}

	businesses, total, err := uc.businessRepo.FindByParams(ctx, query)
	if err != nil {
		uc.logger.Error("Business Usecase", "function", helpers.GetFunctionName(), "err", err)
		return data, helpers.NewError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	businessesData := make([]dto.BusinessResponse, len(businesses))
	for i, business := range businesses {
		b := dto.BusinessResponse{}
		businessModelToDto(ctx, &b, business, true, &now)
		businessesData[i] = b
	}

	data.Businesses = businessesData
	data.Metadata.Page = params.Page
	data.Metadata.PerPage = params.PerPage
	data.Metadata.Total = total
	data.Metadata.TotalPages = int(math.Ceil(float64(data.Metadata.Total) / float64(data.Metadata.PerPage)))
	data.Metadata.CurrentTime = now

	return data, nil
}

func businessModelToDto(ctx context.Context, dst *dto.BusinessResponse, src domain.Business, useIsOpen bool, currentTime *time.Time) {
	dst.ID = src.ID
	dst.Alias = src.Alias
	dst.Name = src.Name
	dst.ImageURL = src.ImageURL

	if useIsOpen {
		dst.IsOpen = helpers.IsOpen(*currentTime, src.OpenAt, src.CloseAt)
	}

	if !useIsOpen {
		dst.OpenAt = src.OpenAt
		dst.CloseAt = src.CloseAt
	}
	dst.URL = fmt.Sprintf("http://%s/api/v1/businesses/%s", ctx.Value("host"), src.Alias)
	dst.Categories = make([]dto.Category, len(src.Categories))
	for i, category := range src.Categories {
		dst.Categories[i] = dto.Category{
			Title: category.Title,
			Alias: category.Alias,
		}
	}

	dst.Coordinates = dto.Coordinates{
		Longitude: src.Location.Longitude,
		Latitude:  src.Location.Latitude,
	}

	dst.Transactions = make([]string, len(src.Transactions))
	for i, transaction := range src.Transactions {
		dst.Transactions[i] = transaction.Type
	}

	dst.Price = src.Price
	dst.Location = dto.Location{
		Address1: src.Location.Address1,
		Address2: src.Location.Address2,
		Address3: src.Location.Address3,
		City:     src.Location.City,
		ZipCode:  strconv.Itoa(src.Location.ZipCode),
		Country:  src.Location.Country,
		State:    src.Location.State,
	}
	for _, a := range []string{dst.Location.Address1, dst.Location.Address2, dst.Location.Address3} {
		if a != "" {
			dst.Location.DisplayAddress = append(dst.Location.DisplayAddress, a)
		}
	}
	dst.Location.DisplayAddress = append(dst.Location.DisplayAddress, fmt.Sprintf("%s, %s %s", dst.Location.City, dst.Location.State, dst.Location.ZipCode))
	dst.Phone = src.Phone
}
