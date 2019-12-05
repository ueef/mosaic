package dispatcher

import (
	"github.com/ueef/mosaic/pkg/picture"
)

type Response struct {
	Err    error
	Buff   []byte
	Path   string
	Pict   *picture.Picture
	Timing Timer
}

func (r Response) IsSuccessful() bool {
	return r.Err == nil
}

func NewResponse(path string, pict *picture.Picture) *Response {
	return &Response{
		Err:    nil,
		Buff:   nil,
		Path:   path,
		Pict:   pict,
		Timing: NewTimer(),
	}
}

func NewErrorResponse(path string, err error, timing Timer) *Response {
	return &Response{
		Err:    err,
		Buff:   nil,
		Path:   path,
		Pict:   nil,
		Timing: timing,
	}
}
