package services

import (
	"context"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
)

type SchemaService struct {
	schemaRepository repositories.SchemaRepository
}

func NewSchemaService(sr repositories.SchemaRepository) *SchemaService {
	return &SchemaService{
		schemaRepository: sr,
	}
}

func (s *SchemaService) ValidateAndSave(ctx *context.Context, schema *domain.Schema) (*domain.Schema, error) {
	if !schema.ValidateIfFieldsTypesAreValid() {
		return nil, domain.ErrInvalidFieldTypes
	}

	if !schema.ValidateIfFieldsAreUnique() {
		return nil, domain.ErrFieldsNotUnique
	}

	schema.Normalize()

	err := s.schemaRepository.Create(ctx, schema)
	if err != nil {
		return nil, err
	}

	schema, err = s.schemaRepository.FindById(ctx, schema.ID.Hex())
	if err != nil {
		return nil, err
	}

	return schema, nil
}
