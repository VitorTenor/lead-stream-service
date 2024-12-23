package handlers

import (
	"context"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"github.com/vitortenor/lead-stream-service/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewSchemaServiceMock(sr repositories.SchemaRepository) *services.SchemaService {
	return &services.SchemaService{
		SchemaRepository: sr,
	}
}

func NewSchemaRepositoryMock() repositories.SchemaRepository {
	return &schemaRepositoryMock{}
}

type schemaRepositoryMock struct {
}

func (s schemaRepositoryMock) Create(_ *context.Context, schema *domain.Schema) error {
	schema.ID, _ = primitive.ObjectIDFromHex("67696ff2e3f76ec9d8e8dc3b")
	return nil
}

func (s schemaRepositoryMock) FindById(_ *context.Context, id string) (*domain.Schema, error) {
	if id == "67696ff2e3f76ec9d8e8dc3b" {
		val, _ := primitive.ObjectIDFromHex(id)
		return &domain.Schema{
			ID: val,
		}, nil
	}
	return nil, nil
}
