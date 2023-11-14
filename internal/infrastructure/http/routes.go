package http

import (
	"context"
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
	logger *slog.Logger
	server *http.Server
}

func New(db *gorm.DB, config config.App, logger *slog.Logger) *App {
	return &App{
		db:     db,
		host:   config.Host,
		port:   config.Port,
		logger: logger,
	}
}

func (app *App) Start() error {
	r := mux.NewRouter()
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
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

	app.server = &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", app.host, app.port),
	}

	return app.server.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
