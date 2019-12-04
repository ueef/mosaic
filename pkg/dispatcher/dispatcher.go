package dispatcher

import (
	"fmt"
	"github.com/ueef/mosaic/pkg/picture"
	"time"
)

type Dispatcher struct {
	c  cache
	a  awaiters
	p  picture.Pictures
	s  bool
	ch struct {
		s chan *Response
		l chan *Response
		p chan *Response
		r chan *Response
	}
}

func (d *Dispatcher) Dispatch(host, path string) (<-chan *Response, error) {
	if !d.s {
		return nil, fmt.Errorf("the dispatcher must be started before use")
	}

	pict, err := d.p.Match(host, path)
	if err != nil {
		return nil, err
	}

	c := make(chan *Response, 1)
	if d.a.push(path, c) {
		if r := d.c.get(path); r != nil {
			d.ch.r <- r
		} else {
			d.ch.l <- NewResponse(path, pict)
		}
	}

	return c, nil
}

func (d *Dispatcher) Start(ql, cl int) error {
	d.s = true
	d.c = *newCache(cl)
	d.a = *newAwaiters()
	d.ch.s = make(chan *Response, ql)
	d.ch.l = make(chan *Response, ql)
	d.ch.p = make(chan *Response, ql)
	d.ch.r = make(chan *Response, ql)

	go d.process()

	for i := 0; i < ql; i++ {
		go d.load()
		go d.save()
		go d.send()
	}

	return nil
}

func (d *Dispatcher) load() {
	for r := range d.ch.l {
		t := time.Now()
		r = load(r)
		r.Times["loading"] = time.Since(t)

		if r.IsSuccessful() {
			d.ch.p <- r
		} else {
			d.ch.r <- r
		}
	}
}

func (d *Dispatcher) process() {
	for r := range d.ch.p {
		t := time.Now()
		r = process(r)
		r.Times["processing"] = time.Since(t)

		if r.IsSuccessful() {
			d.ch.s <- r
		} else {
			d.ch.r <- r
		}
	}
}

func (d *Dispatcher) save() {
	for r := range d.ch.s {
		t := time.Now()
		r = save(r)
		r.Times["saving"] = time.Since(t)

		if r.IsSuccessful() {
			d.c.set(r.Path, r)
		}

		d.ch.r <- r
	}
}

func (d *Dispatcher) send() {
	for r := range d.ch.r {
		for {
			c := d.a.pop(r.Path)
			if c == nil {
				break
			}

			c <- r
			close(c)
		}
	}
}

func NewDispatcher(p picture.Pictures) *Dispatcher {
	return &Dispatcher{
		p: p,
	}
}
