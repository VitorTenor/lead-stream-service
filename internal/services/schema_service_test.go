package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSchemaService_ValidateAndSave(t *testing.T) {
	ctx := context.Background()
	service := NewSchemaService(NewSchemaRepositoryMock())

	_ = t.Run("success", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		phoneField := domain.SchemaField{
			Name:     "phone",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		emailField := domain.SchemaField{
			Name:     "email",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, emailField, phoneField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.NoError(t, err) {
			id, _ := primitive.ObjectIDFromHex("67696ff2e3f76ec9d8e8dc3b")
			_ = assert.Equal(t, id, schema.ID)
		}
	})

	_ = t.Run("invalid field type", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		phoneField := domain.SchemaField{
			Name:     "phone",
			Type:     "strAng",
			Required: true,
			Unique:   true,
		}

		emailField := domain.SchemaField{
			Name:     "email",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, emailField, phoneField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.Error(t, err) {
			_ = assert.Equal(t, domain.ErrInvalidFieldTypes, err)
		}
	})

	_ = t.Run("fields not unique", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		emailField := domain.SchemaField{
			Name:     "email",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, emailField, emailField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.Error(t, err) {
			_ = assert.Equal(t, domain.ErrFieldsNotUnique, err)
		}
	})

	_ = t.Run("required fields not present", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		nameField := domain.SchemaField{
			Name:     "name",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, nameField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.Error(t, err) {
			_ = assert.Equal(t, domain.ErrRequiredFieldsNotPresent, err)
		}
	})

	_ = t.Run("success, test if fields are normalized", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		phoneField := domain.SchemaField{
			Name:     "PhOnE",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		emailField := domain.SchemaField{
			Name:     "email",
			Type:     "string",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, emailField, phoneField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.NoError(t, err) {
			_ = assert.Equal(t, "email", schema.Fields[0].Name)
			_ = assert.Equal(t, "string", schema.Fields[0].Type)
			_ = assert.Equal(t, "phone", schema.Fields[1].Name)
			_ = assert.Equal(t, "string", schema.Fields[1].Type)
		}
	})

	_ = t.Run("success, test if types are case insensitive and normalized", func(t *testing.T) {
		// arrange
		var fields []domain.SchemaField

		phoneField := domain.SchemaField{
			Name:     "phone",
			Type:     "IntegeR",
			Required: true,
			Unique:   true,
		}

		emailField := domain.SchemaField{
			Name:     "email",
			Type:     "strIng",
			Required: true,
			Unique:   true,
		}

		fields = append(fields, emailField, phoneField)

		schema := &domain.Schema{Fields: fields}

		// act
		_, err := service.ValidateAndSave(&ctx, schema)

		// assert
		if assert.NoError(t, err) {
			_ = assert.Equal(t, "email", schema.Fields[0].Name)
			_ = assert.Equal(t, "string", schema.Fields[0].Type)
			_ = assert.Equal(t, "phone", schema.Fields[1].Name)
			_ = assert.Equal(t, "integer", schema.Fields[1].Type)
		}
	})

}
