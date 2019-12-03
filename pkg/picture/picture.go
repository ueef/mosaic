package picture

import (
	"errors"
	"github.com/ueef/mosaic/pkg/encoder"
	"github.com/ueef/mosaic/pkg/filter"
	"github.com/ueef/mosaic/pkg/loader"
	"github.com/ueef/mosaic/pkg/parse"
	"github.com/ueef/mosaic/pkg/saver"
	"regexp"
)

type Picture struct {
	Saver       saver.Saver
	Loader      loader.Loader
	Filter      filter.Filter
	Encoder     encoder.Encoder
	HostPattern *regexp.Regexp
	PathPattern *regexp.Regexp
}

func (p Picture) Match(host, path string) bool {
	return (p.HostPattern == nil || p.HostPattern.MatchString(host)) && (p.PathPattern == nil || p.PathPattern.MatchString(path))
}

type Pictures []*Picture

func (p Pictures) Match(host, path string) (*Picture, error) {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Match(host, path) {
			return p[i], nil
		}
	}

	return nil, errors.New("there aren't any matching pictures")
}

func New(saver saver.Saver, loader loader.Loader, filter filter.Filter, encoder encoder.Encoder, hostPattern *regexp.Regexp, pathPattern *regexp.Regexp) *Picture {
	return &Picture{
		Saver:       saver,
		Loader:      loader,
		Filter:      filter,
		Encoder:     encoder,
		HostPattern: hostPattern,
		PathPattern: pathPattern,
	}
}

func NewPictureFromConfig(c interface{}) (*Picture, error) {
	m, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must be of the type map[string]interface{}")
	}

	i, err := parse.GetRequiredInterfaceFromMap("saver", m)
	if err != nil {
		return nil, err
	}
	s, err := saver.NewFromConfig(i)
	if err != nil {
		return nil, err
	}

	i, err = parse.GetRequiredInterfaceFromMap("loader", m)
	if err != nil {
		return nil, err
	}
	l, err := loader.NewFromConfig(i)
	if err != nil {
		return nil, err
	}

	i, err = parse.GetRequiredInterfaceFromMap("filter", m)
	if err != nil {
		return nil, err
	}
	f, err := filter.NewFromConfig(i)
	if err != nil {
		return nil, err
	}

	i, err = parse.GetRequiredInterfaceFromMap("encoder", m)
	if err != nil {
		return nil, err
	}
	e, err := encoder.NewFromConfig(i)
	if err != nil {
		return nil, err
	}

	h, _, err := parse.GetRegexpFromMap("host_pattern", m)
	if err != nil {
		return nil, err
	}
	p, _, err := parse.GetRegexpFromMap("path_pattern", m)
	if err != nil {
		return nil, err
	}

	return New(s, l, f, e, h, p), nil
}

func NewPicturesFromConfig(c []interface{}) (Pictures, error) {
	var err error

	p := make(Pictures, len(c))
	for i := range c {
		p[i], err = NewPictureFromConfig(c[i])
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}
