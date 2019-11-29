package dispatcher

import "sync"

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
