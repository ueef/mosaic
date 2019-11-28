package saver

import (
	"github.com/ueef/mosaic/pkg/parse"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Direct struct {
	Dir string
}

func (s Direct) Save(path string, data []byte) error {
	p := s.GetFilePath(path)
	d := filepath.Dir(p)

	err := os.MkdirAll(d, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (s Direct) GetFilePath(path string) string {
	return s.Dir + "/" + path
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
