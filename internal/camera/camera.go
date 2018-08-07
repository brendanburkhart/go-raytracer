package camera

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"sync"

	"github.com/BrendanBurkhart/raytracer/internal/scene"

	"github.com/BrendanBurkhart/raytracer/pkg/raytracing"
)

type empty struct{}
type semaphore chan empty

// Scope provides the ability to point and target
type Scope struct {
	Position raytracing.Vector  `json:"position"`
	Target   *raytracing.Vector `json:"target"`
	Roll     float64            `json:"roll"`
	Right    *raytracing.Vector `json:"right"`
	Up       *raytracing.Vector `json:"up"`
	Forward  *raytracing.Vector `json:"forward"`
}

// GetUp returns the value of the normalized up vector for the scope
func (s *Scope) GetUp() raytracing.Vector {
	return *s.Up
}

// GetRight returns the value of the normalized right vector for the scope
func (s *Scope) GetRight() raytracing.Vector {
	return *s.Right
}

// GetForward returns the value of the normalized forward vector for the scope
func (s *Scope) GetForward() raytracing.Vector {
	return *s.Forward
}

// Initialize must be called before the Scope is used
func (s *Scope) Initialize() error {
	var ok bool
	var right, up, forward raytracing.Vector

	if s.Target != nil {
		forward, ok = s.Target.Subtract(s.Position).Normalize()
		if !ok {
			return fmt.Errorf("target and position are the same")
		}
		vertical, _ := raytracing.Vector{X: 0, Y: 1, Z: 0}.Normalize()
		right = forward.Cross(vertical)
		up = right.Cross(forward)

		var err error
		if up, err = up.Rotate(-s.Roll, forward); err != nil {
			return fmt.Errorf("scope rotation failed: %s", err)
		}
		right = forward.Cross(up)

		s.Up = &up
		s.Right = &right
		s.Forward = &forward
	}

	right, ok = s.Right.Normalize()
	if !ok {
		return fmt.Errorf("vector 'right' is a zero vector")
	}
	up, ok = s.Up.Normalize()
	if !ok {
		return fmt.Errorf("vector 'up' is a zero vector")
	}
	forward, ok = s.Forward.Normalize()
	if !ok {
		return fmt.Errorf("vector 'forward' is a zero vector")
	}

	s.Up = &up
	s.Right = &right
	s.Forward = &forward

	return nil
}

// Camera renders a scene using a specific view and perspective
type Camera struct {
	imageWidth  int
	imageHeight int

	output *image.RGBA

	AntiAliasingFactor *int `json:"antiAliasingFactor"`

	Lens
	Scope
}

// UnmarshalJSON unmarshals a Camera and resolves implementations of Lens
func (c *Camera) UnmarshalJSON(b []byte) error {
	type Alias Camera
	alias := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(b, &alias); err != nil {
		return err
	}

	if err := c.Scope.Initialize(); err != nil {
		return err
	}

	if c.AntiAliasingFactor != nil && *c.AntiAliasingFactor < 1 {
		return fmt.Errorf("anti-aliasing factor must be at least one ")
	}
	if c.AntiAliasingFactor == nil {
		antiAliasingFactor := 1
		c.AntiAliasingFactor = &antiAliasingFactor
	}

	var err error
	c.Lens, err = CreateLens(b)
	return err
}

// SetImageSize sets the width and height for rendered images
func (c *Camera) SetImageSize(width int, height int) (err error) {
	c.imageWidth = width
	c.imageHeight = height

	err = c.Lens.setAspectRatio(float64(width) / float64(height))
	if err != nil {
		return err
	}

	c.output = image.NewRGBA(image.Rect(0, 0, c.imageWidth, c.imageHeight))
	return
}

// Save encodes the internal image into a png file and writes to w
func (c *Camera) Save(w io.Writer) error {
	if c.output == nil {
		return fmt.Errorf("image must be rendered before saving it")
	}
	return png.Encode(w, c.output)
}

// Render creates a rendering of the Scene from the view of the Camera, use Save to save that image
func (c *Camera) Render(s *scene.Scene, maxRayReflections int, threads int) error {
	if c.output == nil {
		return fmt.Errorf("camera cannot perform render until image size is set (using SetImageSize)")
	}

	var wg sync.WaitGroup

	sema := make(semaphore, threads)

	antiAliasingIncrement := 1.0 / float64(*c.AntiAliasingFactor)

	for pixelY := 0; pixelY < c.imageHeight; pixelY++ {
		for pixelX := 0; pixelX < c.imageWidth; pixelX++ {
			var rays []raytracing.Ray
			for i := 0; i < *c.AntiAliasingFactor; i++ {
				for j := 0; j < *c.AntiAliasingFactor; j++ {
					pixelX := (float64(pixelX) + float64(i)*antiAliasingIncrement) / float64(c.imageWidth)
					pixelY := (float64(pixelY) + float64(j)*antiAliasingIncrement) / float64(c.imageHeight)
					screenX := 2.0*(pixelX) - 1.0
					screenY := -2.0*(pixelY) + 1.0
					ray := c.generateLightRay(screenX, screenY, c.Scope)
					rays = append(rays, ray)
				}
			}
			wg.Add(1)
			go c.renderRays(s, rays, pixelX, pixelY, maxRayReflections, &wg, sema)
		}
	}

	wg.Wait()

	return nil
}

// renderRay traces given starting rays through the scene and records the result. If a non-nil
// WaitGroup is passed in, Done will be called on it once the ray tracing is complete.
// This is threadsafe and can be executed in a goroutine.
func (c *Camera) renderRays(s *scene.Scene, rays []raytracing.Ray, pixelX int, pixelY int, maxRayReflections int, wg *sync.WaitGroup, sema semaphore) {
	if wg != nil {
		defer wg.Done()
	}

	sema <- empty{}

	var colors []raytracing.Color

	for _, ray := range rays {
		colors = append(colors, s.TraceRay(ray, 1.0, maxRayReflections))
	}

	pixelColor := raytracing.AverageColors(colors)

	red := math.Min(pixelColor.Red*255.0, 255.0)
	green := math.Min(pixelColor.Green*255.0, 255.0)
	blue := math.Min(pixelColor.Blue*255.0, 255.0)

	c.output.Set(pixelX, pixelY, color.RGBA{uint8(red), uint8(green), uint8(blue), 255.0})

	<-sema
}
