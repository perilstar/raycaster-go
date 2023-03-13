package shaders

import (
	"embed"
	"io"
)

//go:embed pass1.vert pass1.frag
var f embed.FS

func LoadPass1() (string, string) {
	vertFile, _ := f.Open("pass1.vert")
	fragFile, _ := f.Open("pass1.frag")
	defer vertFile.Close()
	defer fragFile.Close()

	vertByteValue, _ := io.ReadAll(vertFile)
	vertByteValue = append(vertByteValue, '\x00')
	vertSource := string(vertByteValue)

	fragByteValue, _ := io.ReadAll(fragFile)
	fragByteValue = append(fragByteValue, '\x00')
	fragSource := string(fragByteValue)

	return vertSource, fragSource
}
