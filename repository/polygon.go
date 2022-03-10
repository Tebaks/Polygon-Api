package repository

import (
	"app/polygon"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrPolygonNotFound = errors.New("polygon not found with given name")
)

type PolygonRepository interface {
	FindByName(ctx context.Context, name string) (polygon.Polygon, error)
	InsertOne(ctx context.Context, polygon polygon.Polygon) error
}

type polygonRepository struct {
	collection *mongo.Collection
}

func NewPolygonRepository(db *mongo.Database) PolygonRepository {
	return &polygonRepository{
		collection: db.Collection("polygons"),
	}
}

func (r *polygonRepository) FindByName(ctx context.Context, name string) (polygon.Polygon, error) {
	var (
		polygon polygon.Polygon
		err     error
	)

	err = r.collection.FindOne(ctx, map[string]interface{}{"name": name}).Decode(&polygon)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return polygon, ErrPolygonNotFound
		}
		return polygon, err
	}

	return polygon, err
}

func (r *polygonRepository) InsertOne(ctx context.Context, polygon polygon.Polygon) error {
	var (
		err error
	)

	_, err = r.collection.InsertOne(ctx, polygon)
	return err
}
