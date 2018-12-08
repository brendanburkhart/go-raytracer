package camera

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/brendanburkhart/raytracer/pkg/raytracing"
)

// Lens calculates light rays from the camera into the scene
type Lens interface {
	generateLightRay(screenX float64, screenY float64, scope Scope) raytracing.Ray
	setAspectRatio(ratio float64) error
	GetLensName() string
}

// namedLens provides an embeddable label for lenses
type namedLens struct {
	name string
}

// GetName returns a string with the name of the lens
func (nl namedLens) GetLensName() string {
	return nl.name
}

// ViewPort defines a view port based on a view width and an aspect ratio
type ViewPort struct {
	ViewWidth  float64 `json:"viewWidth"`
	viewHeight float64
}

// setAspectRatio sets the view port height to the specified aspect ratio
// One of setAspectRatio or setBounds should be called to define a usable view port before use
func (v *ViewPort) setAspectRatio(ratio float64) error {
	v.viewHeight = v.ViewWidth / ratio
	return nil
}

// OrthographicLens provides light ray generation for orthographic rendering
type OrthographicLens struct {
	*ViewPort
	*namedLens
}

// generateLightRay creates a light ray from the lens passing through the point represented by (screenX, screenY)
// screenX and screenY range from -1.0 in the lower left corner to 1.0 in the upper right
func (l *OrthographicLens) generateLightRay(screenX float64, screenY float64, scope Scope) raytracing.Ray {
	lightRay := raytracing.Ray{}

	horizontal := scope.GetRight().Scale(screenX * l.ViewWidth * 0.5)
	vertical := scope.GetUp().Scale(screenY * l.viewHeight * 0.5)

	lightRay.Position = scope.Position.Add(horizontal).Add(vertical)
	lightRay.Direction = scope.GetForward()
	return lightRay
}

// FisheyeLens provides light ray generation for fisheye rendering
type FisheyeLens struct {
	HFOV float64 `json:"hfov"`
	VFOV float64 `json:"vfov"`
	*namedLens
}

// setAspectRatio sets the view port height to the specified aspect ratio
func (l *FisheyeLens) setAspectRatio(ratio float64) error {
	l.VFOV = l.HFOV / ratio
	return nil
}

// generateLightRay creates a light ray from the lens passing through the point represented by (screenX, screenY)
// screenX and screenY range from -1.0 in the lower left corner to 1.0 in the upper right
func (l *FisheyeLens) generateLightRay(screenX float64, screenY float64, scope Scope) raytracing.Ray {
	lightRay := raytracing.Ray{}

	horizontalAngle := -screenX * l.HFOV / 2.0
	verticalAngle := screenY * l.VFOV / 2.0

	direction := scope.GetForward()
	direction, _ = direction.Rotate(verticalAngle, scope.GetRight())
	direction, _ = direction.Rotate(horizontalAngle, scope.GetUp())
	direction, _ = direction.Normalize()

	lightRay.Position = scope.Position
	lightRay.Direction = direction
	return lightRay
}

// PerspectiveLens provides light ray generation for perspective rendering
type PerspectiveLens struct {
	OpticalRadius *float64 `json:"opticalRadius"`
	FocalLength   *float64 `json:"focalLength"`
	HFOV          float64  `json:"hfov"`
	ViewWidth     float64  `json:"viewWidth"`
	viewHeight    float64
	*namedLens
}

// setAspectRatio sets the view port height to the specified aspect ratio
func (l *PerspectiveLens) setAspectRatio(ratio float64) error {
	if l.HFOV != 0.0 {
		hfovRadian := l.HFOV / 180.0 * math.Pi

		if l.FocalLength != nil {
			opticalRadius := (1.0 / math.Cos(hfovRadian*0.5)) * *l.FocalLength
			l.OpticalRadius = &opticalRadius
		} else if l.OpticalRadius != nil {
			focalLength := math.Cos(hfovRadian*0.5) * *l.OpticalRadius
			l.FocalLength = &focalLength
		} else {
			return fmt.Errorf("with persepctive lens, at least one of focalLength or opticalRadius must be specified")
		}

		l.ViewWidth = math.Sin(hfovRadian*0.5) * *l.OpticalRadius * 2.0
	} else if l.FocalLength == nil {
		return fmt.Errorf("when using perspective lens with viewWidth, focalLength must be specified")
	}

	l.viewHeight = l.ViewWidth / ratio
	return nil
}

// generateLightRay creates a light ray from the lens passing through the point represented by (screenX, screenY)
// screenX and screenY range from -1.0 in the lower left corner to 1.0 in the upper right
func (l *PerspectiveLens) generateLightRay(screenX float64, screenY float64, scope Scope) raytracing.Ray {
	lightRay := raytracing.Ray{}

	direction := scope.GetForward().Scale(*l.FocalLength)
	direction = direction.Add(scope.GetRight().Scale(screenX * l.ViewWidth * 0.5))
	direction = direction.Add(scope.GetUp().Scale(screenY * l.viewHeight * 0.5))
	direction, _ = direction.Normalize()

	lightRay.Position = scope.Position
	lightRay.Direction = direction
	return lightRay
}

// CreateLens takes JSON data and returns an implementation of Lens matching that data
func CreateLens(b []byte) (Lens, error) {
	lens := &struct {
		Type string `json:"projection"`
	}{}

	if err := json.Unmarshal(b, &lens); err != nil {
		return nil, err
	}

	switch lens.Type {
	case "fisheye":
		var lens FisheyeLens
		if err := json.Unmarshal(b, &lens); err != nil {
			return nil, err
		}
		lens.namedLens = &namedLens{name: "fisheye"}
		return &lens, nil
	case "perspective":
		var lens PerspectiveLens
		if err := json.Unmarshal(b, &lens); err != nil {
			return nil, err
		}
		lens.namedLens = &namedLens{name: "perspective"}
		return &lens, nil
	default:
		var lens OrthographicLens
		if err := json.Unmarshal(b, &lens); err != nil {
			return nil, err
		}

		lens.namedLens = &namedLens{name: "orthographic"}
		return &lens, nil
	}
}
