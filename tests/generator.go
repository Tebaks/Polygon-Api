//go:build functional
// +build functional
package tests

import (
	p "app/polygon"
	"app/repository"
	"app/service"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func createRandomVertices(count int) []p.Vertex {
	vertices := make([]p.Vertex, count)
	for i := 0; i < count; i++ {
		vertices[i] = p.Vertex{
			X: float64(rand.Intn(100)),
			Y: float64(rand.Intn(100)),
		}
	}
	return vertices
}

func SeedRandomPolygons(count int) error {
	repository := repository.NewPolygonRepository(DB)
	service := service.NewPolygonService(repository)
	for i := 0; i < count; i++ {
		vertices := createRandomVertices(rand.Intn(8) + 3)
		_, err := service.CreateNewPolygon(context.Background(), vertices)
		if err != nil {
			return err
		}
	}

	polygonCount, err := DB.Collection("polygons").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	if polygonCount < int64(count) {
		return errors.New("Seeding polygons failed")
	}

	// Get random 3 polygons from DB
	var polygons []p.Polygon
	pipleline := []bson.D{bson.D{{"$sample", bson.D{{"size", 3}}}}}
	polygonsCursor, err := DB.Collection("polygons").Aggregate(context.Background(), pipleline)
	if err != nil {
		return err
	}

	err = polygonsCursor.All(context.Background(), &polygons)
	if err != nil {
		return err
	}

	// Check if all polygons are valid
	for _, polygon := range polygons {
		polyJSON, err := json.MarshalIndent(polygon, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(polyJSON))
	}

	return nil
}
