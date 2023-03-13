package engine

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"cinderwolf.net/raycaster/mapdata"
	"cinderwolf.net/raycaster/vector"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type TextureStorage struct {
	Texture uint32
}

type Engine struct {
	TextureSlice   []uint8
	Frames         int
	Window         glfw.Window
	TextureStorage *TextureStorage
	Ctx            *context.Context
	Cancel         *context.CancelFunc
	Shader1        uint32
	Shader2        uint32
	Player         *Player
	MapData        *mapdata.MapData
}

var (
	width     int = 379
	height    int = 377
	texWidth  uint32
	texHeight uint32
)

func NewEngine() *Engine {
	return &Engine{Frames: 0}
}

func updateTexSize() {
	texWidth = uint32(width)
	texHeight = uint32(height)

	texWidth--
	texWidth |= texWidth >> 1
	texWidth |= texWidth >> 2
	texWidth |= texWidth >> 4
	texWidth |= texWidth >> 8
	texWidth |= texWidth >> 16
	texWidth++

	texHeight--
	texHeight |= texHeight >> 1
	texHeight |= texHeight >> 2
	texHeight |= texHeight >> 4
	texHeight |= texHeight >> 8
	texHeight |= texHeight >> 16
	texHeight++
}

func (e *Engine) Start() {
	e.Player = &Player{
		Position: *vector.NewVector(3, 3),
		Heading:  *vector.NewVector(MOVEMENT_SPEED, 0).SetHeading(0.0174529252 * 30),
	}

	e.MapData = mapdata.Load()

	updateTexSize()
	fmt.Printf("texWidth: %v\n", texWidth)
	fmt.Printf("texHeight: %v\n", texHeight)

	if err := glfw.Init(); err != nil {
		log.Fatalf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)

	glfw.WindowHint(glfw.Focused, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(int(width), int(height), "Test", nil, nil)

	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	runtime.LockOSThread()

	window.SetSizeLimits(200, 200, glfw.DontCare, glfw.DontCare)

	go func(e *Engine) {
		lastFrameCount := 0
		for {
			time.Sleep(time.Second)
			fmt.Printf("FPS: %d   Frames: %d\n", e.Frames-lastFrameCount, e.Frames)

			lastFrameCount = e.Frames

		}
	}(e)

	window.SetKeyCallback(e.HandleKey)

	window.SetFramebufferSizeCallback(func(window *glfw.Window, newWidth int, newHeight int) {
		width = newWidth
		height = newHeight
		updateTexSize()

		gl.Viewport(0, 0, int32(texWidth), int32(texHeight))

		newTextureSlice := make([]uint8, texWidth*texHeight*3)

		e.TextureSlice = newTextureSlice

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGB,
			int32(texWidth),
			int32(texHeight),
			0,
			gl.RGB,
			gl.UNSIGNED_BYTE,
			nil,
		)
	})
	e.DoSetup(window)

	glfw.SwapInterval(1)
	for !window.ShouldClose() {
		e.DrawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	e.DoCleanup(window)

}
