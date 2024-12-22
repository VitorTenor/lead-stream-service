package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/vitortenor/lead-stream-service/internal/api/handlers"
)

func InitRoutes(humaApi huma.API, sh *handlers.SchemaHandler, fh *handlers.FileHandler) {
	handlers.InitSchemaRoutes(humaApi, sh)
	handlers.InitFileRoutes(humaApi, fh)
}
