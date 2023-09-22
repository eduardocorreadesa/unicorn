package main

import (
	"context"
	"os"
	"unicorn/internal/configuration"
	"unicorn/internal/controller"
	"unicorn/internal/repository"
	"unicorn/internal/service"

	app "unicorn/internal/api"

	"log"
)

func main() {

	env := os.Getenv("ENV")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Started Unicorn API...")

	config := configuration.NewEnvConfig(ctx, env)

	unicornRepo := repository.UnicornData(config)
	unicornService := service.NewUnicornService(ctx, config, unicornRepo)
	unicornController := controller.NewUnicornController(ctx, config, unicornService)

	app := getApp(ctx, config, unicornController)
	app.SetupServer()
}

func getApp(ctx context.Context, config configuration.Config, unicornController controller.UnicornController) *app.API {
	return app.NewAPI(ctx, config, unicornController)
}
