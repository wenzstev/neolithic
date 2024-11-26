package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"math/rand"
	"os"
)

const (
	grass1Path = "assets/grass_1.png"
	grass2Path = "assets/grass_2.png"
	grass3Path = "assets/grass_3.png"
)

type SpriteGround struct {
	Images []*ebiten.Image
	Index  int
}

func NewGrassGround() (*SpriteGround, error) {
	return NewSpriteGround([]string{grass1Path, grass2Path, grass3Path})
}

var _ Ground = &SpriteGround{}

func NewSpriteGround(paths []string) (*SpriteGround, error) {
	images := make([]*ebiten.Image, len(paths))

	for i, path := range paths {
		img, err := loadSprite(path)
		if err != nil {
			return nil, err
		}

		images[i] = ebiten.NewImageFromImage(img)
	}

	ind := rand.Intn(len(images))

	return &SpriteGround{
		Images: images,
		Index:  ind,
	}, nil
}

func (s *SpriteGround) Draw(screen *ebiten.Image, transform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *transform
	screen.DrawImage(s.Images[s.Index], op)
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
