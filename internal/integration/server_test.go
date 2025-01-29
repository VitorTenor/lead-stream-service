package integration

import (
	"context"
	"log"
	"net/http/httptest"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/docker/go-connections/nat"
	"github.com/labstack/echo/v4"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vitortenor/lead-stream-service/internal/api"
	"github.com/vitortenor/lead-stream-service/internal/api/handlers"
	"github.com/vitortenor/lead-stream-service/internal/domain"
	"github.com/vitortenor/lead-stream-service/internal/repositories"
	"github.com/vitortenor/lead-stream-service/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitServerTest() (*httptest.Server, error) {
	ctx := context.Background()

	mongoC, err := buildDockerContainer(&ctx, "mongo:latest", "27017/tcp")
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

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("lead-stream-service-test")
	log.Println("Connected to in-memory database")

	err = seed(&ctx, db.Collection("schemas"), createSchema())
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

func buildDockerContainer(ctx *context.Context, image string, port string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)),
	}
	return testcontainers.GenericContainer(*ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func seed(ctx *context.Context, coll *mongo.Collection, data []domain.Schema) error {
	for _, document := range data {
		_, err := coll.InsertOne(*ctx, document)
		if err != nil {
			return err
		}
	}
	return nil
}

func createSchema() []domain.Schema {
	createField := func(name, fieldType string, required, unique bool) domain.SchemaField {
		return domain.SchemaField{
			Name:     name,
			Type:     fieldType,
			Required: required,
			Unique:   unique,
		}
	}

	schemaFields := []domain.SchemaField{
		createField("email", "string", true, true),
		createField("name", "string", true, false),
		createField("phone", "integer", true, true),
		createField("lastname", "string", false, false),
	}

	objectId, _ := primitive.ObjectIDFromHex("67808a19c567c857d77d7f12")

	schema := domain.Schema{
		ID:        objectId,
		Fields:    schemaFields,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	return []domain.Schema{schema}
}
