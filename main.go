package main

import (
	"app/config"
	"app/handler"
	"app/log"
	"app/repository"
	"app/service"
	"app/util"
	"context"

	_ "app/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Polygon API
// @version 1.0
// @description This is a API for Polygon.

// @contact.name Kenan Abbak
// @contact.url www.kenanabbak.com
// @contact.email kenanabbak@hotmail.com
func main() {
	config := config.NewConfigurations()
	// Connect database with credentials
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Mongo.ConnectionString))
	if err != nil {
		panic(err)
	}

	db := client.Database("polygon")

	logger := log.NewLogger()
	repository := repository.NewPolygonRepository(db)
	service := service.NewPolygonService(repository)
	handler := handler.NewPolygonHandler(service, repository, logger)
	e := echo.New()
	e.Validator = util.NewValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Add polygon prefix
	g := e.Group("/polygon")
	g.GET("/:name", handler.GetPolygonByName)
	g.POST("/", handler.CreateNewPolygonRequest)

	e.Logger.Fatal(e.Start(":" + config.Server.Port))
}
