package config

import (
	"encoding/json"
	"errors"
	"github.com/ueef/mosaic/pkg/picture"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ParsePath(path string) (picture.Pictures, error) {
	pics := picture.Pictures{}
	paths, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		p, err := parseFile(path)
		if err != nil {
			return nil, err
		}

		pics = merge(pics, p)
	}

	return pics, nil
}

func ParsePaths(paths []string) (picture.Pictures, error) {
	pics := picture.Pictures{}
	for _, path := range paths {
		p, err := ParsePath(path)
		if err != nil {
			return nil, err
		}

		pics = merge(pics, p)
	}

	return pics, nil
}

func merge(a, b picture.Pictures) picture.Pictures {
	c := make(picture.Pictures, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)

	return c
}

func parseFile(path string) (picture.Pictures, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(path) {
	case ".json":
		return parseJson(buf)
	case ".yml":
		fallthrough
	case ".yaml":
		return parseYaml(buf)
	default:
		return nil, errors.New(path + " has an unsupported type")
	}
}

func parseJson(d []byte) (picture.Pictures, error) {
	c := make([]interface{}, 0)
	err := json.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return picture.NewPicturesFromConfig(c)
}

func parseYaml(d []byte) (picture.Pictures, error) {
	c := make([]interface{}, 0)
	err := yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return picture.NewPicturesFromConfig(c)
}
