package services

import (
	"context"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		return &domain.Schema{
			ID: primitive.NewObjectID(),
		}, nil
	}
	return nil, nil
}
