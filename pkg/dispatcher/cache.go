package dispatcher

import "sync"

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
