package object

import (
	"encoding/json"
	"math"

	"github.com/brendanburkhart/raytracer/pkg/raytracing"
)

// Box is a representation of an axis aligned box
type Box struct {
	*Material
	MinCorner raytracing.Vector `json:"minCorner"`
	MaxCorner raytracing.Vector `json:"maxCorner"`
	center    raytracing.Vector
	extent    raytracing.Vector
}

func boxFactory(data *json.RawMessage) (Object, error) {
	obj := Box{}
	if err := json.Unmarshal(*data, &obj); err != nil {
		return obj, err
	}
	obj.Initialize()
	return obj, nil
}

// Initialize performs precomputation and preprocessing
func (b *Box) Initialize() {
	if b.MinCorner.X > b.MaxCorner.X {
		b.MinCorner.X, b.MaxCorner.X = b.MaxCorner.X, b.MinCorner.X
	}
	if b.MinCorner.Y > b.MaxCorner.Y {
		b.MinCorner.Y, b.MaxCorner.Y = b.MaxCorner.Y, b.MinCorner.Y
	}
	if b.MinCorner.Z > b.MaxCorner.Z {
		b.MinCorner.Z, b.MaxCorner.Z = b.MaxCorner.Z, b.MinCorner.Z
	}

	b.center = b.MaxCorner.Add(b.MinCorner).Scale(0.5)
	b.extent = b.MaxCorner.Subtract(b.center)
}

// Intersect returns whether there is an intersection with r within maxRange,
// and if so where it occurred. If there is no intersection, the scaling value will be maxRange
func (b Box) Intersect(r raytracing.Ray, maxRange float64) (bool, float64) {
	var tMin, tMax float64

	x1 := (b.MinCorner.X - r.Position.X) / r.Direction.X
	x2 := (b.MaxCorner.X - r.Position.X) / r.Direction.X

	tMin = math.Min(x1, x2)
	tMax = math.Max(x1, x2)

	y1 := (b.MinCorner.Y - r.Position.Y) / r.Direction.Y
	y2 := (b.MaxCorner.Y - r.Position.Y) / r.Direction.Y

	tMin = math.Max(tMin, math.Min(y1, y2))
	tMax = math.Min(tMax, math.Max(y1, y2))

	z1 := (b.MinCorner.Z - r.Position.Z) / r.Direction.Z
	z2 := (b.MaxCorner.Z - r.Position.Z) / r.Direction.Z

	tMin = math.Max(tMin, math.Min(z1, z2))
	tMax = math.Min(tMax, math.Max(z1, z2))

	if tMin < tMax && tMin > 1e-4 && tMin < maxRange {
		return true, tMin
	}
	return false, maxRange
}

// SurfaceNormal returns the normal vector to the box
func (b Box) SurfaceNormal(point raytracing.Vector) (normal raytracing.Vector) {
	relativePoint := point.Subtract(b.center)

	minDistance := math.Abs(math.Abs(relativePoint.X) - b.extent.X)
	normal = raytracing.Vector{X: signum(relativePoint.X), Y: 0, Z: 0}

	distance := math.Abs(math.Abs(relativePoint.Y) - b.extent.Y)
	if distance < minDistance {
		minDistance = distance
		normal = raytracing.Vector{X: 0, Y: signum(relativePoint.Y), Z: 0}
	}

	distance = math.Abs(math.Abs(relativePoint.Z) - b.extent.Z)
	if distance < minDistance {
		normal = raytracing.Vector{X: 0, Y: 0, Z: signum(relativePoint.Z)}
	}
	return
}

func signum(f float64) float64 {
	if f < 0.0 {
		return -1.0
	}
	return 1.0
}
