// Linear algebra
// Right-handed Cartesian coordinate system

package raytracing

import (
	"fmt"
	"math"
)

// Vector is a 3 dimensional component representation of a vector
type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Dot product
func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross product
func (v Vector) Cross(other Vector) Vector {
	return Vector{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Magnitude of vector
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// Add returns the sum of this and the other vector
func (v Vector) Add(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Subtract returns the difference of this and the other vector
func (v Vector) Subtract(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Negative returns the negation of the vector
func (v Vector) Negative() Vector {
	return Vector{
		X: -v.X,
		Y: -v.Y,
		Z: -v.Z,
	}
}

// Scale returns a scaled vector
func (v Vector) Scale(scale float64) Vector {
	return Vector{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale}
}

// Normalize returns a boolean indicating success
func (v Vector) Normalize() (Vector, bool) {
	mag := v.Magnitude()
	// Zero-length vector has no direction and therefore can not be normalized
	if mag <= 0.0 {
		return v, false
	}
	return v.Scale(1.0 / mag), true
}

// Rotate returns the vector rotated around the specified axis
// Rotation is counter-clockwise when axis vector points towards observer
func (v Vector) Rotate(degrees float64, axis Vector) (Vector, error) {
	var ok bool
	if axis, ok = axis.Normalize(); !ok {
		return Vector{}, fmt.Errorf("the zero vector cannot be used as an axis of rotation")
	}

	angle := degrees / 180.0 * math.Pi

	rotationMatrix := [][]float64{
		{
			math.Cos(angle) + axis.X*axis.X*(1-math.Cos(angle)),
			axis.X*axis.Y*(1-math.Cos(angle)) - axis.Z*math.Sin(angle),
			axis.X*axis.Z*(1-math.Cos(angle)) + axis.Y*math.Sin(angle),
		}, {
			axis.Y*axis.X*(1-math.Cos(angle)) + axis.Z*math.Sin(angle),
			math.Cos(angle) + axis.Y*axis.Y*(1-math.Cos(angle)),
			axis.Y*axis.Z*(1-math.Cos(angle)) - axis.X*math.Sin(angle),
		}, {
			axis.Z*axis.X*(1-math.Cos(angle)) - axis.Y*math.Sin(angle),
			axis.Z*axis.Y*(1-math.Cos(angle)) + axis.X*math.Sin(angle),
			math.Cos(angle) + axis.Z*axis.Z*(1-math.Cos(angle)),
		},
	}

	return Multiply(rotationMatrix, v), nil
}

// Multiply will return the vector result of multiplying the matrix by the vector as a column vector
// Assumes matrix has correct dimensions to multiply with vector
func Multiply(matrix [][]float64, vector Vector) (product Vector) {
	product.X = matrix[0][0]*vector.X + matrix[0][1]*vector.Y + matrix[0][2]*vector.Z
	product.Y = matrix[1][0]*vector.X + matrix[1][1]*vector.Y + matrix[1][2]*vector.Z
	product.Z = matrix[2][0]*vector.X + matrix[2][1]*vector.Y + matrix[2][2]*vector.Z
	return
}
