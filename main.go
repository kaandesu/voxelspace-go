package main

import (
	"image"
	"image/png"
	"log/slog"
	"math"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDHT  = 640
	SCREEN_HEIGHT = 480
)

type (
	Scene struct {
		colorMap      image.Image
		camera        *Camera
		heightMapPath string
		colorMapPath  string
		heightMaps    [][]uint8
		yBuffer       []float32
		widht         int
		height        int
	}
	Camera struct {
		position struct {
			x, y int
		}
		speed        int
		height       int
		horizon_pos  int
		scale_height int
		max_dist     int
	}
)

func NewScene(mapPath, colorPath string) *Scene {
	buffer := make([]float32, SCREEN_WIDHT)
	for i := range SCREEN_WIDHT {
		buffer[i] = SCREEN_HEIGHT
	}
	return &Scene{
		heightMapPath: mapPath,
		colorMapPath:  colorPath,
		yBuffer:       buffer,
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

func (scene *Scene) initCamera() {
	scene.camera = &Camera{
		position:     struct{ x, y int }{x: scene.widht / 2, y: scene.height / 2},
		speed:        3,
		horizon_pos:  120,
		height:       100,
		scale_height: 200,
		max_dist:     700,
	}
}

func (scene *Scene) updateCamera() {
	camSpeed := scene.camera.speed
	if rl.IsKeyDown(rl.KeyW) {
		scene.camera.position.y -= camSpeed
	}
	if rl.IsKeyDown(rl.KeyS) {
		scene.camera.position.y += camSpeed
	}
	if rl.IsKeyDown(rl.KeyA) {
		scene.camera.position.x -= camSpeed
	}
	if rl.IsKeyDown(rl.KeyD) {
		scene.camera.position.x += camSpeed
	}

	if rl.IsKeyDown(rl.KeyE) {
		scene.camera.height += camSpeed
	}
	if rl.IsKeyDown(rl.KeyQ) {
		scene.camera.height -= camSpeed
	}

	if rl.IsKeyDown(rl.KeyUp) {
		scene.camera.horizon_pos += camSpeed * 2
	}
	if rl.IsKeyDown(rl.KeyDown) {
		scene.camera.horizon_pos -= camSpeed * 2
	}
}

func (scene *Scene) LoadSetup() {
	scene.loadHeightMap()
	scene.loadColorMap()
	scene.initCamera()
}

// NOTE: lower is more queality less peformance
const QoL int = 10

func (scene *Scene) render() {
	for i := range scene.yBuffer {
		scene.yBuffer[i] = SCREEN_HEIGHT
	}
	for z := 0; z < scene.camera.max_dist; {
		z += ((z + 1) * QoL / scene.camera.max_dist) + 1
		pleftX := -z + scene.camera.position.x
		pleftY := -z + scene.camera.position.y
		prightX := z + scene.camera.position.x
		// prightY := -z + scene.camera.position[1]

		dx := float32(prightX-pleftX) / float32(SCREEN_WIDHT)
		px := float32(pleftX)

		for i := 0; i < SCREEN_WIDHT; i++ {
			x := int(px)
			y := pleftY

			if x >= 0 && x < scene.widht && y >= 0 && y < scene.height {
				heightOnScreen := float32(scene.camera.height-int(scene.heightMaps[y][x])) / float32(z) * float32(scene.camera.scale_height)
				heightOnScreen += float32(scene.camera.horizon_pos)

				// NOTE: line below makes the height-map convex
				heightOnScreen = heightOnScreen * float32(math.Abs(float64(scene.camera.max_dist+z))/float64(scene.camera.max_dist))

				color := scene.colorMap.At(x, y)
				r, g, b, _ := color.RGBA()
				col := rl.NewColor(uint8(r), uint8(g), uint8(b), 255)

				if heightOnScreen < scene.yBuffer[i] {
					rl.DrawLine(int32(i), int32(heightOnScreen), int32(i), int32(scene.yBuffer[i]), col)
					scene.yBuffer[i] = heightOnScreen
				}
			}
			px += dx
		}
	}
}

func main() {
	scene := NewScene("./maps/D1.png", "./maps/C1W.png")
	scene.LoadSetup()
	rl.InitWindow(SCREEN_WIDHT, SCREEN_HEIGHT, "voxelspace")

	rl.SetTargetFPS(60)
	renderTexture := rl.LoadRenderTexture(700, 700)
	defer rl.UnloadRenderTexture(renderTexture)

	for !rl.WindowShouldClose() {
		scene.updateCamera()
		rl.BeginTextureMode(renderTexture)
		rl.ClearBackground(rl.SkyBlue)

		scene.render()

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
