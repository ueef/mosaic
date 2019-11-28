package stamp

import (
	"errors"
	"github.com/golang/freetype/truetype"
	"github.com/ueef/mosaic/pkg/parse"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
)

type Stamp interface {
	Draw(x, y int, c color.Color, i draw.Image)
	GetWidth() int
	GetHeight() int
}

type stamp struct {
	f font.Face
	t string
	w int
	h int
	x int
	y int
}

func (s *stamp) Draw(x, y int, c color.Color, i draw.Image) {
	d := font.Drawer{
		Dst:  i,
		Src:  image.NewUniform(c),
		Face: s.f,
		Dot: fixed.Point26_6{
			X: fixed.I(x + s.x),
			Y: fixed.I(y + s.y),
		},
	}
	d.DrawString(s.t)
}

func (s *stamp) GetWidth() int {
	return s.w
}

func (s *stamp) GetHeight() int {
	return s.h
}

func New(s float64, t string, h font.Hinting, f *truetype.Font) Stamp {
	ff := truetype.NewFace(f, &truetype.Options{
		Size:    s,
		Hinting: h,
	})

	b, _ := font.BoundString(ff, t)

	return &stamp{
		f: ff,
		t: t,
		w: b.Max.X.Ceil() - b.Min.X.Ceil(),
		h: b.Max.Y.Ceil() - b.Min.Y.Ceil(),
		x: 0,
		y: -b.Min.Y.Ceil(),
	}
}

func NewFromMap(m map[string]interface{}) (Stamp, error) {
	t, err := parse.GetRequiredStringFromMap("text", m)
	if err != nil {
		return nil, err
	}

	f, err := parse.GetRequiredFontFromMap("font", m)
	if err != nil {
		return nil, err
	}

	h, ok, err := parse.GetFontHintingFromMap("hinting", m)
	if err != nil {
		return nil, err
	}
	if !ok {
		h = font.HintingFull
	}

	fs, err := parse.GetRequiredFloatFromMap("font_size", m)
	if err != nil {
		return nil, err
	}

	return New(fs, t, h, f), nil
}

func NewFromConfig(c interface{}) (s Stamp, err error) {
	m, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must be a map")
	}

	return NewFromMap(m)
}
