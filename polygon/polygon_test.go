//go:build unit
package polygon

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestPolygon_GenerateName(t *testing.T) {
	tests := []struct {
		prefix string
		p      *Polygon
	}{
		{
			prefix: "Polygon_3_",
			p: &Polygon{
				Vertices: []Vertex{
					{X: 1, Y: 2},
					{X: 3, Y: 4},
					{X: 5, Y: 6},
				},
			},
		},
		{
			prefix: "Polygon_4_",
			p: &Polygon{
				Vertices: []Vertex{
					{X: 1, Y: 2},
					{X: 3, Y: 4},
					{X: 5, Y: 6},
					{X: 7, Y: 8},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.prefix, func(t *testing.T) {
			tt.p.GenerateName()
			if !strings.HasPrefix(tt.p.Name, tt.prefix) {
				t.Errorf("Polygon.GenerateName() = %v, want prefix %v", tt.p.Name, tt.prefix)
			}
		})
	}
}

func Test_getDistance(t *testing.T) {
	type args struct {
		a Vertex
		b Vertex
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Find distance between two points",
			args: args{
				a: Vertex{
					X: 1,
					Y: 1,
				},
				b: Vertex{
					X: 2,
					Y: 2,
				},
			},
			want: math.Sqrt(2),
		},
		{
			name: "Find distance between two points",
			args: args{
				a: Vertex{
					X: 1,
					Y: 1,
				},
				b: Vertex{
					X: 1,
					Y: 2,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDistance(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("getDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolygon_findCenter(t *testing.T) {
	tests := []struct {
		name string
		p    *Polygon
		want Vertex
	}{
		{
			name: "Find center of a polygon",
			p: &Polygon{
				Vertices: []Vertex{
					{
						X: 0,
						Y: 1,
					},
					{
						X: 1,
						Y: 0,
					},
					{
						X: 0,
						Y: 0,
					},
					{
						X: 1,
						Y: 1,
					},
				},
			},
			want: Vertex{
				X: 0.5,
				Y: 0.5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.findCenter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Polygon.findCenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolygon_CalculateArea(t *testing.T) {
	tests := []struct {
		name string
		p    *Polygon
		want float64
	}{
		{
			name: "Calculate area of a polygon",
			p: &Polygon{
				Vertices: []Vertex{
					{
						X: 0,
						Y: 1,
					},
					{
						X: 1,
						Y: 0,
					},
					{
						X: 0,
						Y: 0,
					},
					{
						X: 1,
						Y: 1,
					},
				},
			},
			want: 1,
		},
		{
			name: "Calculate area of a polygon",
			p: &Polygon{
				Vertices: []Vertex{
					{
						X: 0,
						Y: 1,
					},
					{
						X: 3,
						Y: 0,
					},
					{
						X: 0,
						Y: 0,
					},
				},
			},
			want: 1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.CalculateArea()
			if tt.p.Area != tt.want {
				t.Errorf("Polygon.CalculateArea() = %v, want %v", tt.p.Area, tt.want)
			}
		})
	}
}

func TestPolygon_sortPoints(t *testing.T) {
	tests := []struct {
		name string
		p    *Polygon
		want []Vertex
	}{
		{
			name: "Sort points of a polygon",
			p: &Polygon{
				Vertices: []Vertex{
					{
						X: 0,
						Y: 1,
					},
					{
						X: 1,
						Y: 0,
					},
					{
						X: 0,
						Y: 0,
					},
					{
						X: 1,
						Y: 1,
					},
				},
			},
			want: []Vertex{
				{
					X: 1,
					Y: 0,
				},
				{
					X: 0,
					Y: 0,
				},
				{
					X: 0,
					Y: 1,
				},
				{
					X: 1,
					Y: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.sortPoints()
			if !reflect.DeepEqual(tt.p.Vertices, tt.want) {
				t.Errorf("Polygon.sortPoints() = %v, want %v", tt.p.Vertices, tt.want)
			}
		})
	}
}
