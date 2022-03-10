//go:build unit
// +build unit
package repository

import (
	p "app/polygon"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestPolygonRepository_FindByName(t *testing.T) {
	var (
		mt = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer mt.Close()

	mt.Run("Success", func(mt *mtest.T) {
		var (
			err          error
			polygon, res p.Polygon
			repo         = NewPolygonRepository(mt.DB)
		)

		polygon = p.Polygon{
			Name: "Polygon_3",
			Vertices: []p.Vertex{
				{X: 1, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			Area: 1,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"name", polygon.Name},
			{"vertices", polygon.Vertices},
			{"area", polygon.Area},
		}))

		res, err = repo.FindByName(context.Background(), polygon.Name)

		assert.Nil(t, err)
		assert.Equal(t, polygon.Name, res.Name)
		assert.Equal(t, polygon.Vertices, res.Vertices)
		assert.Equal(t, polygon.Area, res.Area)
	})

	mt.Run("Return ErrPolygonNotFound when no data found", func(mt *mtest.T) {
		var (
			err  error
			res  p.Polygon
			repo = NewPolygonRepository(mt.DB)
		)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		res, err = repo.FindByName(context.Background(), "foo")

		assert.Equal(t, ErrPolygonNotFound, err)
		assert.Equal(t, p.Polygon{}, res)
	})

	mt.Run("Return error when find returns error", func(mt *mtest.T) {
		var (
			err  error
			res  p.Polygon
			repo = NewPolygonRepository(mt.DB)
		)

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "foo",
		}))

		res, err = repo.FindByName(context.Background(), "foo")

		assert.NotNil(t, err)
		assert.Equal(t, p.Polygon{}, res)
	})
}

func TestPolygonRepository_InsertOne(t *testing.T) {
	var (
		mt = mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	)
	defer mt.Close()

	mt.Run("Success", func(mt *mtest.T) {
		var (
			err     error
			polygon p.Polygon
			repo    = NewPolygonRepository(mt.DB)
		)

		polygon = p.Polygon{
			Name: "Polygon_3",
			Vertices: []p.Vertex{
				{X: 1, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			Area: 1,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err = repo.InsertOne(context.Background(), polygon)

		assert.Nil(t, err)
	})

	mt.Run("Return error when insert returns error", func(mt *mtest.T) {
		var (
			err     error
			polygon p.Polygon
			repo    = NewPolygonRepository(mt.DB)
		)

		polygon = p.Polygon{
			Name: "Polygon_3",
			Vertices: []p.Vertex{
				{X: 1, Y: 1},
				{X: 1, Y: 0},
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			Area: 1,
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		err = repo.InsertOne(context.Background(), polygon)

		assert.NotNil(t, err)
	})
}
