package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"math/rand"
	"os"
)

const (
	grass1Path = "assets/grass_1.png"
	grass2Path = "assets/grass_2.png"
	grass3Path = "assets/grass_3.png"
)

// Ground represents the base terrain of a tile
type Ground struct {
	Image *ebiten.Image
}

// NewRGBGround creates a new ground with a single color.
func NewRGBGround(col color.RGBA, cellSize int) *Ground {
	image := ebiten.NewImage(cellSize, cellSize)
	image.Fill(col)
	return &Ground{
		Image: image,
	}
}

// NewGrassGround creates a new ground sprite with a Grass texture
func NewGrassGround() (*Ground, error) {
	return NewSpriteGround([]string{grass1Path, grass2Path, grass3Path})
}

func NewSpriteGround(paths []string) (*Ground, error) {
	images := make([]*ebiten.Image, len(paths))

	for i, path := range paths {
		img, err := loadSprite(path)
		if err != nil {
			return nil, err
		}

		images[i] = ebiten.NewImageFromImage(img)
	}

	ind := rand.Intn(len(images))

	return &Ground{
		Image: images[ind],
	}, nil
}

func loadSprite(filePath string) (*ebiten.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	sprite := ebiten.NewImageFromImage(img)
	return sprite, nil
}
