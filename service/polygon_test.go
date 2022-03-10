//go:build unit
// +build unit
package service

import (
	"app/mocks"
	"app/polygon"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPolygonService_CreateNewPolygon(t *testing.T) {
	var (
		repo           = new(mocks.PolygonRepository)
		polygonService = NewPolygonService(repo)
	)

	teardown := func() {
		repo.ExpectedCalls = nil
	}

	t.Run("Create and insert new polygon", func(t *testing.T) {
		var (
			ctx    = context.Background()
			vertex []polygon.Vertex
			poly   polygon.Polygon
			err    error
		)

		defer teardown()

		vertex = []polygon.Vertex{
			{
				X: 1,
				Y: 1,
			},
			{
				X: 2,
				Y: 1,
			},
			{
				X: 2,
				Y: 2,
			},
		}

		repo.On("InsertOne", ctx, mock.Anything).Return(nil)

		poly, err = polygonService.CreateNewPolygon(ctx, vertex)

		assert.NoError(t, err)
		assert.Equal(t, "Polygon_3", poly.Name[:len(poly.Name)-6])
		assert.Equal(t, 3, len(poly.Vertices))
		assert.Equal(t, 0.5, poly.Area)
		assert.Equal(t, []polygon.Vertex{
			{
				X: 2,
				Y: 1,
			},
			{
				X: 1,
				Y: 1,
			},
			{
				X: 2,
				Y: 2,
			}}, poly.Vertices)
	})

	t.Run("Return error if InsertOne returns error", func(t *testing.T) {
		var (
			ctx    = context.Background()
			vertex []polygon.Vertex
			err    error
		)

		defer teardown()

		vertex = []polygon.Vertex{
			{
				X: 1,
				Y: 1,
			},
			{
				X: 2,
				Y: 1,
			},
			{
				X: 2,
				Y: 2,
			},
		}

		repo.On("InsertOne", ctx, mock.Anything).Return(errors.New("error"))

		_, err = polygonService.CreateNewPolygon(ctx, vertex)

		assert.Error(t, err)
	})

}
