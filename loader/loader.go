package loader

import (
	"errors"
	"github.com/ueef/mosaic/parse"
)

const TypeHttp = "http"
const TypeDirect = "direct"

type Loader interface {
	Load(path string) ([]byte, error)
}

func New(t string, m map[string]interface{}) (s Loader, err error) {
	switch t {
	case TypeHttp:
		return NewHttpFromMap(m)
	case TypeDirect:
		return NewDirectFromMap(m)
	}

	return nil, errors.New("type of loader \"" + t + "\" is undefined")
}

func NewFromConfig(c interface{}) (s Loader, err error) {
	m, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must be of the type map[string]interface{}")
	}

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
