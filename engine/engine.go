package engine

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
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
}

var (
	width  = 379
	height = 377
)

func NewEngine() *Engine {
	return &Engine{Frames: 0}
}

func (e *Engine) Start() {
	if err := glfw.Init(); err != nil {
		log.Fatalf("failed to initialize glfw: %v", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.Focused, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Test", nil, nil)

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

		gl.Viewport(0, 0, int32(width), int32(height))

		newTextureSlice := make([]uint8, width*height*3)

		e.TextureSlice = newTextureSlice
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
