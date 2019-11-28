package filter

import (
	"image"
)

const TypeNull = "null"

type null struct{}

func (f null) Apply(img image.Image) (image.Image, error) {
	return img, nil
}

func NewNull() Filter {
	return &null{}
}

func NewNullFromMap(m map[string]interface{}) (Filter, error) {
	return NewNull(), nil
}

func init() {
	RegisterFilter(TypeNull, NewNullFromMap)
}
