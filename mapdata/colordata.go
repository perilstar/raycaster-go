package mapdata

import (
	"encoding/json"
	"math"
)

type ColorData [2][3]uint8

type ColorHSL [3]uint
type ColorRGB [3]uint8

func (c *ColorData) UnmarshalJSON(data []byte) error {
	var v [3]uint

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	c[0] = *hslToRgb(v)
	c[1] = *hslToRgb([3]uint{v[0], v[1], uint(math.Min(float64(v[2])*1.5, 100.0))})

	return nil
}

func hue2rgb(p float32, q float32, t float32) uint8 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return uint8(255*p + (q-p)*6*t)
	}
	if t < 1.0/2.0 {
		return uint8(255 * q)
	}
	if t < 2.0/3.0 {
		return uint8(255*p + (q-p)*(2.0/3.0-t)*6.0)
	}
	return uint8(255 * p)
}

func hslToRgb(color [3]uint) *[3]uint8 {
	h := float32(color[0]) / 360
	s := float32(color[1]) / 100
	l := float32(color[2]) / 100

	if s == 0 {
		channel := uint8(l * 255)
		return &[3]uint8{channel, channel, channel}
	}

	var q float32
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}

	p := 2*l - q
	r := hue2rgb(p, q, h+1.0/3.0)
	g := hue2rgb(p, q, h)
	b := hue2rgb(p, q, h-1.0/3.0)
	return &[3]uint8{r, g, b}
}

func (c ColorData) R(l uint8) uint8 {
	return c[l][0]
}

func (c ColorData) G(l uint8) uint8 {
	return c[l][1]
}

func (c ColorData) B(l uint8) uint8 {
	return c[l][2]
}
