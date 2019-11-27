package encoder

import (
	"errors"
	"github.com/ueef/mosaic/parse"
	"image"
)

const TypePng = "png"
const TypeJpeg = "jpeg"

type Encoder interface {
	Encode(img image.Image) ([]byte, error)
	GetMime() string
}

func New(t string, m map[string]interface{}) (Encoder, error) {
	switch t {
	case TypePng:
		return NewPngEncoderFromMap(m)
	case TypeJpeg:
		return NewJpegEncoderFromMap(m)
	}

	return nil, errors.New("type of encoder \"" + t + "\" is undefined")
}

func NewFromConfig(c interface{}) (s Encoder, err error) {
	m, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must be of the type map[string]interface{}")
	}

	t, err := parse.GetRequiredStringFromMap("type", m)
	if err != nil {
		return nil, err
	}

	m, err = parse.GetRequiredMapFromMap("config", m)
	if err != nil {
		return nil, err
	}

	return New(t, m)
}

type buffer []byte

func (b *buffer) Write(p []byte) (n int, err error) {
	c := make([]byte, len(*b)+len(p))
	bf := c[:len(*b)]
	copy(bf, *b)
	bf = c[len(*b):]
	copy(bf, p)
	*b = c

	return len(p), nil
}
