package engine

import (
	"math"
	"sync"

	"cinderwolf.net/raycaster/mapdata"
	"cinderwolf.net/raycaster/vector"
)

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
	if x < 0 {
		x = 0
	}
	if x > int(texWidth) {
		x = int(texWidth) - 1
	}

	if y < 0 {
		y = 0
	}
	if y > height {
		y = height - 1
	}

	return (x + (y * int(texWidth))) * 3
}

func (e *Engine) DoGeneratePixels() {
	var wg sync.WaitGroup
	cpus := width / 6

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
	if intersectionData == nil {
		return
	}
	e.drawFloorForCol(col, intersectionData.D, floorRay)
	e.renderLine(col, intersectionData)
}

func (e Engine) renderLine(col int, intersectionData *IntersectionData) {
	halfHeight := float64(height) / intersectionData.D / 2
	texture := e.MapData.Textures[*intersectionData.T]
	texCol := texture[int(intersectionData.X*10)%len(texture)]

	texPixHeight := float64(width) / intersectionData.D / 20
	for k := 0; k < 20; k += 1 {
		c := e.MapData.Colors[texCol[k%20]]
		yOffset := height/2 - int(halfHeight)
		for ty := 0; float64(ty) < texPixHeight; ty++ {
			idx := getIdx(col, ty+yOffset+int(texPixHeight*float64(k)))
			e.TextureSlice[idx+0] = c.R(0)
			e.TextureSlice[idx+1] = c.G(0)
			e.TextureSlice[idx+2] = c.B(0)
		}
	}
}

func (e Engine) drawFloorForCol(col int, dist float64, ray *Ray) {
	halfHeight := float64(height) / dist / 2
	minH := halfHeight

	texture := e.MapData.Textures[e.MapData.FloorTexture]
	for y := 0; float64(y) < float64(height)/2-minH; y++ {

		p := float64(y) - float64(height)/2
		posZ := float64(-height)

		rowDistance := posZ / p
		floorStepX := rowDistance * (ray.X2 - ray.X1) / float64(width)
		floorStepY := rowDistance * (ray.Y2 - ray.Y1) / float64(width)

		floorX := e.Player.Position.X + rowDistance*ray.X1
		floorY := e.Player.Position.Y + rowDistance*ray.Y1
		floorX += floorStepX * float64(col)
		floorY += floorStepY * float64(col)
		tx := ((uint8(10*floorX) % 20) + 20) % 20
		ty := ((uint8(10*floorY) % 20) + 20) % 20

		t := texture[ty][tx]
		c := e.MapData.Colors[t]

		idx1 := getIdx(col, y)
		idx2 := getIdx(col, height-y)

		e.TextureSlice[idx1+0] = c.R(1)
		e.TextureSlice[idx1+1] = c.G(1)
		e.TextureSlice[idx1+2] = c.B(1)

		e.TextureSlice[idx2+0] = c.R(0)
		e.TextureSlice[idx2+1] = c.G(0)
		e.TextureSlice[idx2+2] = c.B(0)
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

func (e Engine) collideAndMove(dir float64) {
	movedPos := e.Player.Position.Copy().Add(e.Player.Heading.Copy().Mul(dir * MOVEMENT_SPEED))
	var closestDist *float64
	var closestSegment mapdata.Segment
	for _, segment := range e.MapData.Segments {
		if e.orientation(vector.NewVector(segment.X1, segment.Y1), vector.NewVector(segment.X2, segment.Y2), &e.Player.Position) != 1 {
			continue
		}

		dist := e.distToLS(segment, *movedPos)

		// Not close enough to need to check for collision
		if dist > MOVEMENT_SPEED*1.5 {
			continue
		}
		if closestDist == nil || dist < *closestDist {
			closestDist = &dist
			closestSegment = segment
		}
	}
	if closestDist != nil && *closestDist < MOVEMENT_SPEED*1.5 {
		normal := vector.NewVector((-1 * (closestSegment.Y2 - closestSegment.Y1)), closestSegment.X2-closestSegment.X1)
		normal.Normalize()
		normal.Mul(
			float64(
				e.orientation(
					vector.NewVector(closestSegment.X1, closestSegment.Y1),
					vector.NewVector(closestSegment.X2, closestSegment.Y2),
					&e.Player.Position,
				),
			),
		)

		intersectionRay := &mapdata.Segment{
			X1: closestSegment.X1 + normal.X*-1*MOVEMENT_SPEED*1.5,
			Y1: closestSegment.Y1 + normal.Y*-1*MOVEMENT_SPEED*1.5,
			X2: closestSegment.X2 + normal.X*-1*MOVEMENT_SPEED*1.5,
			Y2: closestSegment.Y2 + normal.Y*-1*MOVEMENT_SPEED*1.5,
		}

		collisionRay := &Ray{}
		collisionRay.X1 = e.Player.Position.X
		collisionRay.Y1 = e.Player.Position.Y
		collisionRay.X2 = movedPos.X
		collisionRay.Y2 = movedPos.Y

		collision := e.getIntersectionData(collisionRay, intersectionRay, 0, false)
		collisionLocation := vector.NewVector(collision.Iix, collision.Iiy)
		length := movedPos.Copy().Sub(collisionLocation).Dot(normal) * dir
		movedPos.Add(normal.Copy().Mul(-1 * dir * length))

		for _, segment := range e.MapData.Segments {
			// TODO figure out what the heck I meant calling this "dist"
			dist := e.distToLS(segment, *movedPos)
			if dist < MOVEMENT_SPEED*1.5 {
				return
			}
		}
	}

	// for _, segment := range e.MapData.Segments {
	// 	if e.orientation(vector.NewVector(segment.X1, segment.Y1), vector.NewVector(segment.X2, segment.Y2), &e.Player.Position) != 1 {
	// 		continue
	// 	}

	// 	dist := e.distToLS(segment, *movedPos)

	// 	// Not close enough to need to check for collision
	// 	if dist > MOVEMENT_SPEED*1.5 {
	// 		continue
	// 	}
	// 	fmt.Printf("closestDist: %v\n", closestDist)
	// 	if closestDist == nil || dist < *closestDist {
	// 		closestDist = &dist
	// 		closestSegment = &segment
	// 	}
	// }

	e.Player.Position.X = movedPos.X
	e.Player.Position.Y = movedPos.Y
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

func (e Engine) orientation(a *vector.Vector, b *vector.Vector, c *vector.Vector) int {
	if (b.X-a.X)*(c.Y-a.Y)-(c.X-a.X)*(b.Y-a.Y) < 0 {
		return 1
	}
	return -1
}
