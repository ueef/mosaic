package filter

import (
	"github.com/anthonynsimon/bild/clone"
	"github.com/ueef/mosaic/pkg/parse"
	"github.com/ueef/mosaic/pkg/stamp"
	"image"
	"image/color"
	"math"
)

const TypeWatermark = "watermark"

type watermark struct {
	s  []stamp.Stamp
	c  color.Color
	cs int
}

func (f watermark) Apply(img image.Image) (image.Image, error) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = clone.AsRGBA(img)
	}

	g := newGrid(f.cs, rgba)
	g.clear()
	g.trim()

	f.drawWatermarks(g, rgba)

	return rgba, nil
}

func (f watermark) drawWatermarks(g *grid, i *image.RGBA) {
	cx, cy := g.x+g.w/2, g.y+g.h/2
	ml := g.w
	if g.w < g.h {
		ml = g.h
	}

	for _, s := range f.s {
		gx, gy := cx, cy
		f.drawWatermark(gx, gy, s, g, i)

		for l, v, vs := 2, &gx, 1; l <= ml; l += 2 {
			gy--
			gx--
			for r := 0; r < 4; r++ {
				switch r {
				case 0:
					v = &gx
					vs = 1
				case 1:
					v = &gy
					vs = 1
				case 2:
					v = &gx
					vs = -1
				case 3:
					v = &gy
					vs = -1
				}

				for ls := 0; ls < l; ls++ {
					*v += vs
					f.drawWatermark(gx, gy, s, g, i)
				}
			}
		}
	}
}

func (f watermark) drawWatermark(gx, gy int, s stamp.Stamp, g *grid, i *image.RGBA) {
	if !g.isCorrect(gx, gy) || g.get(gx, gy) {
		return
	}

	sw, sh := s.GetWidth(), s.GetHeight()
	gw, gh := sw/g.cs, sh/g.cs
	if sw%g.cs > 0 {
		gw++
	}
	if sh%g.cs > 0 {
		gh++
	}

	minx, miny := gx-gw/2, gy-gh/2
	maxx, maxy := minx+gw, miny+gh
	if !g.isCorrect(minx, miny) || !g.isCorrect(maxx-1, maxy-1) {
		return
	}

	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			if g.get(x, y) {
				return
			}
		}
	}

	s.Draw(minx*g.cs, miny*g.cs, f.c, i)

	mf := float64(sh) / float64(sw) * 6
	m := int(math.Round(float64(gw) * mf))
	minx = minx - m
	if minx < g.x {
		minx = g.x
	}

	maxx = maxx + m
	if maxx > g.x+g.w {
		maxx = g.x + g.w
	}

	miny = miny - m
	if miny < g.y {
		miny = g.y
	}

	maxy = maxy + m
	if maxy > g.y+g.h {
		maxy = g.y + g.h
	}

	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			g.set(x, y, true)
		}
	}
}

func NewWatermark(cs int, c color.Color, s []stamp.Stamp) Filter {
	if c == nil {
		c = color.Black
	}

	return &watermark{
		s:  s,
		c:  c,
		cs: cs,
	}
}

func NewWatermarkFromMap(m map[string]interface{}) (Filter, error) {
	cs, err := parse.GetRequiredIntFromMap("cell_size", m)
	if err != nil {
		return nil, err
	}

	v, err := parse.GetRequiredSliceOfInterfacesFromMap("stamps", m)
	if err != nil {
		return nil, err
	}

	c, _, err := parse.GetColorFromMap("color", m)
	if err != nil {
		return nil, err
	}

	s := make([]stamp.Stamp, len(v))
	for i := range v {
		s[i], err = stamp.NewFromConfig(v[i])
		if err != nil {
			return nil, err
		}
	}

	return NewWatermark(cs, c, s), nil
}

type grid struct {
	c  []bool
	x  int
	y  int
	w  int
	h  int
	s  int
	cs int
	xs int
	ys int
}

func (g grid) draw(i *image.RGBA) {
	for x := g.x; x < g.x+g.w; x++ {
		for y := g.y; y < g.y+g.h; y++ {
			for ix := x * g.cs; ix < (x+1)*g.cs; ix++ {
				for iy := y * g.cs; iy < (y+1)*g.cs; iy++ {
					p := iy*i.Stride + ix*4
					if g.get(x, y) {
						i.Pix[p+0] = 255
						i.Pix[p+1] = 255
						i.Pix[p+2] = 255
					} else {
						i.Pix[p+0] = 0
						i.Pix[p+1] = 0
						i.Pix[p+2] = 0
					}
				}
			}
		}
	}
}

func (g *grid) trim() {
	p := 0
	if g.w > g.h {
		p = g.h / 20
	} else {
		p = g.w / 20
	}

	pt, pb, pr, pl := -1, -1, -1, -1
	for y := g.y; y < g.y+g.h; y++ {
		for x := g.x; x < g.x+g.w; x++ {
			if !g.get(x, y) {
				continue
			}

			if -1 == pt {
				pt = y
			}

			if y > pb {
				pb = y
			}

			if -1 == pl || x < pl {
				pl = x
			}

			if x > pr {
				pr = x
			}
		}
	}

	if pt < p {
		pt = p
	}

	if pb == -1 || pb > g.h-p {
		pb = g.h - p
	}

	if pl < p {
		pl = p
	}

	if pr == -1 || pr > g.w-p {
		pr = g.w - p
	}

	g.x = pl
	g.y = pt
	g.w = pr - pl
	g.h = pb - pt
}

