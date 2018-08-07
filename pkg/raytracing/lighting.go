package raytracing

import (
	"math"
)

// Ray is a 3 dimensional ray
type Ray struct {
	Position  Vector `json:"position"`
	Direction Vector `json:"direction"`
}

// Color is a RGB color
type Color struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
}

// Material describes a syrface based on diffusion color and reflectance
type Material struct {
	Specular    Color   `json:"specular"`
	Diffuse     Color   `json:"diffuse"`
	Ambient     Color   `json:"ambient"`
	Alpha       float64 `json:"alpha"`
	Reflectance float64 `json:"reflectance"`
}

// Light describes a light source
type Light struct {
	Position Vector `json:"position"`
	Specular Color  `json:"specular"`
	Diffuse  Color  `json:"diffuse"`
	Ambient  Color  `json:"ambient"`
}

// LambertianReflectance calculates the Lambertian reflectance model. The surface normal should be normalized.
func LambertianReflectance(lights []Light, position Vector, normal Vector, material Material) (color Color) {
	for _, light := range lights {
		dist := light.Position.Subtract(position)

		// Light doesn't reach surface - angle between surface normal and light is more than 90
		if normal.Dot(dist) <= 0.0 {
			continue
		}

		// Normalize light ray
		lightVec, ok := dist.Normalize()
		if !ok {
			continue
		}

		// Lambert diffusion
		surfaceLightLevel := lightVec.Dot(normal)
		color.Red += surfaceLightLevel * light.Diffuse.Red * material.Diffuse.Red
		color.Green += surfaceLightLevel * light.Diffuse.Green * material.Diffuse.Green
		color.Blue += surfaceLightLevel * light.Diffuse.Blue * material.Diffuse.Blue
	}
	return
}

// PhongReflectance calculates the Phong reflectance model. The surface normal should be normalized.
func PhongReflectance(lights []Light, ambientLight Color, viewer Vector, position Vector, normal Vector, material Material) (color Color) {
	for _, light := range lights {
		dist := light.Position.Subtract(position)

		// Light doesn't reach surface - angle between surface normal and light is more than 90
		if normal.Dot(dist) <= 0.0 {
			continue
		}

		// Normalize light ray
		lightVec, ok := dist.Normalize()
		if !ok {
			continue
		}

		// Normalized reflected ray
		reflectDiff := normal.Scale(2.0 * lightVec.Dot(normal))
		reflectedLight := reflectDiff.Subtract(lightVec)

		reflectedLight, ok = reflectedLight.Normalize()
		if !ok {
			continue
		}

		diffCoef := math.Max(0.0, lightVec.Dot(normal))
		diffuse := Color{
			Red:   diffCoef * light.Diffuse.Red * material.Diffuse.Red,
			Green: diffCoef * light.Diffuse.Green * material.Diffuse.Green,
			Blue:  diffCoef * light.Diffuse.Blue * material.Diffuse.Blue}

		specBase := math.Max(0.0, reflectedLight.Dot(viewer))
		specCoef := math.Pow(specBase, material.Alpha)
		if specBase <= 0.0 {
			specCoef = 0.0
		}

		specular := Color{
			Red:   specCoef * light.Specular.Red * material.Specular.Red,
			Green: specCoef * light.Specular.Green * material.Specular.Green,
			Blue:  specCoef * light.Specular.Blue * material.Specular.Blue}

		color.Red += diffuse.Red + specular.Red
		color.Green += diffuse.Green + specular.Green
		color.Blue += diffuse.Blue + specular.Blue
	}

	color.Red += ambientLight.Red * material.Ambient.Red
	color.Green += ambientLight.Green * material.Ambient.Green
	color.Blue += ambientLight.Blue * material.Ambient.Blue

	return
}

// AverageColors returns the average of a slice of Color
func AverageColors(colors []Color) (average Color) {
	for _, color := range colors {
		average.Red += color.Red
		average.Green += color.Green
		average.Blue += color.Blue
	}
	average.Red /= float64(len(colors))
	average.Green /= float64(len(colors))
	average.Blue /= float64(len(colors))
	return
}
