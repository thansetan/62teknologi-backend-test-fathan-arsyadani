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

func NewTransactionRoutes(db *gorm.DB, r *mux.Router, logger *slog.Logger) {
	repo := repository.NewTransactionRepository(db)
	uc := usecase.NewTransactionUsecase(repo, logger)
	handler := handler.NewTransactionHandler(uc)

	r.HandleFunc("", handler.GetAll).Methods(http.MethodGet)
}
