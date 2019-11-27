package loader

import (
	"github.com/ueef/mosaic/parse"
	"io/ioutil"
	"os"
)

type Direct struct {
	Dir string
}

func (s Direct) Load(path string) ([]byte, error) {
	f, err := os.Open(s.Dir + "/" + path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func NewDirect(dir string) *Direct {
	return &Direct{
		Dir: dir,
	}
}

func NewDirectFromMap(m map[string]interface{}) (*Direct, error) {
	dir, err := parse.GetRequiredStringFromMap("dir", m)
	if err != nil {
		return nil, err
	}

	return NewDirect(dir), nil
}
