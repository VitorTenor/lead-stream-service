package services

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type FileService struct {
	SchemaRepository repositories.SchemaRepository
	LeadRepository   repositories.LeadRepository
}

func NewFileService(sr repositories.SchemaRepository, lr repositories.LeadRepository) *FileService {
	return &FileService{
		SchemaRepository: sr,
		LeadRepository:   lr,
	}
}

func (fs *FileService) ProcessAndSave(ctx *context.Context, file *domain.File) error {
	schema, err := fs.SchemaRepository.FindById(ctx, file.SchemaId)
	if err != nil {
		return err
	}

	openedFile, err := file.File.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	reader := csv.NewReader(openedFile)
	headers, err := nextValue(reader)
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
	uniqueFieldsMap := make(map[string]map[string]bool)

	for _, field := range schema.Fields {
		if field.Unique {
			uniqueFieldsMap[field.Name] = make(map[string]bool)
		}
	}

	for {
		record, err := nextValue(reader)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		for i, value := range record {
			header := headers[i]
			if uniqueFields, ok := uniqueFieldsMap[header]; ok {
				if uniqueFields[value] {
					return domain.ErrDuplicatedValue
				}
				uniqueFields[value] = true
			}
		}

		doc, err := leadFromRecord(record, headers, *schema)
		if err != nil {
			return err
		}
		leads = append(leads, doc)
	}

	err = fs.LeadRepository.CreateMany(ctx, leads)
	if err != nil {
		return err
	}

	return nil
}

func leadFromRecord(record []string, headers []string, schema domain.Schema) (*bson.D, error) {
	doc := bson.D{}
	seen := make(map[string]string)
	for _, field := range schema.Fields {
		seen[field.Name] = field.Type
	}

	doc = append(doc, bson.E{Key: "schema_id", Value: schema.ID})

	dateTime := primitive.NewDateTimeFromTime(time.Now())

	for i, value := range record {
		parsedValue, err := domain.ValueFromType(value, seen[headers[i]])
		if err != nil {
			return nil, domain.ErrInvalidFieldValues
		}
		doc = append(doc, bson.E{Key: headers[i], Value: parsedValue})
	}

	doc = append(doc, bson.E{Key: "created_at", Value: dateTime})
	doc = append(doc, bson.E{Key: "updated_at", Value: dateTime})

	return &doc, nil
}

func nextValue(reader *csv.Reader) ([]string, error) {
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return nil, nil
			}
			if errors.Is(err, csv.ErrFieldCount) {
				return record, nil
			}
			return nil, err
		}
		if len(record) > 0 && strings.HasPrefix(record[0], "#") {
			continue
		}
		return record, nil
	}
}
