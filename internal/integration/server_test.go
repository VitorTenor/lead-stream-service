package integration

import (
	"context"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/docker/go-connections/nat"
	"github.com/labstack/echo/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vitortenor/lead-stream-service/internal/api"
	"github.com/vitortenor/lead-stream-service/internal/api/handlers"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"github.com/vitortenor/lead-stream-service/internal/services"
	"github.com/vitortenor/lead-stream-service/internal/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func InitServerTest() (*httptest.Server, error) {
	ctx := context.Background()

	mongoC, err := buildDockerContainer(ctx, "mongo:latest", "27017/tcp")
	if err != nil {
		return nil, err
	}

	host, err := mongoC.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	mongoURI := "mongodb://" + host + ":" + port.Port()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database("lead-stream-service-test")
	log.Println("Connected to in-memory database")

	if cleanDatabaseCollections(ctx, db, []string{"schemas", "leads"}) != nil {
		return nil, err
	}

	path, err := tools.FindProjectRoot()
	if err != nil {
		return nil, err
	}

	err = insertDataFromJSON(ctx, db, filepath.Join(path, "resources", "db", "test_data.json"))
	if err != nil {
		return nil, err
	}

	schemaHandler := handlers.NewSchemaHandler(
		services.NewSchemaService(
			repositories.NewSchemaRepository("schemas", db),
		),
	)

	fileHandler := handlers.NewFileHandler(
		services.NewFileService(
			repositories.NewSchemaRepository("schemas", db),
			repositories.NewLeadRepository("leads", db),
		),
	)

	e := echo.New()
	humaApi := humaecho.New(e, huma.DefaultConfig("api", "v1"))

	api.InitRoutes(humaApi, schemaHandler, fileHandler)

	ts := httptest.NewServer(e)

	return ts, nil
}

func buildDockerContainer(ctx context.Context, image string, port string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func cleanDatabaseCollections(ctx context.Context, db *mongo.Database, collNames []string) error {
	for _, collection := range collNames {
		_, err := db.Collection(collection).DeleteMany(ctx, bson.D{})
		if err != nil {
			return err
		}
	}
	return nil
}

func insertDataFromJSON(ctx context.Context, db *mongo.Database, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string][]map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	for collectionName, documents := range data {
		collection := db.Collection(collectionName)
		for _, document := range documents {
			if id, ok := document["_id"].(string); ok {
				objectId, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return err
				}
				document["_id"] = objectId
			}
			if _, err := collection.InsertOne(ctx, document); err != nil {
				return err
			}
		}
	}
	return nil
}
