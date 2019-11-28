package loader

import (
	"errors"
	"github.com/ueef/mosaic/pkg/parse"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Http struct {
	host    string
	scheme  string
	replace string
	pattern *regexp.Regexp
}

func (s Http) Load(path string) ([]byte, error) {
	if s.pattern != nil {
		path = s.pattern.ReplaceAllString(path, s.replace)
	}

	r, err := http.Get(s.scheme + "://" + s.host + path)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}

	return ioutil.ReadAll(r.Body)
}

func NewHttp(host, scheme, replace string, pattern *regexp.Regexp) *Http {
	return &Http{
		host:    host,
		scheme:  scheme,
		replace: replace,
		pattern: pattern,
	}
}

func NewHttpFromMap(m map[string]interface{}) (*Http, error) {
	h, err := parse.GetRequiredStringFromMap("host", m)
	if err != nil {
		return nil, err
	}

	s, err := parse.GetRequiredStringFromMap("scheme", m)
	if err != nil {
		return nil, err
	}

	r, _, err := parse.GetStringFromMap("replace", m)
	if err != nil {
		return nil, err
	}

	p, _, err := parse.GetRegexpFromMap("pattern", m)
	if err != nil {
		return nil, err
	}

	return NewHttp(h, s, r, p), nil
}
