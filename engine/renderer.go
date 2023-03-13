package engine

import (
	"log"

	"cinderwolf.net/raycaster/shaders"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func (e *Engine) DoSetup(window *glfw.Window) {
	width, height = window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(texWidth), int32(texHeight))

	var textureStorage TextureStorage = TextureStorage{}
	gl.GenTextures(1, &textureStorage.Texture)
	gl.BindTexture(gl.TEXTURE_2D, textureStorage.Texture)
	e.TextureStorage = &textureStorage

	e.TextureSlice = make([]uint8, texWidth*texHeight*3)

	gl.Enable(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	updateTexSize()
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

	vert1, frag1 := shaders.LoadPass1()
	shader1, err := CreateProgram(vert1, frag1)
	if err != nil {
		log.Fatalf("shader1 fail: %v\n", err)
	}

	e.Shader1 = shader1
}

func (e *Engine) DoCleanup(window *glfw.Window) {
	gl.Disable(gl.TEXTURE_2D)
}

func (e *Engine) DrawScene() {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	e.DoGeneratePixels()
	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0,
		0,
		0,
		int32(texWidth),
		int32(texHeight),
		gl.RGB,
		gl.UNSIGNED_BYTE,
		gl.Ptr(e.TextureSlice),
	)

	vertices := []float32{
		-1, -1, 0,
		0, 0,
		1, -1, 0,
		1, 0,
		1, 1, 0,
		1, 1,
		-1, 1, 0,
		0, 1,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.UseProgram(e.Shader1)
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)
}
