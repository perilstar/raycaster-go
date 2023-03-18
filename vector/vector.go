package vector

import (
	"math"
)

type Vector struct {
	X float64
	Y float64
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}

func (v *Vector) Set(other *Vector) *Vector {
	v.X = other.X
	v.Y = other.Y
	return v
}

func (v *Vector) Copy() *Vector {
	return NewVector(v.X, v.Y)
}

func (v *Vector) Add(other *Vector) *Vector {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v *Vector) Sub(other *Vector) *Vector {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v *Vector) Mul(scalar float64) *Vector {
	v.X *= scalar
	v.Y *= scalar
	return v
}

func (v *Vector) Div(scalar float64) *Vector {
	v.X /= scalar
	v.Y /= scalar
	return v
}

func (v *Vector) Dot(o *Vector) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v *Vector) Mag() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Normalize() *Vector {
	mag := v.Mag()

	if mag == 0 {
		return v
	}

	return v.Div(mag)
}

func (v *Vector) Heading() float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func (v *Vector) SetMag(length float64) *Vector {
	return v.Normalize().Mul(length)
}

func (v *Vector) SetHeading(angle float64) *Vector {
	mag := v.Mag()
	v.X = math.Cos(angle)
	v.Y = math.Sin(angle)
	return v.Normalize().Mul(mag)
}

func (v *Vector) Abs() *Vector {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	return v
}
