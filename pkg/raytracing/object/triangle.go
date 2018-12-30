package object

import (
	"encoding/json"

	"github.com/brendanburkhart/raytracer/pkg/raytracing"
)

// Triangle is a triangle in 3 dimensions
type Triangle struct {
	*Material
	normal raytracing.Vector
	edge1  raytracing.Vector
	edge2  raytracing.Vector
	A      raytracing.Vector `json:"A"`
	B      raytracing.Vector `json:"B"`
	C      raytracing.Vector `json:"C"`
}

func triangleFactory(data *json.RawMessage) (Object, error) {
	obj := Triangle{}
	if err := json.Unmarshal(*data, &obj); err != nil {
		return obj, err
	}
	obj.edge1 = obj.B.Subtract(obj.A)
	obj.edge2 = obj.C.Subtract(obj.A)

	obj.normal = obj.edge1.Cross(obj.edge2)
	obj.Normalize()
	return obj, nil
}

// Intersect returns whether there is an intersection with r within maxRange,
// and if so where it occurred. If there is no intersection, the scaling value will be maxRange
func (tr Triangle) Intersect(r raytracing.Ray, maxRange float64) (bool, float64) {
	h := r.Direction.Cross(tr.edge2)

	det := tr.edge1.Dot(h)
	if det < 1e-8 && det > -1e-8 {
		return false, maxRange
	}

	f := 1.0 / det

	transform := r.Position.Subtract(tr.A)

	u := transform.Dot(h) * f
	if u < 0.0 || u > 1.0 {
		return false, maxRange
	}

	q := transform.Cross(tr.edge1)

	v := r.Direction.Dot(q) * f
	if v < 0.0 || (u+v) > 1.0 {
		return false, maxRange
	}

	t := tr.edge2.Dot(q) * f
	if t > 1e-4 && t < maxRange {
		return true, t
	}
	return false, maxRange
}

// SurfaceNormal returns the normal vector to the triangle
func (tr Triangle) SurfaceNormal(r raytracing.Ray) raytracing.Vector {
	if r.Direction.Dot(tr.normal) < 0.0 {
		return tr.normal
	}
	return tr.normal.Negative()
}

// Normalize performs an in-place normalization of certain vectors normalized
// Position vectors, etc. are left un-normalized
func (tr *Triangle) Normalize() {
	normal, _ := tr.normal.Normalize()
	tr.normal = normal
}
