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

func NewCategoryRoutes(db *gorm.DB, r *mux.Router, logger *slog.Logger) {
	repo := repository.NewCategoryRepository(db)
	uc := usecase.NewCategoryUsecase(repo, logger)
	handler := handler.NewCategoryHandler(uc)

	r.HandleFunc("", handler.GetAll).Methods(http.MethodGet)
}
