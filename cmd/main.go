package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/vitortenor/lead-stream-service/internal/api"
	"github.com/vitortenor/lead-stream-service/internal/api/handlers"
	"github.com/vitortenor/lead-stream-service/internal/configuration"
	"github.com/vitortenor/lead-stream-service/internal/infrastructure"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"github.com/vitortenor/lead-stream-service/internal/services"
	"github.com/vitortenor/lead-stream-service/internal/tools"
)

func main() {
	ctx := context.Background()

	configPath, err := tools.FindProjectRoot()
	if err != nil {
		log.Fatal("Failed to find project root: ", err)
	}
	envConfig, err := configuration.InitConfig(ctx, filepath.Join(configPath, "config.yaml"))
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}
	log.Println("Configuration loaded")

	db, err := infrastructure.CreateConnection(ctx, envConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	log.Println("Connected to database")

	defer func() {
		if err := db.Client().Disconnect(ctx); err != nil {
			log.Println("Failed to disconnect from database: ", err)
		}
	}()

	schemaHandler := handlers.NewSchemaHandler(
		services.NewSchemaService(
			repositories.NewSchemaRepository(envConfig.Database.Collection["schemas"], db),
		),
	)

	fileHandler := handlers.NewFileHandler(
		services.NewFileService(
			repositories.NewSchemaRepository(envConfig.Database.Collection["schemas"], db),
			repositories.NewLeadRepository(envConfig.Database.Collection["leads"], db),
		),
	)

	e := echo.New()
	humaApi := humaecho.New(e, huma.DefaultConfig(envConfig.Server.API.Name, envConfig.Server.API.Version))

	api.InitRoutes(humaApi, schemaHandler, fileHandler)

	address := fmt.Sprintf("%s:%d", envConfig.Server.Host, envConfig.Server.Port)
	log.Println("Server started on " + address)

	if http.ListenAndServe(address, e) != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
