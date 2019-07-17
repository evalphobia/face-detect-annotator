package engine

import (
	"image"
	_ "image/jpeg"
	"os"
)

func GetImageSize(imgPath string) (width, height int, err error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	c, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return c.Width, c.Height, nil
}
