//go:build functional
package tests

import (
	"app/polygon"
	"context"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestPolygon = polygon.Polygon{
	Name: "Polygon_3_a2S4D",
	Vertices: []polygon.Vertex{
		{
			X: 0,
			Y: 0,
		},
		{
			X: 0,
			Y: 1,
		},
		{
			X: 1,
			Y: 1,
		},
		{
			X: 1,
			Y: 0,
		},
	},
	Area: 1,
}

var DB *mongo.Database

func CreateContainer() testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:5.0",
		ExposedPorts: []string{"27016:27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "admin",
			"MONGO_INITDB_ROOT_PASSWORD": "test",
		},
		Name: "mongodb",
	}

	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	return container
}

func InitDB() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:test@localhost:27016"))
	if err != nil {
		panic(err)
	}
	DB = client.Database("polygon")
}

func SeedPolygon(db *mongo.Database) {
	_, err := db.Collection("polygons").InsertOne(context.Background(), TestPolygon)
	if err != nil {
		panic(err)
	}
}
