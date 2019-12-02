package loader

import (
	"github.com/ueef/mosaic/pkg/parse"
	"io/ioutil"
	"os"
	"regexp"
)

type Direct struct {
	d string
	r string
	p *regexp.Regexp
}

func (s Direct) Load(path string) ([]byte, error) {
	if nil != s.p {
		path = s.p.ReplaceAllString(path, s.r)
	}

	f, err := os.Open(s.d + path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func NewDirect(dir, replace string, pattern *regexp.Regexp) *Direct {
	return &Direct{
		d: dir,
		r: replace,
		p: pattern,
	}
}

func NewDirectFromMap(m map[string]interface{}) (*Direct, error) {
	d, err := parse.GetRequiredStringFromMap("dir", m)
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

	return NewDirect(d, r, p), nil
}
