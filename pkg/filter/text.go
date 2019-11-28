package filter

import (
	"fmt"
	"github.com/anthonynsimon/bild/clone"
	"github.com/ueef/mosaic/pkg/parse"
	"github.com/ueef/mosaic/pkg/stamp"
	"image"
	"image/color"
	"image/draw"
)

const TypeText = "text"

type text struct {
	g  string
	s  stamp.Stamp
	tc color.Color
	bc color.Color
}

func (f *text) Apply(img image.Image) (image.Image, error) {
	p := 4
	w, h := img.Bounds().Dx(), f.s.GetHeight()+p*2

	si := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: w,
			Y: h,
		},
	})

	draw.Draw(si, si.Bounds(), image.NewUniform(f.bc), image.Point{}, draw.Src)
	f.s.Draw((w-f.s.GetWidth())/2, p, f.tc, si)

	ni, ok := img.(*image.RGBA)
	if !ok {
		ni = clone.AsRGBA(img)
	}

	buf := make([]uint8, len(ni.Pix)+len(si.Pix))
	switch f.g {
	case GravityNorth:
		copy(buf, si.Pix)
		copy(buf[len(si.Pix):], ni.Pix)
	case GravitySouth:
		copy(buf, ni.Pix)
		copy(buf[len(ni.Pix):], si.Pix)
	}

	ni.Pix = buf
	ni.Rect.Max.Y += h

	return ni, nil
}

func NewText(g string, s stamp.Stamp, tc color.Color, bc color.Color) Filter {
	return &text{
		g:  g,
		s:  s,
		tc: tc,
		bc: bc,
	}
}

func NewTextFromMap(m map[string]interface{}) (Filter, error) {
	g, err := parse.GetRequiredStringFromMap("gravity", m)
	if err != nil {
		return nil, err
	}

	if g != GravityNorth && g != GravitySouth {
		return nil, fmt.Errorf("a value of g only must be \"%s\" or \"%s\"", GravityNorth, GravitySouth)
	}

	sm, err := parse.GetRequiredMapFromMap("stamp", m)
	if err != nil {
		return nil, err
	}
	s, err := stamp.NewFromMap(sm)
	if err != nil {
		return nil, err
	}

	tc, err := parse.GetRequiredColorFromMap("text_color", m)
	if err != nil {
		return nil, err
	}

	bc, err := parse.GetRequiredColorFromMap("background_color", m)
	if err != nil {
		return nil, err
	}

	return NewText(g, s, tc, bc), nil
}

func init() {
	RegisterFilter(TypeText, NewTextFromMap)
}
