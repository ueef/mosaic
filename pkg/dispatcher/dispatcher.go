package dispatcher

import (
	"bytes"
	"fmt"
	"github.com/ueef/mosaic/pkg/picture"
	"image"
	"sync"
)

type Response struct {
	Err  error
	Buff []byte
	Path string
	Pict *picture.Picture
}

func (r Response) IsSuccessful() bool {
	return r.Err == nil
}

func newResponse(path string, pict *picture.Picture) *Response {
	return &Response{
		Err:  nil,
		Buff: nil,
		Path: path,
		Pict: pict,
	}
}

func newErrorResponse(path string, err error) *Response {
	return &Response{
		Err:  err,
		Buff: nil,
		Path: path,
		Pict: nil,
	}
}

type Dispatcher struct {
	cache    cache
	awaiters awaiters
	pictures picture.Pictures
	channels struct {
		saving     chan *Response
		loading    chan *Response
		processing chan *Response
		responding chan *Response
		requesting chan *request
	}
}

func (d *Dispatcher) Dispatch(host, path string) (<-chan *Response, error) {
	pict, err := d.pictures.Match(host, path)
	if err != nil {
		return nil, err
	}

	c := make(chan *Response, 1)
	d.channels.requesting <- newRequest(path, pict, c)

	return c, nil
}

func (d *Dispatcher) Start(ql, cl int) error {
	d.cache = *newCache(cl)
	d.awaiters = *newAwaiters()
	d.channels.saving = make(chan *Response, ql)
	d.channels.loading = make(chan *Response, ql)
	d.channels.processing = make(chan *Response, ql)
	d.channels.responding = make(chan *Response, ql)
	d.channels.requesting = make(chan *request, ql*10)

	go d.process()

	for i := 0; i < ql; i++ {
		go d.handle()
		go d.load()
		go d.save()
		go d.send()
	}

	return nil
}

func (d *Dispatcher) handle() {
	for rq := range d.channels.requesting {
		r := d.cache.get(rq.path)
		if r != nil {
			rq.ch <- r
			close(rq.ch)
			continue
		}

		if d.awaiters.push(rq.path, rq.ch) {
			d.channels.loading <- newResponse(rq.path, rq.pict)
		}
	}
}

func (d *Dispatcher) load() {
	for r := range d.channels.loading {
		b, err := r.Pict.Loader.Load(r.Path)
		if err != nil {
			d.channels.responding <- newErrorResponse(r.Path, err)
			continue
		}
		r.Buff = b

		d.channels.processing <- r
	}
}

func (d *Dispatcher) process() {
	for r := range d.channels.processing {
		img, _, err := image.Decode(bytes.NewReader(r.Buff))
		if err != nil {
			d.channels.responding <- newErrorResponse(r.Path, err)
			continue
		}
		r.Buff = nil

		img, err = r.Pict.Filter.Apply(img)
		if err != nil {
			d.channels.responding <- newErrorResponse(r.Path, err)
			continue
		}

		r.Buff, err = r.Pict.Encoder.Encode(img)
		if err != nil {
			d.channels.responding <- newErrorResponse(r.Path, err)
		}

		d.channels.responding <- r
	}
}

func (d *Dispatcher) save() {
	for r := range d.channels.saving {
		err := r.Pict.Saver.Save(r.Path, r.Buff)
		if err != nil {
			fmt.Println(err)
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

type request struct {
	ch   chan *Response
	path string
	pict *picture.Picture
}

func newRequest(path string, pict *picture.Picture, ch chan *Response) *request {
	return &request{
		ch:   ch,
		path: path,
		pict: pict,
	}
}

type awaiters struct {
	m sync.Mutex
	r map[string][]chan *Response
}

func (a *awaiters) pop(k string) chan *Response {
	a.m.Lock()
	defer a.m.Unlock()

	r, ok := a.r[k]
	if !ok {
		return nil
	}

	c := r[0]
	r = r[1:]

	if len(r) > 0 {
		a.r[k] = r
	} else {
		delete(a.r, k)
	}

	return c
}

func (a *awaiters) push(p string, c chan *Response) bool {
	a.m.Lock()
	defer a.m.Unlock()

	r, ok := a.r[p]
	if ok {
		r = append(r, c)
	} else {
		r = []chan *Response{c}
	}
	a.r[p] = r

	return !ok
}

func newAwaiters() *awaiters {
	return &awaiters{
		m: sync.Mutex{},
		r: map[string][]chan *Response{},
	}
}

type cache struct {
	m sync.Mutex
	p []string
	r []*Response
	i map[string]int
	n int
	l int
}

func (c *cache) get(k string) *Response {
	c.m.Lock()
	defer c.m.Unlock()

	i, ok := c.i[k]
	if !ok {
		return nil
	}

	v := c.r[i]
	if i+1 < c.n {
		p := c.p[i+1]
		r := c.r[i+1]

		c.p[i+1] = k
		c.r[i+1] = v

		c.p[i] = p
		c.r[i] = r

		c.i[k]++
		c.i[p]--
	}

	return v
}

func (c *cache) set(k string, r *Response) {
	c.m.Lock()
	defer c.m.Unlock()

	i, ok := c.i[k]
	if ok {
		c.r[i] = r
		return
	}

	for c.n >= c.l {
		delete(c.i, c.p[0])
		c.p = c.p[1:]
		c.r = c.r[1:]
		c.n--
	}

	c.p = append(c.p, k)
	c.r = append(c.r, r)
	c.i[k] = c.n
	c.n++
}

func newCache(l int) *cache {
	return &cache{
		m: sync.Mutex{},
		p: []string{},
		r: []*Response{},
		i: map[string]int{},
		n: 0,
		l: l,
	}
}
