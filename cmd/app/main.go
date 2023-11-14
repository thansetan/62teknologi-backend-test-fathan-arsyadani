package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/docs"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/infrastructure/database"
	app "github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/infrastructure/http"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/utils/config"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/utils/logger"
)

//	@title		62-Teknologi Back-end Test
//	@version	1.0
//
// @BasePath	/api/v1
func main() {
	config, err := config.Load(".env")
	if err != nil {
		panic(err)
	}

	logger := logger.New()

	db, err := database.New(config.DB)
	if err != nil {
		panic(err)
	}

	app := app.New(db, config.App, logger)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Starting App", "error", err.Error())
		}
	}()

	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting Down App")
	if err := app.Shutdown(ctx); err != nil {
		panic(err)
	}

	logger.Info("App Terminated")
}
