package asset

import (
	"encoding/json"
	"errors"
	"github.com/faiface/pixel"
	"github.com/zcubbs/zworld/pkg/packer"
	"image"
	_ "image/png"
	"io"
	"io/fs"
)

type Load struct {
	filesystem fs.FS
}

func NewLoad(filesystem fs.FS) *Load {
	return &Load{filesystem: filesystem}
}

func (load *Load) Open(path string) (fs.File, error) {
	return load.filesystem.Open(path)
}

func (load *Load) Image(path string) (image.Image, error) {
	file, err := load.filesystem.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (load *Load) Sprite(path string) (*pixel.Sprite, error) {
	img, err := load.Image(path)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)

	return pixel.NewSprite(pic, pic.Bounds()), nil
}

func (load *Load) Json(path string, data interface{}) error {
	file, err := load.filesystem.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	return nil
}

func (load *Load) SpriteSheet(path string) (*SpriteSheet, error) {
	// load JSON
	serializedSpreadSheet := packer.SerializedSpritesheet{}
	err := load.Json(path, &serializedSpreadSheet)
	if err != nil {
		return nil, err
	}

	// load the image
	img, err := load.Image(serializedSpreadSheet.ImageName)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)

	// create the spritesheet object
	bounds := pic.Bounds()
	lookup := make(map[string]*pixel.Sprite)
	for k, v := range serializedSpreadSheet.Frames {
		rect := pixel.R(
			v.Frame.X,
			bounds.H()-v.Frame.Y,
			v.Frame.X+v.Frame.W,
			bounds.H()-(v.Frame.Y+v.Frame.H),
		).Norm()

		lookup[k] = pixel.NewSprite(pic, rect)
	}

	return NewSpriteSheet(pic, lookup), nil
}

type SpriteSheet struct {
	picture pixel.Picture
	lookup  map[string]*pixel.Sprite
}

func NewSpriteSheet(pic pixel.Picture, lookup map[string]*pixel.Sprite) *SpriteSheet {
	return &SpriteSheet{
		picture: pic,
		lookup:  lookup,
	}
}

func (s *SpriteSheet) Get(name string) (*pixel.Sprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New("invalid sprite name")
	}

	return sprite, nil
}

func (s *SpriteSheet) Picture() pixel.Picture {
	return s.picture
}
