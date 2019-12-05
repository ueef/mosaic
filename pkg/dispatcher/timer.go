package dispatcher

import (
	"fmt"
	"time"
)

type Timer interface {
	Stop()
	Start(name string)
	fmt.Stringer
}

type pin struct {
	n string
	s time.Time
	d time.Duration
}

type timer struct {
	t []pin
}

func (t *timer) Stop() {
	for i := len(t.t) - 1; i >= 0; i-- {
		if t.t[i].d == 0 {
			t.t[i].d = time.Since(t.t[i].s)
			return
		}
	}
}

func (t *timer) Start(name string) {
	t.t = append(t.t, pin{
		n: name,
		s: time.Now(),
		d: 0,
	})
}

func (t timer) String() string {
	s := ""
	f := false
	for _, p := range t.t {
		s += p.n + ": " + p.d.String()
		if f {
			f = true
		} else {
			s += ", "
		}
	}

	return s
}

func NewTimer() Timer {
	return &timer{
		t: make([]pin, 0, 64),
	}
}
