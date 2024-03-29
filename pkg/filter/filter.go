package filter

import (
	"errors"
	"fmt"
	"github.com/ueef/mosaic/pkg/parse"
	"image"
)

const GravityEast string = "east"
const GravityWest string = "west"
const GravityNorth string = "north"
const GravityNorthEast string = "north_east"
const GravityNorthWest string = "north_west"
const GravitySouth string = "south"
const GravitySouthEast string = "south_east"
const GravitySouthWest string = "south_west"
const GravityCenter string = "center"

var registered = map[string]func(m map[string]interface{}) (Filter, error){}

type Filter interface {
	Apply(img image.Image) (image.Image, error)
}

type Filters []Filter

func (f Filters) Apply(img image.Image) (image.Image, error) {
	var err error
	for _, v := range f {
		img, err = v.Apply(img)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

func New(t string, m map[string]interface{}) (Filter, error) {
	c, ok := registered[t]
	if !ok {
		return nil, fmt.Errorf("type of filter \"%s\" is unregistered", t)
	}

	f, err := c(m)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func NewFromMap(m map[string]interface{}) (Filter, error) {
	t, err := parse.GetRequiredStringFromMap("type", m)
	if err != nil {
		return nil, err
	}

	m, err = parse.GetRequiredMapFromMap("config", m)
	if err != nil {
		return nil, err
	}

	return New(t, m)
}

func NewFromConfig(c interface{}) (Filter, error) {
	m, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must contain a value of the type map[string]interface{}")
	}

	f, err := NewFromMap(m)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func RegisterFilter(t string, c func(m map[string]interface{}) (Filter, error)) {
	registered[t] = c
}
