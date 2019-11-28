package encoder

import (
	"github.com/ueef/mosaic/pkg/parse"
	"image"
	"image/jpeg"
)

type JpegEncoder struct {
	Quality int
}

func (e JpegEncoder) Encode(img image.Image) ([]byte, error) {
	b := buffer{}

	err := jpeg.Encode(&b, img, &jpeg.Options{
		Quality: e.Quality,
	})
	if err != nil {
		return nil, nil
	}

	return b, nil
}

func (e JpegEncoder) GetMime() string {
	return "image/jpeg"
}

func NewJpegEncoder(quality int) *JpegEncoder {
	return &JpegEncoder{
		Quality: quality,
	}
}

func NewJpegEncoderFromMap(m map[string]interface{}) (*JpegEncoder, error) {
	quality, err := parse.GetRequiredIntFromMap("quality", m)
	if err != nil {
		return nil, err
	}

	return NewJpegEncoder(quality), nil
}
