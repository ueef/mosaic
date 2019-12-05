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
	Filters     []filter.Filter
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

func New(saver saver.Saver, loader loader.Loader, filters []filter.Filter, encoder encoder.Encoder, hostPattern *regexp.Regexp, pathPattern *regexp.Regexp) *Picture {
	return &Picture{
		Saver:       saver,
		Loader:      loader,
		Filters:     filters,
		Encoder:     encoder,
		HostPattern: hostPattern,
		PathPattern: pathPattern,
	}
}

func NewPictureFromConfig(c interface{}) (*Picture, error) {
	mv, ok := c.(map[string]interface{})
	if !ok {
		return nil, errors.New("a config must be of the type map[string]interface{}")
	}

	iv, err := parse.GetRequiredInterfaceFromMap("saver", mv)
	if err != nil {
		return nil, err
	}
	s, err := saver.NewFromConfig(iv)
	if err != nil {
		return nil, err
	}

	iv, err = parse.GetRequiredInterfaceFromMap("loader", mv)
	if err != nil {
		return nil, err
	}
	l, err := loader.NewFromConfig(iv)
	if err != nil {
		return nil, err
	}

	iv, err = parse.GetRequiredInterfaceFromMap("filter", mv)
	if err != nil {
		return nil, err
	}

	sv, ok := iv.([]interface{})
	if !ok {
		sv = []interface{}{iv}
	}
	f := make([]filter.Filter, len(sv))
	for i, iv := range sv {
		f[i], err = filter.NewFromConfig(iv)
		if err != nil {
			return nil, err
		}
	}

	iv, err = parse.GetRequiredInterfaceFromMap("encoder", mv)
	if err != nil {
		return nil, err
	}
	e, err := encoder.NewFromConfig(iv)
	if err != nil {
		return nil, err
	}

	h, _, err := parse.GetRegexpFromMap("host_pattern", mv)
	if err != nil {
		return nil, err
	}
	p, _, err := parse.GetRegexpFromMap("path_pattern", mv)
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
