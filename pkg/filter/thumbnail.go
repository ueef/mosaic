package filter

import (
	"errors"
	"fmt"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/transform"
	"github.com/ueef/mosaic/pkg/parse"
	"image"
)

const TypeThumbnail = "thumbnail"

type thumbnail struct {
	w int
	h int
	g string
}

func (filter thumbnail) Apply(img image.Image) (image.Image, error) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = clone.AsRGBA(img)
	}

	f := float32(rgba.Bounds().Dx()) / float32(rgba.Bounds().Dy())

	w, h := filter.w, int(float32(filter.w)/f)
	if w < filter.w || h < filter.h {
		w, h = int(float32(filter.h)*f), filter.h
	}

	rgba = transform.Resize(rgba, w, h, transform.Linear)

	x, y := 0, 0
	switch filter.g {
	case GravityEast:
		x = w - filter.w
		y = (h - filter.h) / 2
	case GravityWest:
		x = 0
		y = (h - filter.h) / 2
	case GravitySouth:
		x = (w - filter.w) / 2
		y = h - filter.h
	case GravitySouthEast:
		x = w - filter.w
		y = h - filter.h
	case GravitySouthWest:
		x = 0
		y = h - filter.h
	case GravityNorth:
		x = (w - filter.w) / 2
		y = 0
	case GravityNorthEast:
		x = w - filter.w
		y = 0
	case GravityNorthWest:
		x = 0
		y = 0
	case GravityCenter:
		x = (w - filter.w) / 2
		y = (h - filter.h) / 2
	default:
		return nil, errors.New("unexpected value of g")
	}

	return crop(rgba, x, y, x+filter.w, y+filter.h)
}

func NewThumbnail(w, h int, g string) Filter {
	if g == "" {
		g = GravityCenter
	}

	return &thumbnail{w, h, g}
}

func NewThumbnailFromMap(m map[string]interface{}) (Filter, error) {
	w, err := parse.GetRequiredIntFromMap("width", m)
	if err != nil {
		return nil, err
	}

	h, err := parse.GetRequiredIntFromMap("height", m)
	if err != nil {
		return nil, err
	}

	g, _, err := parse.GetStringFromMap("gravity", m)
	if err != nil {
		return nil, err
	}

	return NewThumbnail(w, h, g), nil
}

func crop(src *image.RGBA, x0, y0, x1, y1 int) (*image.RGBA, error) {
	w, h := x1-x0, y1-y0
	if x1 > src.Bounds().Max.X || y1 > src.Bounds().Max.Y {
		return nil, fmt.Errorf("the cropping area is out of bounds")
	}

	dst := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: w,
			Y: h,
		},
	})

	for y := 0; y < h; y++ {
		a0, a1 := y*dst.Stride, y*dst.Stride+w*4
		b0, b1 := (y+y0)*src.Stride+x0*4, (y+y0)*src.Stride+x1*4
		copy(dst.Pix[a0:a1], src.Pix[b0:b1])
	}

	return dst, nil
}

func init() {
	RegisterFilter(TypeThumbnail, NewThumbnailFromMap)
}
