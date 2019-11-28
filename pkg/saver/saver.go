package saver

import (
	"errors"
	"github.com/ueef/mosaic/pkg/parse"
)

const TypeNull = "null"
const TypeDirect = "direct"
const TypeHashed = "hashed"

type Saver interface {
	Save(path string, data []byte) error
}

func New(t string, m map[string]interface{}) (s Saver, err error) {
	switch t {
	case TypeNull:
		return NewNullFromMap(m)
	case TypeDirect:
		return NewDirectFromMap(m)
	case TypeHashed:
		return NewHashedFromMap(m)
	}

	return nil, errors.New("type of Saver \"" + t + "\" is undefined")
}

func NewFromConfig(c interface{}) (s Saver, err error) {
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
