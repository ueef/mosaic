package dispatcher

import (
	"github.com/ueef/mosaic/pkg/picture"
	"time"
)

type Response struct {
	Err   error
	Buff  []byte
	Path  string
	Pict  *picture.Picture
	Times map[string]time.Duration
}

func (r Response) IsSuccessful() bool {
	return r.Err == nil
}

func NewResponse(path string, pict *picture.Picture) *Response {
	return &Response{
		Err:   nil,
		Buff:  nil,
		Path:  path,
		Pict:  pict,
		Times: map[string]time.Duration{},
	}
}

func NewErrorResponse(path string, err error) *Response {
	return &Response{
		Err:   err,
		Buff:  nil,
		Path:  path,
		Pict:  nil,
		Times: map[string]time.Duration{},
	}
}
