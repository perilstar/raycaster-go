package engine

import (
	"math"
	"runtime"
	"sync"

	"cinderwolf.net/raycaster/mapdata"
	"cinderwolf.net/raycaster/vector"
	"go.timothygu.me/math/v2/imath"
)

const MOVEMENT_SPEED float64 = 0.1
const TURN_SPEED float64 = 2
const FOV float64 = 90 * 0.0174533

type Ray struct {
	X1 float64
	Y1 float64
	X2 float64
	Y2 float64
}
type Player struct {
	Position vector.Vector
	Heading  vector.Vector
}

type IntersectionData struct {
	D   float64
	X   float64
	T   *uint
	L   *mapdata.Segment
	Iix float64
	Iiy float64
}

func getIdx(x int, y int) int {
	return (x + (y * int(texWidth))) * 3
}

func (e *Engine) DoGeneratePixels() {
	var wg sync.WaitGroup
	cpus := imath.Min(runtime.NumCPU(), 32)

	for i := 0; i < cpus; i++ {
		wg.Add(1)

		var start int
		if i == 0 {
			start = 0
		} else {
			start = int(float64(i) / float64(cpus) * float64(width))
		}
		var end int
		if i+1 == cpus {
			end = int(width)
		} else {
			end = int(float64(i+1) / float64(cpus) * float64(width))
		}

		minAngle := e.getAngle(0)
		maxAngle := e.getAngle(width - 1)

		floorRay := &Ray{}
		floorRay.X1 = math.Cos(minAngle + e.Player.Heading.Heading())
		floorRay.Y1 = math.Sin(minAngle + e.Player.Heading.Heading())
		floorRay.X2 = math.Cos(maxAngle + e.Player.Heading.Heading())
		floorRay.Y2 = math.Sin(maxAngle + e.Player.Heading.Heading())

		go func(start int, end int) {
			defer wg.Done()
			for x := start; x < end; x++ {
				idx := getIdx(x, height)
				if idx >= int(texWidth*texHeight*3) {
					return
				}
				// e.TextureSlice[idx+0] = uint8(float64(x) / float64(width) * 255)
				// e.TextureSlice[idx+1] = uint8(float64(y) / float64(height) * 255)
				// e.TextureSlice[idx+2] = uint8((math.Sin(float64(e.Frames)*0.03) + 1) * 0.5 * 255)
				e.drawCol(x, floorRay)
			}

		}(start, end)
	}
	wg.Wait()

	e.Frames++
}

func (e Engine) getDist() float64 {
	return float64(width) / 2 / math.Tan(float64(FOV)/2)
}

func (e Engine) getAngle(col int) float64 {
	return math.Atan2(float64(col-width/2), e.getDist())
}

func (e Engine) castRay(angle float64, col int) *IntersectionData {
	ray := Ray{}
	ray.X1 = e.Player.Position.X
	ray.Y1 = e.Player.Position.Y
	ray.X2 = e.Player.Position.X + math.Cos(e.Player.Heading.Heading()+angle)
	ray.Y2 = e.Player.Position.Y + math.Sin(e.Player.Heading.Heading()+angle)

	var closestIntersection *IntersectionData
	for _, segment := range e.MapData.Segments {
		if e.orientation(
			vector.NewVector(segment.X1, segment.Y1),
			vector.NewVector(segment.X2, segment.Y2),
			&e.Player.Position,
		) != 1 {
			continue
		}
		intersectionData := e.getIntersectionData(&ray, &segment, angle, true)
		if intersectionData != nil && (closestIntersection == nil || (intersectionData.D < closestIntersection.D)) {
			closestIntersection = intersectionData
		}
	}
	return closestIntersection
}

