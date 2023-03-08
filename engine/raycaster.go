package engine

import (
	"math"
	"runtime"
	"sync"
)

func getIdx(x int, y int) int {
	return (x + (y * width)) * 3
}

func (e *Engine) DoGeneratePixels() {
	var wg sync.WaitGroup
	cpus := runtime.NumCPU()

	for i := 0; i < cpus; i++ {

		wg.Add(1)

		start := int(float64(i) / float64(cpus) * float64(width))
		var end int
		if i+1 == cpus {
			end = width
		} else {

			end = int(float64(i+1) / float64(cpus) * float64(width))
		}

		go func(start int, end int) {
			defer wg.Done()

			for x := start; x < end; x++ {
				for y := 0; y < height; y++ {
					idx := getIdx(x, y)
					if idx >= len(e.TextureSlice) {
						return
					}
					e.TextureSlice[idx+0] = uint8(float64(x%100) / float64(width) * 255)
					e.TextureSlice[idx+1] = uint8(float64(y%100) / float64(height) * 255)
					e.TextureSlice[idx+2] = uint8((math.Sin(float64(e.Frames)*0.03) + 1) * 0.5 * 255)
				}
			}

		}(start, end)
	}
	wg.Wait()

	e.Frames++
}
