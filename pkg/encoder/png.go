package encoder

import (
	"image"
	"image/png"
)

type PngEncoder struct{}

func (e PngEncoder) Encode(img image.Image) ([]byte, error) {
	b := buffer{}

	err := png.Encode(&b, img)
	if err != nil {
		return nil, nil
	}

	return b, nil
}

func (e PngEncoder) GetMime() string {
	return "image/png"
}

func NewPngEncoder() *PngEncoder {
	return &PngEncoder{}
}

func NewPngEncoderFromMap(m map[string]interface{}) (*PngEncoder, error) {
	return NewPngEncoder(), nil
}