/*
			 * Horrible abomination that has been ported like 3 times now
       * Don't be like me, use sensible variable names
			 * I have no idea what half of this does anymore
*/
func (e Engine) getIntersectionData(r *Ray, l *mapdata.Segment, angle float64, checkWithinEndpoints bool) *IntersectionData {
	i1to2x := r.X2 - r.X1
	i1to2y := r.Y2 - r.Y1
	i3to4x := l.X2 - l.X1
	i3to4y := l.Y2 - l.Y1

	var intersectionData *IntersectionData

	if (i1to2y / i1to2x) == (i3to4y / i3to4x) {
		return intersectionData
	}

	id := i1to2x*i3to4y - i1to2y*i3to4x
	if id == 0 {
		return intersectionData
	}

	i3to1x := r.X1 - l.X1
	i3to1y := r.Y1 - l.Y1
	ir := (i3to1y*i3to4x - i3to1x*i3to4y) / id
	is := (i3to1y*i1to2x - i3to1x*i1to2y) / id

	if checkWithinEndpoints && (ir < 0 || is < 0 || is > 1) {
		return intersectionData
	}

	iix := (r.X2-r.X1)*ir + r.X1
	iiy := (r.Y2-r.Y1)*ir + r.Y1

	dx := iix - r.X1
	dy := iiy - r.Y1

	idx := iix - l.X1
	idy := iiy - l.Y1

	io := math.Sqrt(idx*idx + idy*idy)
	intersectionData = &IntersectionData{}
	intersectionData.D = math.Sqrt(dx*dx+dy*dy) * math.Cos(angle)
	intersectionData.X = io
	intersectionData.T = l.TextureIdx
	intersectionData.L = l
	intersectionData.Iix = iix
	intersectionData.Iiy = iiy

	return intersectionData
}

func (e Engine) drawCol(col int, floorRay *Ray) {
	intersectionData := e.castRay(e.getAngle(col), col)
	// e.renderLine(col, intersectionData)
	e.drawFloorForCol(intersectionData.D, floorRay)
}

func (e Engine) renderLine(col int, intersectionData *IntersectionData) {
	halfHeight := int(float64(height)/intersectionData.D) / 2
	texture := e.MapData.Textures[*intersectionData.T]
	texCol := texture[int(intersectionData.X*10)%len(texture)]
	tsh := float64(width) / intersectionData.D / 20
	jStart := col
	jEnd := col + 1
	for j := jStart; j < jEnd; j++ {
		for k := 0; k < 20; k += 1 {
			c := e.hslToRgb(e.MapData.Colors[texCol[k|0%20]])
			pStart := uint32((height/2-halfHeight+(int(float64(k)*tsh)))*width + j)
			pEnd := uint32((height/2-halfHeight+(int(float64(k+1)*tsh)))*width + j)
			for p := pStart; p < pEnd; p += texWidth * 3 {
				e.TextureSlice[p+0] = uint8(c.R())
				e.TextureSlice[p+1] = uint8(c.G())
				e.TextureSlice[p+2] = uint8(c.B())
			}
		}
	}
}

func (e Engine) drawFloorForCol(dist float64, ray *Ray) {
	halfHeight := int((float64(height) / dist / 2))

	minH := halfHeight

	texture := e.MapData.Textures[e.MapData.FloorTexture]
	for y := 0; y < height/2-minH; y++ {

		p := float64(y) - float64(height)/2
		posZ := float64(-height)

		rowDistance := posZ / p
		floorStepX := rowDistance * (ray.X2 - ray.X1) / float64(width)
		floorStepY := rowDistance * (ray.Y2 - ray.Y1) / float64(width)

		floorX := e.Player.Position.X + rowDistance*ray.X1
		floorY := e.Player.Position.Y + rowDistance*ray.Y1

		//TODO figure out if x should start at 1
		for x := 1; x < width; x++ {
			tx := ((int(10*floorX) % 20) + 20) % 20
			ty := ((int(10*floorY) % 20) + 20) % 20
			floorX += floorStepX
			floorY += floorStepY

			t := texture[ty][tx]
			hsl := e.MapData.Colors[t]
			hslLighter := mapdata.ColorHSL{hsl.H(), hsl.S(), uint(math.Min(float64(hsl.L())*1.5, 100))}

			c := e.hslToRgb(hsl)
			cl := e.hslToRgb(hslLighter)

			idx1 := getIdx(x, y)
			idx2 := getIdx(x, height-y-1)

			e.TextureSlice[idx1+0] = uint8(cl.R())
			e.TextureSlice[idx1+1] = uint8(cl.G())
			e.TextureSlice[idx1+2] = uint8(cl.B())

			e.TextureSlice[idx2+0] = uint8(c.R())
			e.TextureSlice[idx2+1] = uint8(c.G())
			e.TextureSlice[idx2+2] = uint8(c.B())
			// buf32[x+width*y] =
			// 				(255 << 24) |
			// 					(cl[2] << 16) |
			// 					(cl[1] << 8) |
			// 					cl[0]
			// 			buf32[x+width*(height-y-1)] =
			// 				(255 << 24) |
			// 					(c[2] << 16) |
			// 					(c[1] << 8) |
			// 					c[0]
		}
	}
}

