package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/services"
	"net/http"
	"time"
)

func InitSchemaRoutes(humaApi huma.API, schemaHandler *SchemaHandler) {
	huma.Register(humaApi, huma.Operation{
		Path:          "/schema",
		OperationID:   "create-schema",
		Method:        http.MethodPost,
		DefaultStatus: http.StatusCreated,
		Summary:       "Create a new schema",
		Description:   "Create a new schema with the given fields",
	}, schemaHandler.Create)
}

type SchemaHandler struct {
	service *services.SchemaService
}

func NewSchemaHandler(service *services.SchemaService) *SchemaHandler {
	return &SchemaHandler{
		service: service,
	}
}

func (sh *SchemaHandler) Create(ctx context.Context, sr *SchemaRequest) (*SchemaResponse, error) {
	schema, err := sh.service.ValidateAndSave(&ctx, sr.toDomain())
	if err != nil {
		return nil, handleError(err)
	}

	return schemaToResponse(schema), nil
}

type SchemaRequest struct {
	Body struct {
		Fields []struct {
			Name     string `json:"name" required:"true" default:"name" description:"The name of the field"`
			Type     string `json:"type" required:"true" default:"string" description:"The type of the field"`
			Required bool   `json:"required,omitempty" optional:"true" default:"false" description:"Indicates if the field is required"`
			Unique   bool   `json:"unique,omitempty" optional:"true" default:"false" description:"Indicates if the field is unique"`
		} `json:"fields" required:"true" description:"The fields of the schema"`
	}
}

func (sr *SchemaRequest) toDomain() *domain.Schema {
	var fields []domain.SchemaField

	for _, f := range sr.Body.Fields {
		fields = append(fields, domain.SchemaField{
			Name:     f.Name,
			Type:     f.Type,
			Required: f.Required,
			Unique:   f.Unique,
		})
	}

	return &domain.Schema{
		Fields: fields,
	}
}

type SchemaResponse struct {
	Body struct {
		ID        string                 `json:"id" description:"The ID of the schema"`
		Fields    []SchemaResponseFields `json:"fields" description:"The fields of the schema"`
		CreatedAt string                 `json:"created_at" description:"The creation date of the schema"`
		UpdatedAt string                 `json:"updated_at" description:"The last update date of the schema"`
	}
}

type SchemaResponseFields struct {
	Name     string `json:"name" description:"The name of the field"`
	Type     string `json:"type" description:"The type of the field"`
	Required bool   `json:"required" description:"Indicates if the field is required"`
	Unique   bool   `json:"unique" description:"Indicates if the field is unique"`
}

func schemaToResponse(schema *domain.Schema) *SchemaResponse {
	var fields []SchemaResponseFields

	for _, f := range schema.Fields {
		fields = append(fields, SchemaResponseFields{
			Name:     f.Name,
			Type:     f.Type,
			Required: f.Required,
			Unique:   f.Unique,
		})
	}

	return &SchemaResponse{
		Body: struct {
			ID        string                 `json:"id" description:"The ID of the schema"`
			Fields    []SchemaResponseFields `json:"fields" description:"The fields of the schema"`
			CreatedAt string                 `json:"created_at" description:"The creation date of the schema"`
			UpdatedAt string                 `json:"updated_at" description:"The last update date of the schema"`
		}{
			ID:        schema.ID.Hex(),
			Fields:    fields,
			CreatedAt: schema.CreatedAt.Time().Format(time.DateTime),
			UpdatedAt: schema.UpdatedAt.Time().Format(time.DateTime),
		},
	}
}
