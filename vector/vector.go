package vector

import "math"

type Vector struct {
	x float64
	y float64
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{x, y}
}

func (v Vector) copy() *Vector {
	return NewVector(v.x, v.y)
}

func (v Vector) add(other *Vector) *Vector {
	v.x += other.x
	v.y += other.y
	return &v
}

func (v Vector) sub(other *Vector) *Vector {
	v.x -= other.x
	v.y -= other.y
	return &v
}

func (v Vector) mul(scalar float64) *Vector {
	v.x *= scalar
	v.y *= scalar
	return &v
}

func (v Vector) div(scalar float64) *Vector {
	v.x /= scalar
	v.y /= scalar
	return &v
}

func (v Vector) dot(o *Vector) float64 {
	return v.x*o.x + v.y*o.y
}

func (v Vector) mag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v Vector) normalize() *Vector {
	mag := v.mag()

	if mag == 0 {
		return &v
	}

	return v.div(mag)
}

func (v Vector) heading() float64 {
	angle := math.Atan2(v.y, v.x)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func (v Vector) setMag(length float64) *Vector {
	return v.normalize().mul(length)
}

func (v Vector) setHeading(angle float64) *Vector {
	mag := v.mag()
	v.x = math.Cos(angle)
	v.y = math.Sin(angle)
	return v.normalize().mul(mag)
}

func (v Vector) abs() *Vector {
	v.x = math.Abs(v.x)
	v.y = math.Abs(v.y)
	return &v
}
