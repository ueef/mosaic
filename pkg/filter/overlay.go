package filter

import (
	"fmt"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/transform"
	"github.com/ueef/mosaic/pkg/parse"
	"image"
)

const TypeOverlay = "overlay"

type overlay struct {
	p  int
	g  string
	fi *image.RGBA
}

func (f *overlay) Apply(img image.Image) (image.Image, error) {
	si, ok := img.(*image.RGBA)
	if !ok {
		si = clone.AsRGBA(img)
	}

	sb := si.Bounds()
	sw, sh := sb.Dx(), sb.Dy()

	fi := f.fitForegroundImage(sw, sh)
	if fi == nil {
		return img, nil
	}

	fb := fi.Bounds()
	fw, fh := fb.Dx(), fb.Dy()

	sx, sy := f.p, f.p
	switch f.g {
	case GravityEast:
		sx = sw - fw - f.p
		sy = (sh-(fh+f.p*2))/2 + f.p
	case GravityWest:
		sx = f.p
		sy = (sh-(fh+f.p*2))/2 + f.p
	case GravitySouth:
		sx = (sw-(fw+f.p*2))/2 + f.p
		sy = sh - fh - f.p
	case GravitySouthEast:
		sx = sw - fw - f.p
		sy = sh - fh - f.p
	case GravitySouthWest:
		sx = f.p
		sy = sh - fh - f.p
	case GravityNorth:
		sx = (sw-(fw+f.p*2))/2 + f.p
		sy = f.p
	case GravityNorthEast:
		sx = sw - fw - f.p
		sy = f.p
	case GravityNorthWest:
		sx = f.p
		sy = f.p
	case GravityCenter:
		sx = (sw-(fw+f.p*2))/2 + f.p
		sy = (sh-(fh+f.p*2))/2 + f.p
	default:
		return nil, fmt.Errorf("a g \"%s\" is invalid", f.g)
	}

	for y := 0; y < fh; y++ {
		sp := (sy+y)*si.Stride + sx*4
		fp := y * fi.Stride
		for x := 0; x < fw; x++ {
			fa := int(fi.Pix[fp+3])
			for p := 0; p < 3; p++ {
				sc := int(si.Pix[sp+p])
				fc := int(fi.Pix[fp+p])

				si.Pix[sp+p] = uint8((fc*fa)/255 + (255-fa)*sc/255)
			}

			sp += 4
			fp += 4
		}
	}

	return si, nil
}

func (f *overlay) fitForegroundImage(sw, sh int) *image.RGBA {
	if sw < f.p*2 || sh < f.p*2 {
		return nil
	}

	fw, fh := f.fi.Bounds().Dx(), f.fi.Bounds().Dy()
	if fw+f.p*2 < sw && fh+f.p*2 < sh {
		return f.fi
	}

	ff := float32(fw) / float32(fh)
	w := sw - f.p*2
	h := int(float32(w) / ff)
	if w+f.p*2 > sw || h+f.p*2 > sh {
		h = sh - f.p*2
		w = int(float32(h) * ff)
	}

	return transform.Resize(f.fi, w, h, transform.Linear)
}

func NewOverlay(p int, g string, fi *image.RGBA) Filter {
	if g == "" {
		g = GravityCenter
	}

	return &overlay{
		p:  p,
		g:  g,
		fi: fi,
	}
}

func NewOverlayFromMap(m map[string]interface{}) (Filter, error) {
	p, _, err := parse.GetIntFromMap("padding", m)
	if err != nil {
		return nil, err
	}

	g, _, err := parse.GetStringFromMap("gravity", m)
	if err != nil {
		return nil, err
	}

	fi, err := parse.GetRequiredImageFromMap("image", m)
	if err != nil {
		return nil, err
	}

	return NewOverlay(p, g, fi), nil
}

func init() {
	RegisterFilter(TypeOverlay, NewOverlayFromMap)
}
