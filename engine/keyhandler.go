package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func (e *Engine) HandleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape:
		if action == glfw.Press {
			w.SetShouldClose(true)
		}
	case glfw.KeyW, glfw.KeyA, glfw.KeyS, glfw.KeyD:
		if action == glfw.Press {
			e.Keys[key] = true
		} else if action == glfw.Release {
			e.Keys[key] = false
		}
	}
}
