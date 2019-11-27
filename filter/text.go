package filter

import (
	"fmt"
	"github.com/anthonynsimon/bild/clone"
	"github.com/ueef/mosaic/parse"
	"github.com/ueef/mosaic/stamp"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
)

const TypeText = "text"

type Text struct {
	Gravity string
	Stamp stamp.Stamp
	TextColor color.Color
	BackgroundColor color.Color
}

func (f *Text) Apply(img image.Image) (image.Image, error) {
	p := 4
	w, h := img.Bounds().Dx(), f.Stamp.GetHeight() + p*2

	si := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: w,
			Y: h,
		},
	})

	draw.Draw(si, si.Bounds(), image.NewUniform(f.BackgroundColor), image.Point{}, draw.Src)
	f.Stamp.Draw((w-f.Stamp.GetWidth())/2, p, f.TextColor, si)

	ni, ok := img.(*image.RGBA)
	if !ok {
		ni = clone.AsRGBA(img)
	}

	buf := make([]uint8, len(ni.Pix)+len(si.Pix))
	switch f.Gravity {
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

func NewText(g string, s stamp.Stamp, tc color.Color, bc color.Color) *Text {
	return &Text{
		Gravity: g,
		Stamp: s,
		TextColor: tc,
		BackgroundColor: bc,
	}
}

func NewTextFromMap(m map[string]interface{}) (Filter, error) {
	g, err := parse.GetRequiredStringFromMap("gravity", m)
	if err != nil {
		return nil, err
	}

	if g != GravityNorth && g != GravitySouth {
		return nil, fmt.Errorf("a value of gravity only must be \"%s\" or \"%s\"", GravityNorth, GravitySouth)
	}

	f, err := parse.GetRequiredFontFromMap("font", m)
	if err != nil {
		return nil, err
	}

	fs, err := parse.GetRequiredFloatFromMap("font_size", m)
	if err != nil {
		return nil, err
	}

	t, err := parse.GetRequiredStringFromMap("text", m)
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

	return NewText(g, stamp.New(fs, t, font.HintingFull, f), tc, bc), nil
}

func init() {
	RegisterFilter(TypeText, NewTextFromMap)
}