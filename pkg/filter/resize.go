package filter

import (
	"github.com/anthonynsimon/bild/transform"
	"github.com/ueef/mosaic/pkg/parse"
	"image"
)

const TypeResize = "resize"

type resize struct {
	w int
	h int
}

func (f resize) Apply(img image.Image) (image.Image, error) {
	if f.w == 0 && f.h == 0 {
		return img, nil
	}

	b := img.Bounds()
	if (f.w == 0 || b.Dx() < f.w) && (f.h == 0 || b.Dy() < f.h) {
		return img, nil
	}

	fr := float32(b.Dx()) / float32(b.Dy())

	var w, h int
	if f.w == 0 {
		w, h = int(float32(f.h)*fr), f.h
	} else if f.h == 0 {
		w, h = f.w, int(float32(f.w)/fr)
	} else {
		w, h = f.w, int(float32(f.w)/fr)
		if w > f.w || h > f.h {
			w, h = int(float32(f.h)*fr), f.h
		}
	}

	img = transform.Resize(img, w, h, transform.Linear)

	return img, nil
}

func NewResize(w, h int) Filter {
	return &resize{w, h}
}

func NewResizeFromMap(m map[string]interface{}) (Filter, error) {
	width, _, err := parse.GetIntFromMap("w", m)
	if err != nil {
		return nil, err
	}

	height, _, err := parse.GetIntFromMap("h", m)
	if err != nil {
		return nil, err
	}

	return NewResize(width, height), nil
}

func init() {
	RegisterFilter(TypeResize, NewResizeFromMap)
}