func (g *grid) replace(s, pv int, patterns [][]uint8) {
	a := s * s
	f := make([]uint8, s*s)
	for x := g.x; x < g.x+g.w-s; x++ {
		for y := g.y; y < g.y+g.h-s; y++ {

			fi := 0
			im := make([][2]int, 0, s*s)
			for sy := 0; sy < s; sy++ {
				for sx := 0; sx < s; sx++ {
					if g.get(x+sx, y+sy) {
						f[fi] = 1
						im = append(im, [2]int{x + sx, y + sy})
					} else {
						f[fi] = 0
					}
					fi++
				}
			}

			for _, p := range patterns {
				if len(p) != a {
					continue
				}

				ok := true
				v := 0
				for i := range f {
					if f[i] != p[i] {
						v++
						if v > pv {
							ok = false
							break
						}
					}
				}

				if !ok {
					continue
				}

				for i := range im {
					g.set(im[i][0], im[i][1], false)
				}
			}
		}
	}
}

func (g *grid) clear() {

	g.replace(4, 1, [][]uint8{
		{
			0, 1, 0, 0,
			0, 1, 0, 0,
			0, 1, 0, 0,
			0, 1, 0, 0,
		},
		{
			0, 0, 0, 0,
			1, 1, 1, 1,
			0, 0, 0, 0,
			0, 0, 0, 0,
		},

		{
			0, 0, 0, 0,
			1, 1, 1, 0,
			0, 0, 1, 0,
			0, 0, 1, 0,
		},
		{
			0, 0, 1, 0,
			0, 0, 1, 0,
			1, 1, 1, 0,
			0, 0, 0, 0,
		},
		{
			0, 1, 0, 0,
			0, 1, 0, 0,
			0, 1, 1, 1,
			0, 0, 0, 0,
		},
		{
			0, 0, 0, 0,
			0, 1, 1, 1,
			0, 1, 0, 0,
			0, 1, 0, 0,
		},

		{
			0, 1, 0, 0,
			1, 1, 1, 1,
			0, 1, 0, 0,
			0, 1, 0, 0,
		},
	})

	g.replace(3, 0, [][]uint8{
		{
			0, 0, 0,
			0, 1, 0,
			0, 0, 0,
		},
	})
}

func (g grid) isCorrect(x, y int) bool {
	return x >= g.x && x < g.x+g.w && y >= g.y && y < g.y+g.h
}

func (g grid) correctX(x int) int {
	if x < g.x {
		return g.x
	}

	m := g.x + g.w - 1
	if x > m {
		return m
	}

	return x
}

func (g grid) correctY(y int) int {
	if y < g.y {
		return g.y
	}

	m := g.y + g.h - 1
	if y > m {
		return m
	}

	return y
}

func (g grid) get(x, y int) bool {
	return g.c[y*g.s+x]
}

func (g *grid) set(x, y int, v bool) {
	g.c[y*g.s+x] = v
}

func newGrid(cs int, img *image.RGBA) *grid {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	x, y := w%cs/2, h%cs/2

	gw, gh := w/cs, h/cs
	g := make([]bool, gw*gh)

	s := threshold(img)
	for gy := 0; gy < gh; gy++ {
		miny := b.Min.Y + gy*cs
		maxy := miny + cs

		for gx := 0; gx < gw; gx++ {
			minx := b.Min.X + gx*cs
			maxx := minx + cs

		loop:
			for y := miny; y < maxy; y++ {
				for x := minx; x < maxx; x++ {
					si := y*w + x
					if s[si] > 0 {
						g[gy*gw+gx] = true
						break loop
					}
				}
			}
		}
	}

	return &grid{
		c:  g,
		x:  0,
		y:  0,
		w:  gw,
		h:  gh,
		s:  gw,
		xs: x,
		ys: y,
		cs: cs,
	}
}

func threshold(img *image.RGBA) []int {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	a := w * h
	s := make([]int, a)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		i := y * w
		pi := y*img.Stride + b.Min.X*4
		for x := b.Min.X; x < b.Max.X; x++ {
			s[i] = math.MaxUint8 - int(0.299*float64(img.Pix[pi+0])+0.587*float64(img.Pix[pi+1])+0.114*float64(img.Pix[pi+2]))
			i++
			pi += 4
		}
	}

	cs := 32
	ca := cs * cs
	cw, ch := w/cs, h/cs
	if w%cs > 0 {
		cw++
	}
	if h%cs > 0 {
		ch++
	}

	for cy := 0; cy < ch; cy++ {
		miny := b.Min.Y + cy*cs
		maxy := miny + cs
		if maxy > b.Max.Y {
			maxy = b.Max.Y
		}

		for cx := 0; cx < cw; cx++ {
			minx := b.Min.X + cx*cs
			maxx := minx + cs
			if maxx >= b.Max.X {
				maxx = b.Max.X
			}

			min, max, avg := math.MaxUint8, 0, 0
			for y := miny; y < maxy; y++ {
				for x := minx; x < maxx; x++ {
					si := y*w + x
					if s[si] < min {
						min = s[si]
					}
					if s[si] > max {
						max = s[si]
					}
					avg += s[si]
				}
			}

			avg /= ca
			d := max - avg
			avg += d / 2
			for y := miny; y < maxy; y++ {
				for x := minx; x < maxx; x++ {
					si := y*w + x

					if d < 32 {
						s[si] = 0
					} else if s[si] > avg {
						s[si] = 255
					} else {
						s[si] = 0
					}

				}
			}
		}
	}

	return s
}

func init() {
	RegisterFilter(TypeWatermark, NewWatermarkFromMap)
}
