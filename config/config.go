package config

import (
	"encoding/json"
	"errors"
	"github.com/ueef/mosaic/parse"
	"github.com/ueef/mosaic/picture"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Pictures picture.Pictures
}

func New(pc picture.Pictures) *Config {
	return &Config{
		Pictures: pc,
	}
}

func NewFromMap(m map[string]interface{}) (*Config, error) {
	i, err := parse.GetRequiredInterfaceFromMap("pictures", m)
	if err != nil {
		return nil, err
	}
	pc, err := picture.NewPicturesFromConfig(i)
	if err != nil {
		return nil, err
	}

	return New(*pc), nil
}

func ParseFile(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(name) {
	case ".json":
		return parseJson(buf)
	case ".yml":
		fallthrough
	case ".yaml":
		return parseYaml(buf)
	default:
		return nil, errors.New(name + " has an unsupported type")
	}
}

func parseYaml(d []byte) (*Config, error) {
	c := make(map[string]interface{})
	err := yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return NewFromMap(c)
}

func parseJson(d []byte) (*Config, error) {
	c := make(map[string]interface{})
	err := json.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return NewFromMap(c)
}
