package mapdata

type Color [3]uint

func (c Color) R() uint {
	return c[0]
}

func (c Color) G() uint {
	return c[1]
}

func (c Color) B() uint {
	return c[2]
}
