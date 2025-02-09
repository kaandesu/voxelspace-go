package main

import (
	"image"
	"image/png"
	"log/slog"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	Scene struct {
		colorMap      image.Image
		camera        *Camera
		heightMapPath string
		colorMapPath  string
		heightMaps    [][]uint8
		widht         int
		height        int
	}
	Camera struct {
		position     [2]int
		height       int
		horizon_pos  int
		scale_height int
		max_dist     int
	}
)

const (
	SCREEN_WIDHT  = 640
	SCREEN_HEIGHT = 480
)

func NewScene(mapPath, colorPath string) *Scene {
	return &Scene{
		heightMapPath: mapPath,
		colorMapPath:  colorPath,
	}
}

func (scene *Scene) loadHeightMap() {
	file, err := os.Open(scene.heightMapPath)
	if err != nil {
		slog.Error("Could not open the file", "err", err)
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		slog.Error("Could not decode the heightmap", "err", err)
		return
	}

	bounds := img.Bounds()

	scene.widht, scene.height = bounds.Dx(), bounds.Dy()
	scene.heightMaps = make([][]uint8, scene.height)

	for y := range scene.height {
		scene.heightMaps[y] = make([]uint8, scene.widht)
		for x := range scene.widht {
			r, _, _, _ := img.At(x, y).RGBA()
			scene.heightMaps[y][x] = uint8(r)
		}
	}
}

func (scene *Scene) loadColorMap() {
	file, err := os.Open(scene.colorMapPath)
	if err != nil {
		slog.Error("Could not open the file", "err", err)
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		slog.Error("Could not decode the color map", "err", err)
		return
	}

	scene.colorMap = img
}

func (scene *Scene) LoadSetup() {
	scene.loadHeightMap()
	scene.loadColorMap()
}

func main() {
	scene := NewScene("./maps/D1.png", "./maps/C1W.png")
	scene.LoadSetup()
	rl.InitWindow(SCREEN_WIDHT, SCREEN_HEIGHT, "voxelspace")

	rl.SetTargetFPS(60)
	renderTexture := rl.LoadRenderTexture(700, 700)
	defer rl.UnloadRenderTexture(renderTexture)

	for !rl.WindowShouldClose() {
		rl.BeginTextureMode(renderTexture)
		rl.ClearBackground(rl.SkyBlue)

		rl.EndTextureMode()

		rl.BeginDrawing()
		rl.DrawTexturePro(renderTexture.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  float32(renderTexture.Texture.Width),
				Height: -float32(renderTexture.Texture.Height),
			},
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  700,
				Height: 700,
			},
			rl.Vector2{
				X: 0,
				Y: 0,
			},
			0,
			rl.White)

		rl.EndDrawing()
	}
}
