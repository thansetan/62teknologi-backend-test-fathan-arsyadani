package main

import (
	_ "github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/docs"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/infrastructure/database"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/infrastructure/http"
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

	app := http.New(db, config.App, logger)

	if err := app.Start(); err != nil {
		panic(err)
	}

}
