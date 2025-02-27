package bilicomicdownloader

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	_ "github.com/gen2brain/avif"
)

func ImgToPng(imageData []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	if format == "png" {
		return imageData, nil
	}

	var buf bytes.Buffer

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ImgToJpg(imageData []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	if format == "jpeg" {
		return imageData, nil
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
