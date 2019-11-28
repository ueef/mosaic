package filter

import (
	"errors"
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
	x, y, w, h := 0, 0, img.Bounds().Max.X, img.Bounds().Max.Y
	f := float32(w) / float32(h)

	rw, rh := filter.w, int(float32(filter.w)/f)
	if rw < filter.w || rh < filter.h {
		rw, rh = int(float32(filter.h)*f), filter.h
	}

	img = transform.Resize(img, rw, rh, transform.Linear)

	switch filter.g {
	case GravityEast:
		x = rw - filter.w
		y = (rh - filter.h) / 2
	case GravityWest:
		x = 0
		y = (rh - filter.h) / 2
	case GravitySouth:
		x = (rw - filter.w) / 2
		y = rh - filter.h
	case GravitySouthEast:
		x = rw - filter.w
		y = rh - filter.h
	case GravitySouthWest:
		x = 0
		y = rh - filter.h
	case GravityNorth:
		x = (rw - filter.w) / 2
		y = 0
	case GravityNorthEast:
		x = rw - filter.w
		y = 0
	case GravityNorthWest:
		x = 0
		y = 0
	case GravityCenter:
		x = (rw - filter.w) / 2
		y = (rh - filter.h) / 2
	default:
		return nil, errors.New("unexpected value of g")
	}

	img = transform.Crop(img, image.Rect(x, y, x+filter.w, y+filter.h))

	return img, nil
}

func NewThumbnail(w, h int, g string) Filter {
	return &thumbnail{w, h, g}
}

func NewThumbnailFromMap(m map[string]interface{}) (Filter, error) {
	width, err := parse.GetRequiredIntFromMap("w", m)
	if err != nil {
		return nil, err
	}

	height, err := parse.GetRequiredIntFromMap("h", m)
	if err != nil {
		return nil, err
	}

	gravity, err := parse.GetRequiredStringFromMap("g", m)
	if err != nil {
		return nil, err
	}

	return NewThumbnail(width, height, gravity), nil
}

func init() {
	RegisterFilter(TypeThumbnail, NewThumbnailFromMap)
}
