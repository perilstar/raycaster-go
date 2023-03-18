package engine

import (
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
	Shader1        uint32
	Shader2        uint32
	Player         *Player
	MapData        *mapdata.MapData
	Keys           map[glfw.Key]bool
}

const (
	MOVEMENT_SPEED float64 = 0.3
	TURN_SPEED     float64 = 2.0
	FOV            float64 = 90 * 0.0174533
)

var (
	width     int = 1200
	height    int = 800
	texWidth  uint32
	texHeight uint32
)

func NewEngine() *Engine {
	return &Engine{
		Frames: 0,
		Keys: map[glfw.Key]bool{
			glfw.KeyW: false,
			glfw.KeyA: false,
			glfw.KeyS: false,
			glfw.KeyD: false,
		},
		Player: &Player{
			Position: *vector.NewVector(3, 3),
			Heading:  *vector.NewVector(MOVEMENT_SPEED, 0).SetHeading(0.0174529252 * 30),
		},
	}
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
	e.MapData = mapdata.Load()

	updateTexSize()

	if err := glfw.Init(); err != nil {
		log.Fatalf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)

	glfw.WindowHint(glfw.Focused, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(int(width), int(height), "Cinder's Raycaster", nil, nil)

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

	go func() {
		for {
			time.Sleep(time.Second / 50)
			if e.Keys[glfw.KeyA] {
				e.Player.Heading.SetHeading(e.Player.Heading.Heading() + 0.01745329252*(-1*TURN_SPEED))
			}
			if e.Keys[glfw.KeyD] {
				e.Player.Heading.SetHeading(e.Player.Heading.Heading() + 0.01745329252*(TURN_SPEED))
			}

			if e.Keys[glfw.KeyW] {
				e.collideAndMove(1.0)
			}
			if e.Keys[glfw.KeyS] {
				e.collideAndMove(-1.0)
			}
		}
	}()

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
