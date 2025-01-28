package services

import (
	"context"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
)

type SchemaService struct {
	SchemaRepository repositories.SchemaRepository
}

func NewSchemaService(sr repositories.SchemaRepository) *SchemaService {
	return &SchemaService{
		SchemaRepository: sr,
	}
}

func (s *SchemaService) ValidateAndSave(ctx *context.Context, schema *domain.Schema) (*domain.Schema, error) {
	schema.Normalize()

	if !schema.ValidateIfFieldsTypesAreValid() {
		return nil, domain.ErrInvalidFieldTypes
	}

	if !schema.ValidateIfFieldsAreUnique() {
		return nil, domain.ErrFieldsNotUnique
	}

	if !schema.ValidateCreatedAndUpdatedFields() {
		return nil, domain.ErrInvalidFieldValues
	}

	if !schema.ValidateIfRequiredFieldsArePresent() {
		return nil, domain.ErrRequiredFieldsNotPresent
	}

	err := s.SchemaRepository.Create(ctx, schema)
	if err != nil {
		return nil, err
	}

	schema, err = s.SchemaRepository.FindById(ctx, schema.ID.Hex())
	if err != nil {
		return nil, err
	}

	return schema, nil
}
