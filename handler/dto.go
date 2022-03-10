package handler

import "app/polygon"

type CreateNewPolygonRequest struct {
	Vertices []polygon.Vertex `json:"vertices" validate:"required"`
}

type GetPolygonByNameRequest struct {
	Name string `json:"name" param:"name"`
}
