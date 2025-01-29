package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadRepository interface {
	CreateMany(ctx *context.Context, leads []*bson.D) error
	Create(ctx *context.Context, lead *bson.D) error
}

func NewLeadRepository(collName string, db *mongo.Database) LeadRepository {
	return &leadRepository{
		coll: db.Collection(collName),
	}
}

type leadRepository struct {
	coll *mongo.Collection
}

func (lr *leadRepository) Create(ctx *context.Context, lead *bson.D) error {
	_, err := lr.coll.InsertOne(*ctx, lead)

	return err
}

func (lr *leadRepository) CreateMany(ctx *context.Context, leads []*bson.D) error {
	doc := make([]interface{}, len(leads))
	for i, v := range leads {
		doc[i] = v
	}

	_, err := lr.coll.InsertMany(*ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
