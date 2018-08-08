package scene

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/BrendanBurkhart/raytracer/pkg/raytracing"
	"github.com/BrendanBurkhart/raytracer/pkg/raytracing/object"
)

// Scene describes a renderable scene and holds an output image
type Scene struct {
	Materials    []raytracing.Material `json:"materials"`
	Objects      []object.Object       `json:"objects"`
	Lights       []raytracing.Light    `json:"lights"`
	ambientLight raytracing.Color
}

// Initialize must be called before the Scene is used
func (s *Scene) Initialize() (e error) {
	for i, object := range s.Objects {
		materialID := object.MaterialID()
		if materialID < 0 || materialID >= len(s.Materials) {
			msg := fmt.Sprintf("invalid material id in object %d", i)
			e = errors.New(msg)
			return
		}
	}

	s.ambientLight = raytracing.Color{}
	for _, light := range s.Lights {
		s.ambientLight.Red += light.Ambient.Red
		s.ambientLight.Green += light.Ambient.Green
		s.ambientLight.Blue += light.Ambient.Blue
	}
	s.ambientLight.Red /= float64(len(s.Lights))
	s.ambientLight.Green /= float64(len(s.Lights))
	s.ambientLight.Blue /= float64(len(s.Lights))
	return
}

// UnmarshalJSON unmarshals a Scene containing a slice of object.Object interfaces
func (s *Scene) UnmarshalJSON(b []byte) error {
	type Alias Scene
	auxiliary := &struct {
		JSONObjects object.JSONObjects `json:"objects"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(b, &auxiliary); err != nil {
		return err
	}

	s.Objects = auxiliary.JSONObjects
	return nil
}

// FindIntersection finds the closest intersection between the specified ray and the scene.
// Returns whether an intersection was found, and if so where and with what object index.
func (s *Scene) FindIntersection(r raytracing.Ray) (bool, float64, int) {
	currentObject := -1
	t := 20000.0

	var intersected bool
	for i, object := range s.Objects {
		if intersected, t = object.Intersect(r, t); intersected {
			currentObject = i
		}
	}

	intersected = (currentObject != -1)
	return intersected, t, currentObject
}

// TraceRay traces a given ray to its first intersection and performs lighting calculations
func (s *Scene) TraceRay(r raytracing.Ray, lightStrength float64, remainingDepth int, lighting raytracing.LightingModel) (color raytracing.Color) {
	intersected, t, currentObject := s.FindIntersection(r)

	if !intersected {
		return
	}

	scaled := r.Direction.Scale(t)
	intersection := r.Position.Add(scaled)
	normal := s.Objects[currentObject].SurfaceNormal(intersection)
	material := s.Materials[s.Objects[currentObject].MaterialID()]

	viewer := r.Direction.Negative()
	var ok bool
	viewer, ok = viewer.Normalize()
	if !ok {
		return
	}

	visibleLights := []raytracing.Light{}
	for _, light := range s.Lights {
		lightRay := raytracing.Ray{
			Position:  intersection,
			Direction: light.Position.Subtract(intersection),
		}

		intersected, distance, _ := s.FindIntersection(lightRay)
		if !intersected {
			visibleLights = append(visibleLights, light)
		} else if distance > 1.0 {
			visibleLights = append(visibleLights, light)
		}
	}

	surfaceColor := lighting(visibleLights, s.ambientLight, viewer, intersection, normal, material)
	color.Red += surfaceColor.Red * lightStrength
	color.Green += surfaceColor.Green * lightStrength
	color.Blue += surfaceColor.Blue * lightStrength

	r.Position = intersection
	reflect := 2.0 * r.Direction.Dot(normal)
	r.Direction = r.Direction.Subtract(normal.Scale(reflect))
	r.Direction, _ = r.Direction.Normalize()

	var reflectedColor raytracing.Color
	if remainingDepth > 0 {
		reflectedColor = s.TraceRay(r, lightStrength*material.Reflectance, remainingDepth-1, lighting)
	}

	color.Red = color.Red + reflectedColor.Red
	color.Green = color.Green + reflectedColor.Green
	color.Blue = color.Blue + reflectedColor.Blue
	return
}
