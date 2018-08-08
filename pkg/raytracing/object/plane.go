package object

import (
	"encoding/json"
	"math"

	"github.com/BrendanBurkhart/raytracer/pkg/raytracing"
)

// Plane is an algebraic representation of a plane
type Plane struct {
	*Material
	Normal raytracing.Vector `json:"normal"`
	Point  raytracing.Vector `json:"point"`
}

func planeFactory(data *json.RawMessage) (Object, error) {
	obj := Plane{}
	if err := json.Unmarshal(*data, &obj); err != nil {
		return obj, err
	}
	obj.Normalize()
	return obj, nil
}

// Intersect returns whether there is an intersection with r within maxRange,
// and if so where it occurred. If there is no intersection, the scaling value will be maxRange
func (p Plane) Intersect(r raytracing.Ray, maxRange float64) (bool, float64) {
	denominator := r.Direction.Dot(p.Normal)

	if math.Abs(denominator) < 1e-8 {
		return false, maxRange
	}

	delta := p.Point.Subtract(r.Position)
	numerator := delta.Dot(p.Normal)

	t := numerator / denominator

	if t > 1e-4 && t < maxRange {
		return true, t
	}
	return false, maxRange
}

// SurfaceNormal returns the normal vector to the plane
func (p Plane) SurfaceNormal(point raytracing.Vector) raytracing.Vector {
	return p.Normal
}

// Normalize performs an in-place normalization of certain vectors normalized
// Position vectors, etc. are left un-normalized
func (p *Plane) Normalize() {
	normal, _ := p.Normal.Normalize()
	p.Normal = normal
}
