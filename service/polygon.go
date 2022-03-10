package service

import (
	"app/polygon"
	"app/repository"
	"context"
)

type PolygonService interface {
	CreateNewPolygon(ctx context.Context, vertices []polygon.Vertex) (polygon.Polygon, error)
}

type polygonService struct {
	polygonRepository repository.PolygonRepository
}

func NewPolygonService(polygonRepository repository.PolygonRepository) PolygonService {
	return &polygonService{
		polygonRepository: polygonRepository,
	}
}

func (p *polygonService) CreateNewPolygon(ctx context.Context, vertices []polygon.Vertex) (polygon.Polygon, error) {
	var (
		polygon polygon.Polygon
		err     error
	)

	polygon.Vertices = vertices
	polygon.CalculateArea()
	polygon.GenerateName()

	if err = p.polygonRepository.InsertOne(ctx, polygon); err != nil {
		return polygon, err
	}

	return polygon, nil
}
