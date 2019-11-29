package dispatcher

import (
	"bytes"
	"fmt"
	"image"
)

func load(r *Response) *Response {
	b, err := r.Pict.Loader.Load(r.Path)
	if err != nil {
		return NewErrorResponse(r.Path, err)
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
		return NewErrorResponse(r.Path, err)
	}
	r.Buff = nil

	img, err = r.Pict.Filter.Apply(img)
	if err != nil {
		return NewErrorResponse(r.Path, err)
	}

	r.Buff, err = r.Pict.Encoder.Encode(img)
	if err != nil {
		return NewErrorResponse(r.Path, err)
	}

	return r
}
