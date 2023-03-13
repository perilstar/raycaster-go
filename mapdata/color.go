package mapdata

type Color [3]uint

type ColorHSL Color
type ColorRGB Color

func (c ColorHSL) H() uint {
	return c[0]
}

func (c ColorHSL) S() uint {
	return c[1]
}

func (c ColorHSL) L() uint {
	return c[2]
}

func (c ColorRGB) R() uint {
	return c[0]
}

func (c ColorRGB) G() uint {
	return c[1]
}

func (c ColorRGB) B() uint {
	return c[2]
}
