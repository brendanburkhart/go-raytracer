package object

import (
	"encoding/json"
	"math"

	"github.com/brendanburkhart/raytracer/pkg/raytracing"
)

// Sphere is a 3 dimensional sphere
type Sphere struct {
	*Material
	Radius float64           `json:"radius"`
	Center raytracing.Vector `json:"center"`
}

func sphereFactory(data *json.RawMessage) (Object, error) {
	obj := Sphere{}
	err := json.Unmarshal(*data, &obj)
	return obj, err
}

// Intersect returns whether there is an intersection with r within maxRange,
// and if so where it occurred. If there is no intersection, the scaling value will be maxRange
func (s Sphere) Intersect(r raytracing.Ray, maxRange float64) (bool, float64) {
	A := r.Direction.Dot(r.Direction)

	dist := r.Position.Subtract(s.Center)
	B := 2 * r.Direction.Dot(dist)

	C := dist.Dot(dist) - (s.Radius * s.Radius)

	discriminant := B*B - 4*A*C

	if discriminant < 0.0 {
		return false, maxRange
	}

	sqrtdiscr := math.Sqrt(discriminant)
	t0 := (-B + sqrtdiscr) / (2 * A)
	t1 := (-B - sqrtdiscr) / (2 * A)

	t := math.Min(t0, t1)

	if t > 1e-4 && t < maxRange {
		return true, t
	}
	return false, maxRange
}

// SurfaceNormal returns the normal vector to the sphere at the point specified
// by the position of the ray
func (s Sphere) SurfaceNormal(r raytracing.Ray) raytracing.Vector {
	normal, _ := r.Position.Subtract(s.Center).Normalize()
	return normal
}
