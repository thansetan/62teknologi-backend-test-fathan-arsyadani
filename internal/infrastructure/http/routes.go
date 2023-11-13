package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/routes"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/infrastructure/http/middlewares"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/utils/config"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	host   string
	port   int
	r      *mux.Router
	logger *slog.Logger
}

func New(db *gorm.DB, config config.App, logger *slog.Logger) *App {
	r := mux.NewRouter()
	return &App{db, config.Host, config.Port, r, logger}
}

func (app *App) Start() error {

	app.r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	app.r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	api := app.r.PathPrefix("/api/v1").Subrouter()
	api.Use(middlewares.LoggingMiddleware)
	{
		categories := api.PathPrefix("/categories").Subrouter()
		routes.NewCategoryRoutes(app.db, categories, app.logger)
	}

	{
		businesses := api.PathPrefix("/businesses").Subrouter()
		routes.NewBusinessRoutes(app.db, businesses, app.logger)
	}

	{
		transactions := api.PathPrefix("/transactions").Subrouter()
		routes.NewTransactionRoutes(app.db, transactions, app.logger)
	}

	app.logger.Info("Starting server", "host", app.host, "address", app.port)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", app.host, app.port), app.r)
}
