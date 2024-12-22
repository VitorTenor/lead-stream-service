package repositories

import (
	"context"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SchemaRepository interface {
	Create(ctx *context.Context, schema *domain.Schema) error
	FindById(ctx *context.Context, id string) (*domain.Schema, error)
}

func NewSchemaRepository(collName string, db *mongo.Database) SchemaRepository {
	return &schemaRepository{
		coll: db.Collection(collName),
	}
}

type schemaRepository struct {
	coll *mongo.Collection
}

func (r *schemaRepository) Create(ctx *context.Context, schema *domain.Schema) error {
	schema.ID = primitive.NewObjectID()
	schema.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	schema.UpdatedAt = schema.CreatedAt

	_, err := r.coll.InsertOne(*ctx, schema)
	if err != nil {
		return err
	}

	return nil
}

func (r *schemaRepository) FindById(ctx *context.Context, id string) (*domain.Schema, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var schema domain.Schema
	err = r.coll.FindOne(*ctx, primitive.M{"_id": objID}).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
