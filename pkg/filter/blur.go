package filter

import (
	blur2 "github.com/anthonynsimon/bild/blur"
	"github.com/ueef/mosaic/pkg/parse"
	"image"
)

const TypeBlur = "blur"

type blur struct {
	r float64
}

func (f blur) Apply(img image.Image) (image.Image, error) {
	img = blur2.Gaussian(img, f.r)

	return img, nil
}

func NewBlur(r float64) Filter {
	return &blur{r}
}

func NewBlurFromMap(m map[string]interface{}) (Filter, error) {
	radius, err := parse.GetRequiredFloatFromMap("radius", m)
	if err != nil {
		return nil, err
	}

	return NewBlur(radius), nil
}

func init() {
	RegisterFilter(TypeBlur, NewBlurFromMap)
}
