package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
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
		if !validTypes[strings.ToLower(field.Type)] {
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

var validTypes = map[string]bool{
	"string":   true,
	"integer":  true,
	"boolean":  true,
	"date":     true,
	"time":     true,
	"datetime": true,
}
