package filter

import (
	"errors"
	"github.com/anthonynsimon/bild/transform"
	"github.com/ueef/mosaic/parse"
	"image"
)

const TypeThumbnail = "thumbnail"

type Thumbnail struct {
	width   int
	height  int
	gravity string
}

func (filter Thumbnail) Apply(img image.Image) (image.Image, error) {
	x, y, w, h := 0, 0, img.Bounds().Max.X, img.Bounds().Max.Y
	f := float32(w) / float32(h)

	rw, rh := filter.width, int(float32(filter.width)/f)
	if rw < filter.width || rh < filter.height {
		rw, rh = int(float32(filter.height)*f), filter.height
	}

	img = transform.Resize(img, rw, rh, transform.Linear)

	switch filter.gravity {
	case GravityEast:
		x = rw - filter.width
		y = (rh - filter.height) / 2
	case GravityWest:
		x = 0
		y = (rh - filter.height) / 2
	case GravitySouth:
		x = (rw - filter.width) / 2
		y = rh - filter.height
	case GravityNorth:
		x = (rw - filter.width) / 2
		y = 0
	case GravityCenter:
		x = (rw - filter.width) / 2
		y = (rh - filter.height) / 2
	default:
		return nil, errors.New("unexpected value of gravity")
	}

	img = transform.Crop(img, image.Rect(x, y, x+filter.width, y+filter.height))

	return img, nil
}

func NewThumbnail(w, h int, g string) *Thumbnail {
	return &Thumbnail{w, h, g}
}

func NewThumbnailFromMap(m map[string]interface{}) (Filter, error) {
	width, err := parse.GetRequiredIntFromMap("width", m)
	if err != nil {
		return nil, err
	}

	height, err := parse.GetRequiredIntFromMap("height", m)
	if err != nil {
		return nil, err
	}

	gravity, err := parse.GetRequiredStringFromMap("gravity", m)
	if err != nil {
		return nil, err
	}

	return NewThumbnail(width, height, gravity), nil
}

func init() {
	RegisterFilter(TypeThumbnail, NewThumbnailFromMap)
}
