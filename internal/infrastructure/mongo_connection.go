package infrastructure

import (
	"context"
	"log"

	"github.com/vitortenor/lead-stream-service/internal/configuration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateConnection(ctx context.Context, envConfig *configuration.Config) (*mongo.Database, error) {
	log.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(envConfig.Database.URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(envConfig.Database.Name)

	err = createIndex(ctx, db.Collection(envConfig.Database.Collection["leads"]))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createIndex(ctx context.Context, collection *mongo.Collection) error {
	var indexModel []mongo.IndexModel

	indexModel = append(indexModel, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	indexModel = append(indexModel, mongo.IndexModel{
		Keys:    bson.D{{Key: "telephone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	_, err := collection.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}
