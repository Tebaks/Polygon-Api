// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	polygon "app/polygon"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// PolygonService is an autogenerated mock type for the PolygonService type
type PolygonService struct {
	mock.Mock
}

// CreateNewPolygon provides a mock function with given fields: ctx, vertices
func (_m *PolygonService) CreateNewPolygon(ctx context.Context, vertices []polygon.Vertex) (polygon.Polygon, error) {
	ret := _m.Called(ctx, vertices)

	var r0 polygon.Polygon
	if rf, ok := ret.Get(0).(func(context.Context, []polygon.Vertex) polygon.Polygon); ok {
		r0 = rf(ctx, vertices)
	} else {
		r0 = ret.Get(0).(polygon.Polygon)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []polygon.Vertex) error); ok {
		r1 = rf(ctx, vertices)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
