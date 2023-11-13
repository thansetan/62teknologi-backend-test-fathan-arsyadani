package routes

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/handler"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/repository"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/usecase"
	"gorm.io/gorm"
)

func NewBusinessRoutes(db *gorm.DB, r *mux.Router, logger *slog.Logger) {
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	businessRepo := repository.NewBusinessRepository(db)

	uc := usecase.NewBusinessUsecase(businessRepo, transactionRepo, categoryRepo, db, logger)

	handler := handler.NewBusinessHandler(uc)

	r.HandleFunc("", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/search", handler.GetBusinesses).Methods(http.MethodGet)
	r.HandleFunc("/{id_or_alias}", handler.Update).Methods(http.MethodPut)
	r.HandleFunc("/{id_or_alias}", handler.GetBusiness).Methods(http.MethodGet)
	r.HandleFunc("/{id_or_alias}", handler.Delete).Methods(http.MethodDelete)
}
