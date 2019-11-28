package filter

import (
	"image"
)

const TypeNull = "null"

type Null struct{}

func (f Null) Apply(img image.Image) (image.Image, error) {
	return img, nil
}

func NewNull() *Null {
	return &Null{}
}

func NewNullFromMap(m map[string]interface{}) (Filter, error) {
	return NewNull(), nil
}

func init() {
	RegisterFilter(TypeNull, NewNullFromMap)
}
