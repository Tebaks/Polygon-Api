package polygon

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Polygon struct {
	Name     string   `json:"name" example:"Polygon_3_a2cd24"`
	Vertices []Vertex `json:"vertices"`
	Area     float64  `json:"area" example:"12.5"`
}

type Vertex struct {
	X float64 `json:"x" example:"1"`
	Y float64 `json:"y" example:"2"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func (p *Polygon) GenerateName() {
	p.Name = fmt.Sprintf("Polygon_%d_%s", len(p.Vertices), randomString(5))
}

func (p *Polygon) CalculateArea() {
	p.sortPoints()

	lenght := len(p.Vertices)

	sum := 0.0
	for i := 0; i < lenght-1; i++ {
		sum += (p.Vertices[i].X * p.Vertices[i+1].Y) - (p.Vertices[i].Y * p.Vertices[i+1].X)
	}

	sum += (p.Vertices[lenght-1].X * p.Vertices[0].Y) - (p.Vertices[lenght-1].Y * p.Vertices[0].X)

	area := sum / 2

	p.Area = math.Abs(area)
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (p *Polygon) sortPoints() {
	center := p.findCenter()

	for i := range p.Vertices {
		p.Vertices[i].X -= center.X
		p.Vertices[i].Y -= center.Y
	}

	sort.Slice(p.Vertices, func(i, j int) bool {
		return comparePoints(p.Vertices[i], p.Vertices[j])
	})

	for i := range p.Vertices {
		p.Vertices[i].X += center.X
		p.Vertices[i].Y += center.Y
	}
}

func comparePoints(a, b Vertex) bool {
	angleA := getAngle(a, Vertex{X: 0, Y: 0})
	angleB := getAngle(b, Vertex{X: 0, Y: 0})

	if angleA < angleB {
		return false
	}

	distanceA := getDistance(a, Vertex{X: 0, Y: 0})
	distanceB := getDistance(b, Vertex{X: 0, Y: 0})

	if angleA == angleB && distanceA < distanceB {
		return false
	}

	return true
}

func getAngle(a, b Vertex) float64 {
	x := a.X - b.X
	y := a.Y - b.Y

	angle := math.Atan2(y, x)

	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}

func getDistance(a, b Vertex) float64 {
	x := a.X - b.X
	y := a.Y - b.Y

	return math.Sqrt(x*x + y*y)
}

func (p *Polygon) findCenter() Vertex {
	center := Vertex{X: 0, Y: 0}

	for _, v := range p.Vertices {
		center.X += v.X
		center.Y += v.Y
	}

	center.X /= float64(len(p.Vertices))
	center.Y /= float64(len(p.Vertices))

	return center
}
