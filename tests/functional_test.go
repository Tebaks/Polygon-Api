//go:build functional
package tests

import (
	"app/handler"
	"app/log"
	"app/polygon"
	"app/repository"
	"app/service"
	"app/util"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	container := CreateContainer()
	defer container.Terminate(context.Background())
	InitDB()
	SeedPolygon(DB)
	err := SeedRandomPolygons(1000)
	if err != nil {
		panic(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}
func TestPolygon_CreatePolygon(t *testing.T) {
	var (
		ctx            = context.Background()
		repo           = repository.NewPolygonRepository(DB)
		service        = service.NewPolygonService(repo)
		logger         = log.NewLogger()
		polygonHandler = handler.NewPolygonHandler(service, repo, logger)
	)

	t.Run("CreatePolygon", func(t *testing.T) {
		var (
			res     handler.Response
			polygon polygon.Polygon
		)

		e := echo.New()
		req := httptest.NewRequest("POST", "/polygon/", strings.NewReader(`{"vertices":[{"x":1,"y":1},{"x":2,"y":1},{"x":2,"y":2},{"x":1,"y":1}]}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Echo().Validator = util.NewValidator()

		polygonHandler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)
		data := res.Data.(map[string]interface{})

		DB.Collection("polygons").FindOne(ctx, bson.M{"name": data["name"]}).Decode(&polygon)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, 4, len(data["vertices"].([]interface{})))
	})

	t.Run("CreatePolygon with invalid data", func(t *testing.T) {
		var (
			res handler.Response
		)

		e := echo.New()
		req := httptest.NewRequest("POST", "/polygon/", strings.NewReader(`{"vertices":[{"x":1,"y":1},{"x":2,"y":1}]}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Echo().Validator = util.NewValidator()

		polygonHandler.CreateNewPolygonRequest(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestPolygon_GetPolygonByName(t *testing.T) {
	var (
		repo           = repository.NewPolygonRepository(DB)
		service        = service.NewPolygonService(repo)
		logger         = log.NewLogger()
		polygonHandler = handler.NewPolygonHandler(service, repo, logger)
	)

	t.Run("Get polygon by name", func(t *testing.T) {
		var (
			res handler.Response
		)

		e := echo.New()
		req := httptest.NewRequest("POST", "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/polygon/:name")
		c.SetParamNames("name")
		c.SetParamValues(TestPolygon.Name)

		polygonHandler.GetPolygonByName(c)

		json.NewDecoder(rec.Body).Decode(&res)
		data := res.Data.(map[string]interface{})

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Polygon_3_a2S4D", data["name"])
		assert.Equal(t, 4, len(data["vertices"].([]interface{})))
		assert.Equal(t, 1.0, data["area"])
	})

	t.Run("Return polygon not found when no data", func(t *testing.T) {
		var (
			res handler.Response
		)

		e := echo.New()
		req := httptest.NewRequest("POST", "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/polygon/:name")
		c.SetParamNames("name")
		c.SetParamValues("Polygon_not_found")

		polygonHandler.GetPolygonByName(c)

		json.NewDecoder(rec.Body).Decode(&res)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "polygon not found with given name", res.Error)
	})
}
