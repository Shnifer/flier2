package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"os"
	"image"
	_"image/png"
	_"image/jpeg"
	_"image/gif"
	"time"
	"math/rand"
	"math"
	"log"
)

func run(){
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel test!",
		Bounds: pixel.R(0, 0, 1024, 768),
		//VSync:true,

	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	//win.SetSmooth(true)

	spritesheet, err := loadPicture("textures/trees.png")
	if err != nil {
		panic(err)
	}
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	var treesFrames[]pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			treesFrames = append(treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos   = pixel.ZV
		camZoom      = 1.0
		trees    []*pixel.Sprite
		matrices []pixel.Matrix
	)

	const (
		camSpeed=200
		camZoomSpeed = 1.2
	)

	ticker:=time.Tick(time.Second)

	last:=time.Now()
	lastFN,FN:=0,0
	for !win.Closed() {
		FN++
		dt:=time.Since(last).Seconds()
		last=time.Now()

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		cam := pixel.IM.Scaled(pixel.ZV,camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)


		if win.Pressed(pixelgl.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			trees = append(trees, tree)
			mouse:=cam.Unproject(win.MousePosition())
			matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		win.Clear(colornames.Forestgreen)
		batch.Clear()
		for i:=range trees{
			trees[i].Draw(batch,matrices[i])
		}
		batch.MakePicture()
		batch.Draw(win)

		win.Update()
		select{
			case <-ticker:
				log.Println("FPS",FN-lastFN)
				lastFN = FN
			default:
		}
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main(){
	pixelgl.Run(run)
}