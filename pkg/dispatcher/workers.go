package dispatcher

import (
	"bytes"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/ueef/mosaic/pkg/utils"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func load(r *Response) *Response {
	b, err := r.Pict.Loader.Load(r.Path)
	if err != nil {
		return NewErrorResponse(r.Path, err, r.Timing)
	}
	r.Buff = b

	return r
}

func save(r *Response) *Response {
	err := r.Pict.Saver.Save(r.Path, r.Buff)
	if err != nil {
		fmt.Println(err)
	}

	return r
}

func process(r *Response) *Response {
	img, _, err := image.Decode(bytes.NewReader(r.Buff))
	if err != nil {
		return NewErrorResponse(r.Path, err, r.Timing)
	}

	img = fixOrientation(img, r.Buff)
	r.Buff = nil

	img, err = r.Pict.Filter.Apply(img)
	if err != nil {
		return NewErrorResponse(r.Path, err, r.Timing)
	}

	r.Buff, err = r.Pict.Encoder.Encode(img)
	if err != nil {
		return NewErrorResponse(r.Path, err, r.Timing)
	}

	return r
}

func fixOrientation(i image.Image, b []byte) image.Image {
	e, err := exif.Decode(bytes.NewReader(b))
	if err != nil {
		return i
	}

	t, err := e.Get(exif.Orientation)
	if err != nil {
		return i
	}

	o, err := t.Int(0)
	if err != nil {
		return i
	}

	switch o {
	case 3:
		return utils.Rotate180(i)
	case 6:
		return utils.Rotate90(i)
	case 8:
		return utils.Rotate270(i)
	}

	return i
}
