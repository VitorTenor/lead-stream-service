package domain

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	ID        primitive.ObjectID `bson:"_id"`
	Fields    []SchemaField      `bson:"fields"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

type SchemaField struct {
	Name     string `bson:"name"`
	Type     string `bson:"type"`
	Required bool   `bson:"required"`
	Unique   bool   `bson:"unique"`
}

func (s *Schema) ValidateIfFieldsTypesAreValid() bool {
	for _, field := range s.Fields {
		if !validTypes[field.Type] {
			return false
		}
	}
	return true
}

func (s *Schema) ValidateIfFieldsAreUnique() bool {
	seen := make(map[string]bool)
	for _, field := range s.Fields {
		if seen[field.Name] {
			return false
		}
		seen[field.Name] = true
	}
	return true
}

func (s *Schema) Normalize() {
	for i := range s.Fields {
		s.Fields[i].Name = strings.ToLower(s.Fields[i].Name)
		s.Fields[i].Type = strings.ToLower(s.Fields[i].Type)
	}
}

func (s *Schema) ValidateIfRequiredFieldsArePresent() bool {
	seen := make(map[string]bool)
	for _, field := range s.Fields {
		seen[field.Name] = true
	}

	for requiredField := range requiredFields {
		if !seen[requiredField] {
			return false
		}
	}

	return true
}

func (s *Schema) ValidateCreatedAndUpdatedFields() bool {
	for _, field := range s.Fields {
		if field.Name == "created_at" || field.Name == "updated_at" {
			return false
		}
	}
	return true
}

var requiredFields = map[string]bool{
	"phone": true,
	"email": true,
}

var validTypes = map[string]bool{
	"string":   true,
	"integer":  true,
	"boolean":  true,
	"date":     true,
	"time":     true,
	"datetime": true,
}
