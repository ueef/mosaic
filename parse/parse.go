package parse

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var parsedFonts = map[string]*truetype.Font{}

func GetMapFromMap(k string, m map[string]interface{}) (map[string]interface{}, bool, error) {
	o, ok := m[k]
	if !ok {
		return nil, false, nil
	}

	v, ok := o.(map[string]interface{})
	if !ok {
		return nil, true, fmt.Errorf("a value of a key \"%s\" must be a map[string]interface{}, got %T in map %v", k, o, m)
	}

	return v, true, nil
}

func GetRequiredMapFromMap(k string, m map[string]interface{}) (map[string]interface{}, error) {
	v, ok, err := GetMapFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetIntFromMap(k string, m map[string]interface{}) (int, bool, error) {
	o, ok := m[k]
	if !ok {
		return 0, false, nil
	}

	switch v := o.(type) {
	case int:
		return v, true, nil
	case float64:
		return int(v), true, nil
	}

	return 0, true, fmt.Errorf("a value of a key \"%s\" must be an integer, got %T in map %v", k, o, m)
}

func GetRequiredIntFromMap(k string, m map[string]interface{}) (int, error) {
	v, ok, err := GetIntFromMap(k, m)
	if !ok {
		return 0, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetFloatFromMap(k string, m map[string]interface{}) (float64, bool, error) {
	o, ok := m[k]
	if !ok {
		return 0, false, nil
	}

	switch v := o.(type) {
	case int:
		return float64(v), true, nil
	case float64:
		return v, true, nil
	}

	return 0, true, fmt.Errorf("a value of a key \"%s\" must be an float, got %T in map %v", k, o, m)
}

func GetRequiredFloatFromMap(k string, m map[string]interface{}) (float64, error) {
	v, ok, err := GetFloatFromMap(k, m)
	if !ok {
		return 0, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetStringFromMap(k string, m map[string]interface{}) (string, bool, error) {
	o, ok := m[k]
	if !ok {
		return "", false, nil
	}

	v, ok := o.(string)
	if !ok {
		return "", true, fmt.Errorf("a value of a key \"%s\" must be a string, got %T in map %v", k, o, m)
	}

	return v, true, nil
}

func GetRequiredStringFromMap(k string, m map[string]interface{}) (string, error) {
	v, ok, err := GetStringFromMap(k, m)
	if !ok {
		return "", fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetColorFromMap(k string, m map[string]interface{}) (color.Color, bool, error) {
	v, ok, err := GetStringFromMap(k, m)
	if !ok || err != nil {
		return nil, false, err
	}

	ok, err = regexp.MatchString("#[abcdef\\d]{4}", v)
	if err != nil {
		return nil, true, err
	}
	if ok {
		cv, err := strconv.ParseUint(v[1:], 16, 32)
		if err != nil {
			return nil, true, err
		}

		c := color.RGBA{
			R: uint8((15 & (cv >> 12)) * 17),
			G: uint8((15 & (cv >> 8)) * 17),
			B: uint8((15 & (cv >> 4)) * 17),
			A: uint8((15 & cv) * 17),
		}

		return &c, true, nil
	}

	ok, err = regexp.MatchString("#[abcdef\\d]{8}", v)
	if err != nil {
		return nil, true, err
	}
	if ok {
		cv, err := strconv.ParseUint(v[1:], 16, 32)
		if err != nil {
			panic(err)
		}

		c := color.RGBA{
			R: uint8(255 & (cv >> 24)),
			G: uint8(255 & (cv >> 16)),
			B: uint8(255 & (cv >> 8)),
			A: uint8(255 & cv),
		}

		return &c, true, nil
	}

	ok, err = regexp.MatchString("rgba\\(\\d{1,3},\\d{1,3},\\d{1,3},\\d{1,3}\\)", v)
	if err != nil {
		return nil, true, err
	}
	if ok {
		cv := [4]uint8{0,0,0,0}
		for i, s := range strings.Split(v[5:len(v)-1], ",") {
			v, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return nil, true, err
			}
			cv[i] = uint8(v)
		}

		c := color.RGBA{
			R: cv[0],
			G: cv[1],
			B: cv[2],
			A: cv[3],
		}

		return &c, true, nil
	}

	return nil, true, fmt.Errorf("a value of a key \"%s\" have unsupported format. expected format is #ffff, #ffffffff or rgba(255,255,255,255)", k)
}

func GetRequiredColorFromMap(k string, m map[string]interface{}) (color.Color, error) {
	v, ok, err := GetColorFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetFontFromMap(k string, m map[string]interface{}) (*truetype.Font, bool, error) {
	v, ok, err := GetStringFromMap(k, m)
	if !ok || err != nil {
		return nil, false, err
	}

	f, ok := parsedFonts[v]
	if ok {
		return f, true, nil
	}

	b, err := ioutil.ReadFile(v)
	if err != nil {
		return nil, true, err
	}

	f, err = freetype.ParseFont(b)
	if err != nil {
		return nil, true, err
	}

	parsedFonts[v] = f

	return f, true, nil
}

func GetRequiredFontFromMap(k string, m map[string]interface{}) (*truetype.Font, error) {
	v, ok, err := GetFontFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetFontHintingFromMap(k string, m map[string]interface{}) (font.Hinting, bool, error) {
	v, ok, err := GetStringFromMap(k, m)
	if !ok || err != nil {
		return font.HintingNone, false, err
	}

	switch v {
	case "none":
		return font.HintingNone, true, nil
	case "full":
		return font.HintingFull, true, nil
	case "vertical":
		return font.HintingVertical, true, nil
	}

	return font.HintingNone, true, fmt.Errorf("a value of a key \"%s\" is invalid. expected none, full or vertical, got %v", k, v)
}

func GetRequiredFontHintingFromMap(k string, m map[string]interface{}) (font.Hinting, error) {
	v, ok, err := GetFontHintingFromMap(k, m)
	if !ok {
		return font.HintingNone, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetRegexpFromMap(k string, m map[string]interface{}) (*regexp.Regexp, bool, error) {
	v, ok, err := GetStringFromMap(k, m)
	if !ok || err != nil {
		return nil, false, err
	}

	r, err := regexp.Compile(v)
	if err != nil {
		return nil, true, err
	}

	return r, true, nil
}

func GetRequiredRegexpFromMap(k string, m map[string]interface{}) (*regexp.Regexp, error) {
	v, ok, err := GetRegexpFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetInterfaceFromMap(k string, m map[string]interface{}) (interface{}, bool, error) {
	o, ok := m[k]
	if !ok {
		return nil, false, nil
	}

	return o, true, nil
}

func GetRequiredInterfaceFromMap(k string, m map[string]interface{}) (interface{}, error) {
	v, ok, err := GetInterfaceFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}

func GetSliceOfInterfacesFromMap(k string, m map[string]interface{}) ([]interface{}, bool, error) {
	o, ok := m[k]
	if !ok {
		return nil, false, nil
	}

	v, ok := o.([]interface{})
	if !ok {
		return nil, true, fmt.Errorf("a value of a key \"%s\" must be a slice of interfaces, got %T in map %v", k, o, m)
	}

	return v, true, nil
}

func GetRequiredSliceOfInterfacesFromMap(k string, m map[string]interface{}) ([]interface{}, error) {
	v, ok, err := GetSliceOfInterfacesFromMap(k, m)
	if !ok {
		return nil, fmt.Errorf("a key \"%s\" is undefined in map %v", k, m)
	}

	return v, err
}
