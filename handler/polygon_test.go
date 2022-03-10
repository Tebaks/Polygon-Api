//go:build unit
package handler

import (
	"app/log"
	"app/mocks"
	"app/polygon"
	"app/repository"
	"app/util"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPolygonHandler_GetPolygonByName(t *testing.T) {
	var (
		repo           = new(mocks.PolygonRepository)
		polygonService = new(mocks.PolygonService)
		logger         = new(mocks.Logger)
		handler        = NewPolygonHandler(polygonService, repo, logger)
		echo           = echo.New()
	)

	teardown := func() {
		repo.ExpectedCalls = nil
		polygonService.ExpectedCalls = nil
	}

	t.Run("Success", func(t *testing.T) {
		var (
			poly polygon.Polygon
			res  Response
		)

		defer teardown()

		poly.Name = "Polygon_3"
		poly.Vertices = []polygon.Vertex{
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
			},
		}
		poly.Area = 0.5

		repo.On("FindByName", mock.Anything, "Polygon_3").Return(poly, nil)

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/:name")
		c.SetParamNames("name")
		c.SetParamValues("Polygon_3")

		handler.GetPolygonByName(c)

		json.NewDecoder(rec.Body).Decode(&res)
		data := res.Data.(map[string]interface{})

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, true, res.Success)
		assert.Equal(t, "Polygon_3", data["name"])
		assert.Equal(t, 3, len(data["vertices"].([]interface{})))
		assert.Equal(t, 0.5, data["area"])
	})

	t.Run("Return error when no data", func(t *testing.T) {
		var (
			res Response
		)

		defer teardown()

		repo.On("FindByName", mock.Anything, "Polygon_3").Return(polygon.Polygon{}, repository.ErrPolygonNotFound)
		logger.On("ErrorWithFields", "polygon not found", repository.ErrPolygonNotFound, log.Fields{"name": "Polygon_3"})

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/:name")
		c.SetParamNames("name")
		c.SetParamValues("Polygon_3")

		handler.GetPolygonByName(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, 404, rec.Code)
		assert.Equal(t, false, res.Success)
		assert.Equal(t, "polygon not found with given name", res.Error)
	})

	t.Run("Return error when polygonRepository.FindByName returns error", func(t *testing.T) {
		var (
			res Response
		)

		defer teardown()

		repo.On("FindByName", mock.Anything, "Polygon_3").Return(polygon.Polygon{}, errors.New("error"))
		logger.On("ErrorWithFields", "failed to get polygon by name", errors.New("error"), log.Fields{"name": "Polygon_3"})

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/:name")
		c.SetParamNames("name")
		c.SetParamValues("Polygon_3")

		handler.GetPolygonByName(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, 500, rec.Code)
		assert.Equal(t, false, res.Success)
		assert.Equal(t, "error", res.Error)
	})
}

func TestPolygonhandler_CreateNewPolygonRequest(t *testing.T) {
	var (
		repo           = new(mocks.PolygonRepository)
		polygonService = new(mocks.PolygonService)
		logger         = new(mocks.Logger)
		handler        = NewPolygonHandler(polygonService, repo, logger)
		echo           = echo.New()
	)

	teardown := func() {
		repo.ExpectedCalls = nil
		polygonService.ExpectedCalls = nil
	}

	t.Run("Success", func(t *testing.T) {
		var (
			res  Response
			poly polygon.Polygon
		)

		defer teardown()

		poly.Name = "Polygon_3"
		poly.Vertices = []polygon.Vertex{
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
			},
		}
		poly.Area = 0.5

		polygonService.On("CreateNewPolygon", mock.Anything, poly.Vertices).Return(poly, nil)

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"vertices":[{"x":2,"y":1},{"x":1,"y":1},{"x":2,"y":2}]}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/")
		c.Echo().Validator = util.NewValidator()

		handler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)
		data := res.Data.(map[string]interface{})

		assert.Equal(t, 201, rec.Code)
		assert.Equal(t, true, res.Success)
		assert.Equal(t, "Polygon_3", data["name"])
		assert.Equal(t, 3, len(data["vertices"].([]interface{})))
		assert.Equal(t, 0.5, data["area"])
	})

	t.Run("Return error when invalid body", func(t *testing.T) {
		var (
			res Response
		)

		defer teardown()

		logger.On("Error", "invalid request body", mock.Anything)

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"dots":"asd"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/")
		c.Echo().Validator = util.NewValidator()

		handler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.NotNil(t, res.Error)
		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, false, res.Success)
	})

	t.Run("Return error when polygonService.CreateNewPolygon returns error", func(t *testing.T) {
		var (
			res Response
		)

		defer teardown()

		vertices := []polygon.Vertex{
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
			},
		}

		polygonService.On("CreateNewPolygon", mock.Anything, mock.Anything).Return(polygon.Polygon{}, errors.New("error"))
		logger.On("ErrorWithFields", "failed to create new polygon", errors.New("error"), log.Fields{"vertices": vertices})

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"vertices":[{"x":2,"y":1},{"x":1,"y":1},{"x":2,"y":2}]}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/")
		c.Echo().Validator = util.NewValidator()

		handler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, 500, rec.Code)
		assert.Equal(t, false, res.Success)
		assert.Equal(t, "failed to create new polygon", res.Error)
	})

	t.Run("Return error when less than 3 vertices", func(t *testing.T) {
		var (
			res Response
		)

		defer teardown()

		logger.On("Error", "invalid request body", mock.Anything)

		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"vertices":[{"x":2,"y":1},{"x":1,"y":1}]}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := echo.NewContext(req, rec)
		c.SetPath("/polygon/")
		c.Echo().Validator = util.NewValidator()

		handler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, false, res.Success)
		assert.Equal(t, "polygon must have at least 3 vertices", res.Error)
	})
}