func (e Engine) collide(oldPos vector.Vector, dir float64) {
	var closestDist *float64
	var closestSegment *mapdata.Segment
	for _, segment := range e.MapData.Segments {
		if e.orientation(vector.NewVector(segment.X1, segment.Y1), vector.NewVector(segment.X2, segment.Y2), &oldPos) != 1 {
			continue
		}

		dist := e.distToLS(segment, e.Player.Position)

		// Not close enough to need to check for collision
		if dist > MOVEMENT_SPEED {
			continue
		}

		if closestDist != nil || (dist != -1 && dist < *closestDist) {
			closestDist = &dist
			closestSegment = &segment
		}
	}

	if *closestDist < MOVEMENT_SPEED {
		normal := vector.NewVector((-1 * (closestSegment.Y2 - closestSegment.X2)), closestSegment.Y1-closestSegment.X1)
		normal.Normalize()
		normal.Mul(
			float64(
				e.orientation(
					vector.NewVector(closestSegment.X1, closestSegment.Y1),
					vector.NewVector(closestSegment.X2, closestSegment.Y2),
					&oldPos,
				),
			),
		)

		intersectionRay := &mapdata.Segment{
			X1: closestSegment.X1 + normal.X*-1*MOVEMENT_SPEED,
			Y1: closestSegment.Y1 + normal.Y*-1*MOVEMENT_SPEED,
			X2: closestSegment.X2 + normal.X*-1*MOVEMENT_SPEED,
			Y2: closestSegment.Y1 + normal.Y*-1*MOVEMENT_SPEED,
		}

		collisionRay := &Ray{}
		collisionRay.X1 = oldPos.X
		collisionRay.Y1 = oldPos.Y
		collisionRay.X2 = e.Player.Position.X
		collisionRay.Y2 = e.Player.Position.Y

		collision := e.getIntersectionData(collisionRay, intersectionRay, 0, false)
		collisionLocation := vector.NewVector(collision.Iix, collision.Iiy)

		length := e.Player.Position.Copy().Sub(collisionLocation).Dot(normal) * dir

		e.Player.Position.Add(normal.Copy().Mul(-1 * dir * length))

		collisionDuringSlide := false

		for _, ls := range e.MapData.Segments {
			// TODO figure out what the heck I meant calling this "thing"
			thing := e.distToLS(ls, e.Player.Position)
			if thing < MOVEMENT_SPEED*0.8 {
				collisionDuringSlide = true
			}
		}

		if collisionDuringSlide {
			e.Player.Position = oldPos
		}
		return
	}
}

func (e Engine) distToLS(l mapdata.Segment, c vector.Vector) float64 {
	p1 := vector.NewVector(l.X1, l.Y1)
	p2 := vector.NewVector(l.X2, l.Y2)
	ia := c.X - p1.X
	ib := c.Y - p1.Y
	ic := p2.X - p1.X
	id := p2.Y - p1.Y

	dot := ia*ic + ib*id
	len_sq := ic*ic + id*id

	param := -1.0
	if len_sq != 0 {
		param = dot / len_sq
	}

	var xx float64
	var yy float64

	if param < 0 {
		xx = p1.X
		yy = p1.Y
	} else if param > 1 {
		xx = p2.X
		yy = p2.Y
	} else {
		xx = p1.X + param*ic
		yy = p1.Y + param*id
	}

	dx := c.X - xx
	dy := c.Y - yy

	d := math.Sqrt(dx*dx + dy*dy)
	return d
}

func (e Engine) hslToRgb(hsl mapdata.ColorHSL) *mapdata.ColorRGB {
	h := float64(hsl.H()) / 360
	s := float64(hsl.S()) / 100
	l := float64(hsl.L()) / 100
	var r float64
	var g float64
	var b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {

		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}

		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}
	return &mapdata.ColorRGB{uint(r * 255), uint(g * 255), uint(b * 255)}
}

func (e Engine) orientation(a *vector.Vector, b *vector.Vector, c *vector.Vector) int {
	if (b.X-a.X)*(c.Y-a.Y)-(c.X-a.X)*(b.Y-a.Y) < 0 {
		return 1
	}
	return -1
}

func hue2rgb(p float64, q float64, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}
