package services

import (
	"context"
	"encoding/csv"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

type FileService struct {
	schemaRepository repositories.SchemaRepository
	leadRepository   repositories.LeadRepository
}

func NewFileService(sr repositories.SchemaRepository, lr repositories.LeadRepository) *FileService {
	return &FileService{
		schemaRepository: sr,
		leadRepository:   lr,
	}
}

func (fs *FileService) ProcessAndSave(ctx *context.Context, file *domain.File) error {
	schema, err := fs.schemaRepository.FindById(ctx, file.SchemaId)
	if err != nil {
		return err
	}

	openedFile, err := file.File.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	reader := csv.NewReader(openedFile)
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	if !domain.ValidateRequiredFields(headers) {
		return domain.ErrRequiredFieldsMissing
	}

	if !domain.ValidateRequiredFieldsFromSchema(headers, schema.Fields) {
		return domain.ErrRequiredFieldsMissing
	}

	var leads []*bson.D
	seen := make(map[string]map[string]bool)

	for _, field := range schema.Fields {
		if field.Unique {
			seen[field.Name] = make(map[string]bool)
		}
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		for i, value := range record {
			header := headers[i]
			if uniqueFields, ok := seen[header]; ok {
				if uniqueFields[value] {
					return domain.ErrDuplicatedValue
				}
				uniqueFields[value] = true
			}
		}

		document, err := leadFromRecord(record, headers, *schema)
		if err != nil {
			return err
		}
		leads = append(leads, document)
	}

	err = fs.leadRepository.CreateMany(ctx, leads)
	if err != nil {
		return err
	}

	return nil
}

func leadFromRecord(record []string, headers []string, schema domain.Schema) (*bson.D, error) {
	document := bson.D{}
	seen := make(map[string]string)
	for _, field := range schema.Fields {
		seen[field.Name] = field.Type
	}

	document = append(document, bson.E{Key: "schema_id", Value: schema.ID})

	for i, value := range record {
		parsedValue, err := domain.ValueFromType(value, seen[headers[i]])
		if err != nil {
			return nil, domain.ErrInvalidFieldValues
		}
		document = append(document, bson.E{Key: headers[i], Value: parsedValue})
	}

	return &document, nil
}
