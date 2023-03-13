package mapdata

import (
	"embed"
	"encoding/json"
	"io"
	"log"
)

//go:embed map.json
var f embed.FS

type MapData struct {
	Colors       []ColorHSL `json:"p"`
	Textures     []Texture  `json:"t"`
	Segments     []Segment  `json:"s"`
	FloorTexture int        `json:"f"`
}

func Load() *MapData {
	m := &MapData{}
	file, _ := f.Open("map.json")
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	err := json.Unmarshal(byteValue, m)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v\n", err)
	}

	return m
}
