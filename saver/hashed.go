package saver

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/ueef/mosaic/parse"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Hashed struct {
	Dir string
}

func (s Hashed) Save(path string, data []byte) error {
	p, err := s.GetFilePath(path)
	if err != nil {
		return err
	}
	d := filepath.Dir(p)

	err = os.MkdirAll(d, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (s Hashed) GetFilePath(path string) (string, error) {
	hr := md5.New()
	_, err := hr.Write([]byte(path))
	if err != nil {
		return "", err
	}
	hs := base64.RawURLEncoding.EncodeToString(hr.Sum(nil))

	return s.Dir + "/" + hs + ".png", nil
}

func NewHashed(dir string) *Hashed {
	return &Hashed{
		Dir: dir,
	}
}

func NewHashedFromMap(m map[string]interface{}) (*Hashed, error) {
	dir, err := parse.GetRequiredStringFromMap("dir", m)
	if err != nil {
		return nil, err
	}

	return NewHashed(dir), nil
}
