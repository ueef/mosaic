package dispatcher

import (
	"github.com/ueef/mosaic/pkg/picture"
	"time"
)

type Dispatcher struct {
	cache    cache
	awaiters awaiters
	pictures picture.Pictures
	channels struct {
		saving     chan *Response
		loading    chan *Response
		processing chan *Response
		responding chan *Response
	}
}

func (d *Dispatcher) Dispatch(host, path string) (<-chan *Response, error) {
	pict, err := d.pictures.Match(host, path)
	if err != nil {
		return nil, err
	}

	c := make(chan *Response, 1)
	if d.awaiters.push(path, c) {
		if r := d.cache.get(path); r != nil {
			d.channels.responding <- r
		} else {
			d.channels.loading <- NewResponse(path, pict)
		}
	}

	return c, nil
}

func (d *Dispatcher) Start(ql, cl int) error {
	d.cache = *newCache(cl)
	d.awaiters = *newAwaiters()
	d.channels.saving = make(chan *Response, ql)
	d.channels.loading = make(chan *Response, ql)
	d.channels.processing = make(chan *Response, ql)
	d.channels.responding = make(chan *Response, ql)

	go d.process()

	for i := 0; i < ql; i++ {
		go d.load()
		go d.save()
		go d.send()
	}

	return nil
}

func (d *Dispatcher) load() {
	for r := range d.channels.loading {
		t := time.Now()
		r = load(r)
		r.Times["load"] = time.Since(t)

		if r.IsSuccessful() {
			d.channels.processing <- r
		} else {
			d.channels.responding <- r
		}
	}
}

func (d *Dispatcher) process() {
	for r := range d.channels.processing {
		t := time.Now()
		r = process(r)
		r.Times["load"] = time.Since(t)

		if r.IsSuccessful() {
			d.channels.processing <- r
		} else {
			d.channels.responding <- r
		}
	}
}

func (d *Dispatcher) save() {
	for r := range d.channels.saving {
		t := time.Now()
		r = save(r)
		r.Times["load"] = time.Since(t)

		if r.IsSuccessful() {
			d.channels.processing <- r
		} else {
			d.channels.responding <- r
		}
	}
}

func (d *Dispatcher) send() {
	for r := range d.channels.responding {
		if r.IsSuccessful() {
			d.cache.set(r.Path, r)
			d.channels.saving <- r
		}

		for {
			c := d.awaiters.pop(r.Path)
			if c == nil {
				break
			}

			c <- r
			close(c)
		}
	}
}

func NewDispatcher(pictures picture.Pictures) *Dispatcher {
	return &Dispatcher{
		pictures: pictures,
	}
}
